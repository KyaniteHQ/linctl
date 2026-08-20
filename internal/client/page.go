package client

// Page is the cursor state every Linear connection read returns alongside its
// items. List types embed it so the JSON shape stays flat: encoding/json
// promotes an anonymous embedded struct's fields into the outer object.
type Page struct {
	HasNextPage bool    `json:"has_next_page"`
	EndCursor   *string `json:"end_cursor,omitempty"`
}

// HasMore reports whether Linear has further pages beyond the returned items.
// It is the constraint the CLI list pipeline uses instead of reflecting over
// an unconstrained page type.
func (page Page) HasMore() bool {
	return page.HasNextPage
}
