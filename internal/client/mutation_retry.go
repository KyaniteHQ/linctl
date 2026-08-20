package client

// MutationRetryClass is the retry policy for one guarded write. Linear's
// IssueUpdateInput and IssueRelationCreateInput have no stable idempotency
// key (vendored schema.graphql; Linear SDK pin in .linear-sdk-ref). Blind
// replay of an ambiguous result can duplicate a relation or hide a failed
// state change.
type MutationRetryClass string

const (
	// MutationRetryIdempotent means the same input can be sent again.
	MutationRetryIdempotent MutationRetryClass = "idempotent"
	// MutationRetryReconcile means an ambiguous result must be read back
	// before another write. That is the class for issue state assignment
	// and issueRelationCreate.
	MutationRetryReconcile MutationRetryClass = "reconcile-before-retry"
	// MutationRetryNever means the caller must not send the mutation again.
	MutationRetryNever MutationRetryClass = "never-retry"
)

// IssueStateWriteRetryClass is the retry class for issue state assignment.
func IssueStateWriteRetryClass() MutationRetryClass {
	return MutationRetryReconcile
}

// IssueRelationCreateRetryClass is the retry class for issueRelationCreate.
func IssueRelationCreateRetryClass() MutationRetryClass {
	return MutationRetryReconcile
}
