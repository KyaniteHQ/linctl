package client

// MutationRetryClass is the retry policy for one guarded write. Linear's
// IssueUpdateInput and IssueRelationCreateInput have no stable idempotency
// key (vendored schema.graphql; Linear SDK pin in .linear-sdk-ref). Blind
// replay of an ambiguous result can duplicate a relation or hide a failed
// state change.
type MutationRetryClass string

// MutationRetryReconcile means an ambiguous result must be read back
// before another write. That is the class for issue state assignment
// and issueRelationCreate.
const MutationRetryReconcile MutationRetryClass = "reconcile-before-retry"

// IssueStateWriteRetryClass is the retry class for issue state assignment.
func IssueStateWriteRetryClass() MutationRetryClass {
	return MutationRetryReconcile
}

// IssueRelationCreateRetryClass is the retry class for issueRelationCreate.
func IssueRelationCreateRetryClass() MutationRetryClass {
	return MutationRetryReconcile
}

// WorkflowStateWriteRetryClass is the retry class for workflow-state writes.
func WorkflowStateWriteRetryClass() MutationRetryClass {
	return MutationRetryReconcile
}

// TemplateWriteRetryClass is the retry class for issue-template writes.
func TemplateWriteRetryClass() MutationRetryClass {
	return MutationRetryReconcile
}

func applyMutationRetryClass[T any](
	class MutationRetryClass,
	recovered T,
	recoveredOK bool,
	writeErr error,
) (T, error) {
	var zero T
	if class == MutationRetryReconcile && recoveredOK {
		return recovered, nil
	}

	return zero, writeErr
}
