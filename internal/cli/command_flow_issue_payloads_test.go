package cli

import (
	"strings"
)

func commandFlowIssuePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowIssueReadPayload(operation, fake); ok {
		return payload, true
	}

	return commandFlowIssueWritePayload(operation, fake)
}

//nolint:gocyclo // The command-flow fake is intentionally centralized by operation name.
func commandFlowAttachmentIssuePayload(operation string) (string, bool) {
	if !strings.HasPrefix(operation, "attachmentIssue") {
		return "", false
	}

	switch operation {
	case "attachmentIssue":
		return `{"attachmentIssue":` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`}`, true
	case "attachmentIssue_attachments":
		return `{"attachmentIssue":{"attachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_botActor":
		return `{"attachmentIssue":{"id":"issue-id","botActor":` + commandActorBotJSON() + `}}`, true
	case "attachmentIssue_children":
		return `{"attachmentIssue":{"children":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_comments":
		return `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[` +
			commandCommentMetadataJSON("issue-id", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_documents":
		return `{"attachmentIssue":{"documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_formerAttachments":
		return `{"attachmentIssue":{"formerAttachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_formerNeeds":
		return `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","formerNeeds":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_history":
		return `{"attachmentIssue":{"history":{"nodes":[` +
			commandIssueHistoryJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_inverseRelations":
		return `{"attachmentIssue":{"inverseRelations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_labels":
		return `{"attachmentIssue":{"labels":{"nodes":[` +
			commandLabelJSON("label body") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_needs":
		return `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","needs":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_relations":
		return `{"attachmentIssue":{"relations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_releases":
		return `{"attachmentIssue":{"releases":{"nodes":[` +
			commandReleaseJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_sharedAccess":
		return commandIssueSharedAccessPayload("attachmentIssue"), true
	case "attachmentIssue_stateHistory":
		return `{"attachmentIssue":{"id":"issue-id","stateHistory":{"nodes":[` +
			commandIssueStateSpanJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_subscribers":
		return `{"attachmentIssue":{"id":"issue-id","subscribers":{"nodes":[` +
			commandUserJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

//nolint:gocyclo // The command-flow fake is intentionally centralized by operation name.
func commandFlowIssueVCSBranchPayload(operation string) (string, bool) {
	if !strings.HasPrefix(operation, "issueVcsBranchSearch") {
		return "", false
	}

	switch operation {
	case "issueVcsBranchSearch":
		return `{"issueVcsBranchSearch":` +
			commandIssueJSON("LIT-40", "Branch issue", "todo-state", "Todo", "unstarted") +
			`}`, true
	case "issueVcsBranchSearch_attachments":
		return `{"issueVcsBranchSearch":{"attachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_botActor":
		return `{"issueVcsBranchSearch":{"id":"issue-id","botActor":` + commandActorBotJSON() + `}}`, true
	case "issueVcsBranchSearch_children":
		return `{"issueVcsBranchSearch":{"children":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_documents":
		return `{"issueVcsBranchSearch":{"documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_formerAttachments":
		return `{"issueVcsBranchSearch":{"formerAttachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_comments":
		return `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[` +
			commandCommentMetadataJSON("issue-id", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_formerNeeds":
		return `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","formerNeeds":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_history":
		return `{"issueVcsBranchSearch":{"history":{"nodes":[` +
			commandIssueHistoryJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_inverseRelations":
		return `{"issueVcsBranchSearch":{"inverseRelations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_labels":
		return `{"issueVcsBranchSearch":{"labels":{"nodes":[` +
			commandLabelJSON("label body") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_needs":
		return `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","needs":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_relations":
		return `{"issueVcsBranchSearch":{"relations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_releases":
		return `{"issueVcsBranchSearch":{"releases":{"nodes":[` +
			commandReleaseJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_sharedAccess":
		return commandIssueSharedAccessPayload("issueVcsBranchSearch"), true
	case "issueVcsBranchSearch_stateHistory":
		return `{"issueVcsBranchSearch":{"id":"issue-id","stateHistory":{"nodes":[` +
			commandIssueStateSpanJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_subscribers":
		return `{"issueVcsBranchSearch":{"id":"issue-id","subscribers":{"nodes":[` +
			commandUserJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueRelationPayload(operation string) (string, bool) {
	switch operation {
	case "issueRelations":
		return `{"issueRelations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issueRelation":
		return `{"issueRelation":` + commandIssueRelationJSON() + `}`, true
	case "IssueRelationCreate":
		return `{"issueRelationCreate":{"success":true,"issueRelation":{` +
			`"id":"issue-relation-id","type":"related",` +
			`"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,` +
			`"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},` +
			`"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}}}}`, true
	case "IssueRelationDelete":
		return `{"issueRelationDelete":{"success":true,"entityId":"issue-relation-id"}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueReadPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowIssueListPayload(operation, fake); ok {
		return payload, true
	}
	if payload, ok := commandFlowIssueChildPayload(operation, fake); ok {
		return payload, true
	}
	if payload, ok := commandFlowIssueUtilityPayload(operation, fake); ok {
		return payload, true
	}
	if payload, ok := commandFlowIssueRelationPayload(operation); ok {
		return payload, true
	}

	switch operation {
	case "issueSearch":
		if fake.emptyIssueSearch {
			return `{"issueSearch":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"issueSearch":{"nodes":[` + commandIssueJSON("LIT-3", "Search result", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "NextIssuesByTeam":
		if fake.emptyNextIssues {
			return emptyCommandIssuesPayload(), true
		}
		if fake.rankedNextIssues {
			return `{"issues":{"nodes":[` +
				commandIssueWithNextRankJSON("LIT-28", "Low priority standalone", 4, "Low", "2026-05-01T12:00:00Z", 0) + `,` +
				commandIssueWithNextRankJSON("LIT-29", "Urgent standalone", 1, "Urgent", "2026-06-01T12:00:00Z", 0) + `,` +
				commandIssueWithNextRankJSON("LIT-30", "Unblocks checkout", 2, "High", "2026-06-10T12:00:00Z", 2) +
				`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"issues":{"nodes":[` + commandIssueWithNextRankJSON("LIT-27", "Next issue", 0, "No priority", "2026-06-01T12:00:00Z", 0) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issue":
		return `{"issue":` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `}`, true
	case "IssueDependencies":
		return commandFlowIssueDependenciesPayload(), true
	case "issue_comments":
		if fake.emptyIssueComments {
			return `{"issue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		hasNextPage := "false"
		if fake.truncatedExport {
			hasNextPage = "true"
		}
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[{"id":"comment-id","body":"First comment","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","parentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":` + hasNextPage + `,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueUtilityPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "issueFigmaFileKeySearch":
		if fake.emptyIssueFigmaSearch {
			return `{"issueFigmaFileKeySearch":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"issueFigmaFileKeySearch":{"nodes":[` +
			commandIssueJSON("LIT-41", "Figma issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issuePriorityValues":
		return `{"issuePriorityValues":[{"priority":1,"label":"Urgent"},{"priority":0,"label":"No priority"}]}`, true
	case "issueFilterSuggestion":
		return `{"issueFilterSuggestion":{"filter":{"state":{"type":{"eq":"started"}}},"logId":"issue-filter-log-id"}}`, true
	case "issueTitleSuggestionFromCustomerRequest":
		return `{"issueTitleSuggestionFromCustomerRequest":{"title":"Improve exports","logId":"title-log-id"}}`, true
	default:
		return "", false
	}
}

//nolint:gocyclo // The command-flow fake is intentionally centralized by operation name.
func commandFlowIssueChildPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "issue_attachments":
		return `{"issue":{"attachments":{"nodes":[` + commandAttachmentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_botActor":
		return `{"issue":{"id":"issue-id","botActor":` + commandActorBotJSON() + `}}`, true
	case "issue_children":
		if fake.emptyIssueChildren {
			return `{"issue":{"children":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"issue":{"children":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_documents":
		return `{"issue":{"documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_formerAttachments":
		return `{"issue":{"formerAttachments":{"nodes":[` + commandAttachmentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_formerNeeds":
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","formerNeeds":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_history":
		return `{"issue":{"history":{"nodes":[` + commandIssueHistoryJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_inverseRelations":
		return `{"issue":{"inverseRelations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_labels":
		return `{"issue":{"labels":{"nodes":[` + commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_needs":
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","needs":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_relations":
		return `{"issue":{"relations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_releases":
		return `{"issue":{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_sharedAccess":
		return commandIssueSharedAccessPayload("issue"), true
	case "issue_stateHistory":
		return `{"issue":{"id":"issue-id","stateHistory":{"nodes":[` +
			commandIssueStateSpanJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_subscribers":
		return `{"issue":{"id":"issue-id","subscribers":{"nodes":[` +
			commandUserJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandIssueSharedAccessPayload(root string) string {
	return `{"` + root + `":{"id":"issue-id","identifier":"LIT-1","sharedAccess":` +
		commandIssueSharedAccessJSON() + `}}`
}

func commandFlowIssueDependenciesPayload() string {
	return `{"issue":{
		"id":"issue-id",
		"identifier":"LIT-1",
		"parent":` + commandIssueJSON("LIT-25", "Parent issue", "todo-state", "Todo", "unstarted") + `,
		"children":{
			"nodes":[` + commandIssueJSON("LIT-26", "Child issue", "todo-state", "Todo", "unstarted") + `],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		},
		"relations":{
			"nodes":[{"id":"blocks-relation","type":"blocks","relatedIssue":` + commandIssueJSON("LIT-23", "Blocked issue", "todo-state", "Todo", "unstarted") + `}],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		},
		"inverseRelations":{
			"nodes":[{"id":"blocked-by-relation","type":"blocks","issue":` + commandIssueJSON("LIT-24", "Blocker issue", "todo-state", "Todo", "unstarted") + `}],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		}
	}}`
}

func commandFlowIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowBroadIssueListPayload(operation, fake); ok {
		return payload, true
	}
	if payload, ok := commandFlowDependencyIssueListPayload(operation, fake); ok {
		return payload, true
	}

	switch operation {
	case "IssuesByTeamState":
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-2", "Started issue", "started-state", "Started", fake.expectedStateType) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamProject":
		if fake.emptyIssueProject {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-4", "Project issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamAssignee":
		return commandFlowAssigneeIssueListPayload(fake), true
	case "IssuesByTeamLabel":
		if fake.emptyIssueLabel {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-7", "Labeled issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCycle":
		if fake.emptyIssueCycle {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-8", "Cycle issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCreatedAfter":
		if fake.emptyIssueCreatedAfter {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-9", "Recent issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCreatedBefore":
		if fake.emptyIssueCreatedBefore {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-19", "Older issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	default:
		return "", false
	}
}

func commandFlowDependencyIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "IssuesByTeamHasBlockers":
		if fake.emptyIssueHasBlockers {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-21", "Blocked issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamBlocks":
		if fake.emptyIssueBlocks {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-22", "Blocking issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssueBlockedIssues":
		if fake.emptyIssueBlockedBy {
			return `{"issue":{"id":"issue-id","identifier":"LIT-1","relations":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","relations":{"nodes":[{"id":"relation-id","type":"blocks","relatedIssue":` + commandIssueJSON("LIT-23", "Blocked by issue", "todo-state", "Todo", "unstarted") + `}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowBroadIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "issues":
		if fake.emptyIssueAllTeams {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-20", "All-team issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeam":
		return commandFlowUnfilteredIssueListPayload(fake), true
	default:
		return "", false
	}
}

func commandFlowUnfilteredIssueListPayload(fake commandFlowFakeClient) string {
	if fake.emptyIssueList {
		return emptyCommandIssuesPayload()
	}
	if fake.multiIssueList {
		return `{"issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Alpha issue", "todo-state", "Todo", "unstarted") + `,` +
			commandIssueJSON("LIT-2", "Zebra issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
	}

	return `{"issues":{"nodes":[` + commandIssueJSON("LIT-1", "Listed issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func commandFlowAssigneeIssueListPayload(fake commandFlowFakeClient) string {
	if fake.emptyIssueMine {
		return emptyCommandIssuesPayload()
	}
	if fake.expectedAssigneeID == "assignee-id" {
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-6", "Assigned issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
	}

	return `{"issues":{"nodes":[` + commandIssueJSON("LIT-5", "Mine issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func emptyCommandIssuesPayload() string {
	return `{"issues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func commandFlowIssueWritePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "IssueCreate":
		return `{"issueCreate":{"success":true,"issue":` + commandIssueJSON("LIT-2", "Created issue", "todo-state", "Todo", "unstarted") + `}}`, true
	case "IssueUpdate":
		if fake.expectedStartStateID != "" {
			return `{"issueUpdate":{"success":true,"issue":` +
				commandIssueJSON("LIT-1", "Started issue", "started-state", "Started", "started") + `}}`, true
		}
		return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSON("LIT-1", "Updated issue", "todo-state", "Todo", "unstarted") + `}}`, true
	case "IssueCommentCreate":
		return `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"Looks good","url":"https://linear.app/comment/comment-id","issue":` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `}}}`, true
	case "AttachmentLinkURL":
		return `{"attachmentCreate":{"success":true,"attachment":{"id":"attachment-id","title":"Linked PR","subtitle":null,"url":"https://example.com/pr/1","sourceType":null}}}`, true
	case "CompletedWorkflowStates":
		return `{"workflowStates":{"nodes":[{"id":"done-state","name":"Done","type":"completed","position":1}]}}`, true
	case "StartedWorkflowStates":
		return `{"workflowStates":{"nodes":[{"id":"started-state","name":"Started","type":"started","position":1}]}}`, true
	case "WorkflowStatesByType":
		return `{"workflowStates":{"nodes":[{"id":"type-state-id","name":"TypeState","type":"unstarted","position":1}]}}`, true
	case "IssueClose":
		return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSON("LIT-1", "Closed issue", "done-state", "Done", "completed") + `}}`, true
	default:
		return "", false
	}
}
