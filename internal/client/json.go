package client

import "encoding/json"

// countJSONArrayEntries returns the number of top-level entries in a JSON array,
// or zero when raw is empty or is not a JSON array.
func countJSONArrayEntries(raw json.RawMessage) int {
	entryCount := 0
	var entries []json.RawMessage
	if err := json.Unmarshal(raw, &entries); err == nil {
		entryCount = len(entries)
	}

	return entryCount
}
