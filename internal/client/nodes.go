package client

// mapNodes maps a page of GraphQL connection nodes onto domain summaries.
func mapNodes[Node any, Summary any](nodes []Node, mapOne func(Node) Summary) []Summary {
	summaries := make([]Summary, 0, len(nodes))
	for _, node := range nodes {
		summaries = append(summaries, mapOne(node))
	}

	return summaries
}
