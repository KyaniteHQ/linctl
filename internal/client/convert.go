package client

// Pointer and value conversion helpers shared by generated-client call sites.

func intPtr(value int) *int {
	return &value
}

func boolPtr(value bool) *bool {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
