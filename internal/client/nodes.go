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
	Items       []Item
	HasNextPage bool
	EndCursor   *string
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

// listConnection wraps collectNodePages for the common shape: fetch one page of a
// connection, map its nodes to a Summary, and thread its page info through unchanged.
func listConnection[Node, Summary any](
	operation string,
	limit int,
	pageSize int,
	fetch func(pageSize int, after *string) (nodes []Node, hasNextPage bool, endCursor *string, err error),
	mapOne func(Node) Summary,
) (nodePage[Summary], error) {
	return collectNodePages(operation, limit, pageSize, func(pageSize int, after *string) (nodePage[Summary], error) {
		nodes, hasNextPage, endCursor, err := fetch(pageSize, after)
		if err != nil {
			return nodePage[Summary]{}, err
		}

		return nodePage[Summary]{Items: mapNodes(nodes, mapOne), HasNextPage: hasNextPage, EndCursor: endCursor}, nil
	})
}

// defaultListPageSize matches Linear's per-request cap, the same authority
// issueListPageSize cites at issue.go:134.
const defaultListPageSize = 250
