// Package gqlmodel holds neutral GraphQL input models shared by code generation
// and the parent client compatibility aliases.
package gqlmodel

import "encoding/json"

// LinearAttachmentCreateInput is the sparse Linear attachmentCreate payload linctl supports.
type LinearAttachmentCreateInput struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
	URL      string  `json:"url"`
	IssueID  string  `json:"issueId"`
}

// LinearCommentCreateInput is the sparse Linear commentCreate payload linctl supports.
type LinearCommentCreateInput struct {
	Body     *string `json:"body,omitempty"`
	IssueID  *string `json:"issueId,omitempty"`
	ParentID *string `json:"parentId,omitempty"`
}

// LinearCommentUpdateInput is the sparse Linear commentUpdate payload linctl supports.
type LinearCommentUpdateInput struct {
	Body *string `json:"body,omitempty"`
}

// LinearCycleCreateInput is the sparse Linear cycleCreate payload linctl supports.
type LinearCycleCreateInput struct {
	TeamID      string  `json:"teamId"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	StartsAt    string  `json:"startsAt"`
	EndsAt      string  `json:"endsAt"`
	CompletedAt *string `json:"completedAt,omitempty"`
}

// LinearCycleUpdateInput is the sparse Linear cycleUpdate payload linctl supports.
type LinearCycleUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	StartsAt    *string `json:"startsAt,omitempty"`
	EndsAt      *string `json:"endsAt,omitempty"`
	CompletedAt *string `json:"completedAt,omitempty"`
}

// LinearDocumentCreateInput is the sparse Linear documentCreate payload linctl supports.
type LinearDocumentCreateInput struct {
	Title     string  `json:"title"`
	Content   *string `json:"content,omitempty"`
	TeamID    *string `json:"teamId,omitempty"`
	ProjectID *string `json:"projectId,omitempty"`
}

// LinearDocumentUpdateInput is the sparse Linear documentUpdate payload linctl supports.
type LinearDocumentUpdateInput struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

// LinearIssueCreateInput is the sparse Linear issueCreate payload linctl supports.
type LinearIssueCreateInput struct {
	Title              *string  `json:"title,omitempty"`
	Description        *string  `json:"description,omitempty"`
	TeamID             string   `json:"teamId"`
	ProjectID          *string  `json:"projectId,omitempty"`
	StateID            *string  `json:"stateId,omitempty"`
	Priority           *int     `json:"priority,omitempty"`
	AssigneeID         *string  `json:"assigneeId,omitempty"`
	LabelIDs           []string `json:"labelIds,omitempty"`
	DueDate            *string  `json:"dueDate,omitempty"`
	Estimate           *int     `json:"estimate,omitempty"`
	ParentID           *string  `json:"parentId,omitempty"`
	ProjectMilestoneID *string  `json:"projectMilestoneId,omitempty"`
}

// LinearIssueFilter is the sparse Linear IssueFilter linctl composes for team-scoped issue listing.
type LinearIssueFilter struct {
	Team                  *LinearIDFilter                 `json:"team,omitempty"`
	State                 *LinearWorkflowStateTypeFilter  `json:"state,omitempty"`
	Project               *LinearIDFilter                 `json:"project,omitempty"`
	Assignee              *LinearIDFilter                 `json:"assignee,omitempty"`
	Labels                *LinearLabelCollectionFilter    `json:"labels,omitempty"`
	Cycle                 *LinearIDFilter                 `json:"cycle,omitempty"`
	CreatedAt             *LinearDateComparator           `json:"createdAt,omitempty"`
	HasBlockedByRelations *LinearRelationExistsComparator `json:"hasBlockedByRelations,omitempty"`
	HasBlockingRelations  *LinearRelationExistsComparator `json:"hasBlockingRelations,omitempty"`
}

// LinearIDFilter matches an entity by id equality.
type LinearIDFilter struct {
	ID LinearIDComparator `json:"id"`
}

// LinearIDComparator is the sparse Linear IDComparator linctl supports.
type LinearIDComparator struct {
	Eq string `json:"eq"`
}

// LinearWorkflowStateTypeFilter matches workflow states by type equality.
type LinearWorkflowStateTypeFilter struct {
	Type LinearStringComparator `json:"type"`
}

// LinearStringComparator is the sparse Linear StringComparator linctl supports.
type LinearStringComparator struct {
	Eq string `json:"eq"`
}

// LinearLabelCollectionFilter matches issues carrying at least one matching label.
type LinearLabelCollectionFilter struct {
	Some LinearIDFilter `json:"some"`
}

// LinearDateComparator is the sparse Linear DateComparator linctl supports.
type LinearDateComparator struct {
	Gte *string `json:"gte,omitempty"`
	Lte *string `json:"lte,omitempty"`
}

// LinearRelationExistsComparator is the sparse Linear RelationExistsComparator linctl supports.
type LinearRelationExistsComparator struct {
	Eq bool `json:"eq"`
}

// LinearIssueLabelCreateInput is the sparse Linear issueLabelCreate payload linctl supports.
type LinearIssueLabelCreateInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	ParentID    *string `json:"parentId,omitempty"`
	TeamID      *string `json:"teamId,omitempty"`
}

// LinearIssueLabelUpdateInput is the sparse Linear issueLabelUpdate payload linctl supports.
type LinearIssueLabelUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// LinearIssueRelationCreateInput is the sparse Linear issueRelationCreate payload linctl supports.
type LinearIssueRelationCreateInput struct {
	Type           string `json:"type"`
	IssueID        string `json:"issueId"`
	RelatedIssueID string `json:"relatedIssueId"`
}

// LinearIssueUpdateInput is the sparse Linear issueUpdate payload linctl supports.
// DueDate and Estimate use RawMessage so explicit null clears the value while
// an absent value leaves it untouched.
type LinearIssueUpdateInput struct {
	Title              *string         `json:"title,omitempty"`
	Description        *string         `json:"description,omitempty"`
	AssigneeID         *string         `json:"assigneeId,omitempty"`
	DelegateID         *string         `json:"delegateId,omitempty"`
	StateID            *string         `json:"stateId,omitempty"`
	Priority           *int            `json:"priority,omitempty"`
	LabelIDs           []string        `json:"labelIds,omitempty"`
	DueDate            json.RawMessage `json:"dueDate,omitempty"`
	Estimate           json.RawMessage `json:"estimate,omitempty"`
	ProjectMilestoneID json.RawMessage `json:"projectMilestoneId,omitempty"`
}

// LinearProjectCreateInput is the sparse Linear projectCreate payload linctl supports.
type LinearProjectCreateInput struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Content     *string  `json:"content,omitempty"`
	TeamIDs     []string `json:"teamIds"`
}

// LinearProjectLabelCreateInput is the sparse Linear projectLabelCreate payload linctl supports.
type LinearProjectLabelCreateInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// LinearProjectLabelUpdateInput is the sparse Linear projectLabelUpdate payload linctl supports.
type LinearProjectLabelUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// LinearProjectMilestoneCreateInput is the sparse Linear projectMilestoneCreate payload linctl supports.
type LinearProjectMilestoneCreateInput struct {
	ProjectID   string  `json:"projectId"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	TargetDate  *string `json:"targetDate,omitempty"`
}

// LinearProjectMilestoneUpdateInput is the sparse Linear projectMilestoneUpdate payload linctl supports.
type LinearProjectMilestoneUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	TargetDate  *string `json:"targetDate,omitempty"`
}

// LinearProjectUpdateCreateInput is the sparse Linear projectUpdateCreate payload linctl supports.
type LinearProjectUpdateCreateInput struct {
	ProjectID string  `json:"projectId"`
	Body      *string `json:"body,omitempty"`
	Health    *string `json:"health,omitempty"`
}

// LinearProjectUpdateInput is the sparse Linear projectUpdate payload linctl supports.
type LinearProjectUpdateInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Content     *string  `json:"content,omitempty"`
	TeamIDs     []string `json:"teamIds,omitempty"`
}
