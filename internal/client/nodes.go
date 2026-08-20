package client

import "fmt"

// mapNodes maps a page of GraphQL connection nodes onto domain summaries.
func mapNodes[Node any, Summary any](nodes []Node, mapOne func(Node) Summary) []Summary {
	summaries := make([]Summary, 0, len(nodes))
	for _, node := range nodes {
		summaries = append(summaries, mapOne(node))
	}

	return summaries
}

// nodePage is one accumulated connection page: its items plus cursor state.
type nodePage[Item any] struct {
	Items []Item
	Page
}

// collectNodePages accumulates connection items until limit items are
// collected or pages run out, capping each request at pageSize so large
// --limit values page through Linear's per-request cap instead of silently
// truncating. The returned page carries every accumulated item plus the last
// page's cursor state; errors are prefixed with the operation context.
func collectNodePages[Item any](
	operation string,
	limit int,
	pageSize int,
	fetch func(pageSize int, after *string) (nodePage[Item], error),
) (nodePage[Item], error) {
	if limit < 1 {
		return nodePage[Item]{}, fmt.Errorf("%s: limit must be positive", operation)
	}

	items := make([]Item, 0, limit)
	var after *string
	for {
		page, err := fetch(min(pageSize, limit-len(items)), after)
		if err != nil {
			return nodePage[Item]{}, fmt.Errorf("%s: %w", operation, err)
		}
		items = append(items, page.Items...)

		if len(items) >= limit || !page.HasNextPage {
			page.Items = items
			return page, nil
		}
		if page.EndCursor == nil {
			return nodePage[Item]{}, fmt.Errorf("%s: next page has no end cursor: %w", operation, ErrGraphQL)
		}
		after = page.EndCursor
	}
}

// listConnection is listConnectionWithParent for the common shape: a connection
// with no denormalized parent to carry back.
func listConnection[Node, Summary any](
	operation string,
	limit int,
	pageSize int,
	fetch func(pageSize int, after *string) (nodes []Node, hasNextPage bool, endCursor *string, err error),
	mapOne func(Node) Summary,
) (nodePage[Summary], error) {
	withoutParent := func(pageSize int, after *string) ([]Node, struct{}, bool, *string, error) {
		nodes, hasNextPage, endCursor, err := fetch(pageSize, after)

		return nodes, struct{}{}, hasNextPage, endCursor, err
	}
	page, _, err := listConnectionWithParent(operation, limit, pageSize, withoutParent, mapOne)

	return page, err
}

// listConnectionWithParent wraps collectNodePages: fetch one page, map its
// nodes to a Summary, thread its page info through, and carry back the parent
// entity the read denormalizes into every page (a project id, an issue
// identifier, a team key). The fetch returns that parent instead of
// writing it into a captured struct field, so the data flow stays visible in
// the signature. Linear repeats the parent on every page, so the last page's
// value is the one returned; a zero Parent comes back with any error.
func listConnectionWithParent[Node, Summary, Parent any](
	operation string,
	limit int,
	pageSize int,
	fetch func(pageSize int, after *string) ([]Node, Parent, bool, *string, error),
	mapOne func(Node) Summary,
) (nodePage[Summary], Parent, error) {
	var parent Parent
	collect := func(pageSize int, after *string) (nodePage[Summary], error) {
		nodes, pageParent, hasNextPage, endCursor, err := fetch(pageSize, after)
		if err != nil {
			return nodePage[Summary]{}, err
		}
		parent = pageParent

		return nodePage[Summary]{
			Items: mapNodes(nodes, mapOne),
			Page:  Page{HasNextPage: hasNextPage, EndCursor: endCursor},
		}, nil
	}
	page, err := collectNodePages(operation, limit, pageSize, collect)
	if err != nil {
		var zero Parent
		return nodePage[Summary]{}, zero, err
	}

	return page, parent, nil
}

// issueParent, projectParent, and labelParent are the connection parents that
// several reads share. identifier, projectName, and labelName stay empty when
// the operation's selection set does not carry them; the consuming list type
// has no field for them in that case.
type issueParent struct {
	issueID    string
	identifier string
}

type projectParent struct {
	projectID   string
	projectName string
}

type labelParent struct {
	labelID   string
	labelName string
}

// defaultListPageSize matches Linear's per-request cap, the same authority
// issueListPageSize cites at issue.go:134.
const defaultListPageSize = 250
