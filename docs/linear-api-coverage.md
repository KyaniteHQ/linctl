# Linear API coverage ledger

Generated from current local sources and upstream Linear SDK commit `28d6461`.

Sources (paths relative to the upstream Linear SDK checkout):

- Upstream SDK methods: `packages/sdk/src/_generated_sdk.ts`
- Upstream schema roots: `packages/sdk/src/schema.graphql`
- Local generated operations: `internal/client/generated.go`
- Local GraphQL operations: `internal/client/operations/*.graphql`
- Repo domain map: `docs/domain-map.md`

Statuses: `implemented`, `accepted_gap`, `safe_candidate`, `blocked_needs_design`, `intentionally_excluded`.

## Summary

| Surface | Total | Implemented/root-backed | Classified |
| --- | ---: | ---: | ---: |
| Upstream SDK root methods | 458 | 133 | 458 |
| Upstream Query root fields | 158 | 113 | 158 |
| Upstream Mutation root fields | 364 | 21 | 364 |
| Local generated Go operations | 332 | 332 | 332 |
| Domain-map commands | 421 | 295 | 421 |

## Upstream SDK Root Methods

| Method | Kind | Status | Evidence |
| --- | --- | --- | --- |
| `administrableTeams` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentActivities` | method | implemented | local operation or command exists |
| `agentActivity` | method | implemented | local operation or command exists |
| `agentSession` | method | implemented | local operation or command exists |
| `agentSessionCreateOnComment` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnIssue` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionUpdateExternalUrl` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessions` | method | implemented | local operation or command exists |
| `agentSkill` | method | implemented | local operation or command exists |
| `agentSkills` | method | implemented | local operation or command exists |
| `airbyteIntegrationConnect` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `applicationInfo` | method | implemented | local operation or command exists |
| `archiveCustomerNeed` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveCycle` | method | implemented | local operation or command exists |
| `archiveInitiative` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveInitiativeUpdate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveIntegration` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `archiveIssue` | method | implemented | local operation or command exists |
| `archiveNotification` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveProject` | method | implemented | local operation or command exists |
| `archiveProjectStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveProjectUpdate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveRelease` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveReleasePipeline` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveReleaseStage` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveRoadmap` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archiveWorkflowState` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `archivedIntegrations` | getter | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `attachment` | method | implemented | local operation or command exists |
| `attachmentIssue` | method | implemented | local operation or command exists |
| `attachmentLinkDiscord` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkFront` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkGitHubIssue` | method | blocked_needs_design | attachment-to-GitHub linking mutates third-party integration state; needs explicit integration guard semantics |
| `attachmentLinkGitHubPR` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkGitLabMR` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkIntercom` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkJiraIssue` | method | blocked_needs_design | attachment-to-Jira linking mutates third-party integration state; needs explicit integration guard semantics |
| `attachmentLinkSalesforce` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkSlack` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkURL` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkZendesk` | method | blocked_needs_design | mutation needs product and safety design |
| `attachmentSyncToSlack` | method | blocked_needs_design | mutation needs product and safety design |
| `attachments` | method | implemented | local operation or command exists |
| `attachmentsForURL` | method | implemented | local operation or command exists |
| `auditEntries` | method | blocked_needs_design | audit logs can expose actor, IP, country, and request metadata; needs explicit admin/security output model |
| `auditEntryTypes` | getter | implemented | local operation or command exists |
| `authenticationSessions` | getter | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `availableUsers` | getter | intentionally_excluded | available-user picker enumeration is a specialized product resolver; `user list` is the supported user read surface |
| `comment` | method | implemented | local operation or command exists |
| `commentResolve` | method | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `commentUnresolve` | method | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `comments` | method | implemented | local operation or command exists |
| `constructor` | method | blocked_needs_design | SDK method is not matched to a GraphQL root field; explicit classification required |
| `createAgentActivity` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createAgentSkill` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createAttachment` | method | implemented | local operation or command exists |
| `createComment` | method | implemented | local operation or command exists |
| `createContact` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCsvExportReport` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomView` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomer` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomerNeed` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomerStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomerTier` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCycle` | method | implemented | local operation or command exists |
| `createDocument` | method | implemented | local operation or command exists |
| `createEmailIntakeAddress` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createEmoji` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createEntityExternalLink` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createFavorite` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createGitAutomationState` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createGitAutomationTargetBranch` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createInitiative` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createInitiativeRelation` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createInitiativeToProject` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createInitiativeUpdate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createInitiativeUpdateReminder` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createIntegrationGithubCommit` | getter | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `createIntegrationTemplate` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `createIntegrationsSettings` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `createIssue` | method | implemented | local operation or command exists |
| `createIssueBatch` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createIssueLabel` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createIssueRelation` | method | implemented | local operation or command exists |
| `createIssueToRelease` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createNotificationSubscription` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createOrganizationInvite` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProject` | method | implemented | local operation or command exists |
| `createProjectLabel` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProjectMilestone` | method | implemented | local operation or command exists |
| `createProjectRelation` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProjectStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProjectUpdate` | method | implemented | local operation or command exists |
| `createProjectUpdateReminder` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createPushSubscription` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createReaction` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createRelease` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createReleaseNote` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createReleasePipeline` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createReleaseStage` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createRoadmap` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createRoadmapToProject` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createTeam` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createTeamMembership` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createTemplate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createTimeSchedule` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createTriageResponsibility` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createViewPreferences` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createWebhook` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `createWorkflowState` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customView` | method | implemented | local operation or command exists |
| `customViewHasSubscribers` | method | implemented | local operation or command exists |
| `customViews` | method | implemented | local operation or command exists |
| `customer` | method | implemented | local operation or command exists |
| `customerMerge` | method | blocked_needs_design | mutation needs product and safety design |
| `customerNeed` | method | implemented | local operation or command exists |
| `customerNeedCreateFromAttachment` | method | blocked_needs_design | mutation needs product and safety design |
| `customerNeeds` | method | implemented | local operation or command exists |
| `customerStatus` | method | implemented | local operation or command exists |
| `customerStatuses` | method | implemented | local operation or command exists |
| `customerTier` | method | implemented | local operation or command exists |
| `customerTiers` | method | implemented | local operation or command exists |
| `customerUnsync` | method | blocked_needs_design | mutation needs product and safety design |
| `customerUpsert` | method | blocked_needs_design | mutation needs product and safety design |
| `customers` | method | implemented | local operation or command exists |
| `cycle` | method | implemented | local operation or command exists |
| `cycleShiftAll` | method | blocked_needs_design | bulk Cycle date shifting is a state-changing organization operation that needs target-pinned guard semantics |
| `cycleStartUpcomingCycleToday` | method | blocked_needs_design | starting an upcoming Cycle changes team planning state and needs target-pinned guard semantics |
| `cycles` | method | implemented | local operation or command exists |
| `deleteAgentSkill` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteAttachment` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteComment` | method | implemented | local operation or command exists |
| `deleteCustomView` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteCustomer` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteCustomerNeed` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteCustomerStatus` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteCustomerTier` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteDocument` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteEmailIntakeAddress` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteEmoji` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteEntityExternalLink` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteFavorite` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteGitAutomationState` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteGitAutomationTargetBranch` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteInitiative` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteInitiativeRelation` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteInitiativeToProject` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteIntegration` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteIntegrationIntercom` | getter | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteIntegrationTemplate` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteIssue` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteIssueImport` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteIssueLabel` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteIssueRelation` | method | implemented | local operation or command exists |
| `deleteIssueToRelease` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteNotificationSubscription` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteOrganization` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteOrganizationCancel` | getter | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteOrganizationDomain` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteOrganizationInvite` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteProject` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteProjectLabel` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteProjectMilestone` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteProjectRelation` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteProjectUpdate` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deletePushSubscription` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteReaction` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteRelease` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteReleaseNote` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteReleasePipeline` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteRoadmap` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteRoadmapToProject` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteTeam` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteTeamCycles` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteTeamKey` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteTeamMembership` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteTemplate` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteTimeSchedule` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteTriageResponsibility` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteViewPreferences` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteWebhook` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `document` | method | implemented | local operation or command exists |
| `documentContentHistory` | method | blocked_needs_design | content, thread, and archive payload reads can expose body/blob data; needs explicit opt-in projection before CLI exposure |
| `documents` | method | implemented | local operation or command exists |
| `emailIntakeAddress` | method | intentionally_excluded | email intake administration sits outside the ordinary agent CLI read surface |
| `emailIntakeAddressRefreshSesDomainStatus` | method | blocked_needs_design | mutation needs product and safety design |
| `emailIntakeAddressRotate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `emailTokenUserAccountAuth` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emailUnsubscribe` | method | blocked_needs_design | mutation needs product and safety design |
| `emailUserAccountAuthChallenge` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emoji` | method | implemented | local operation or command exists |
| `emojis` | method | implemented | local operation or command exists |
| `entityExternalLink` | method | implemented | local operation or command exists |
| `externalUser` | method | implemented | local operation or command exists |
| `externalUsers` | method | implemented | local operation or command exists |
| `favorite` | method | implemented | local operation or command exists |
| `favorites` | method | implemented | local operation or command exists |
| `fileUpload` | method | implemented | local operation or command exists |
| `googleUserAccountAuth` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `imageUploadFromUrl` | method | blocked_needs_design | mutation needs product and safety design |
| `importFileUpload` | method | blocked_needs_design | mutation needs product and safety design |
| `initiative` | method | implemented | local operation or command exists |
| `initiativeRelation` | method | implemented | local operation or command exists |
| `initiativeRelations` | method | implemented | local operation or command exists |
| `initiativeToProject` | method | implemented | local operation or command exists |
| `initiativeToProjects` | method | implemented | local operation or command exists |
| `initiativeUpdate` | method | implemented | local operation or command exists |
| `initiativeUpdates` | method | implemented | local operation or command exists |
| `initiatives` | method | implemented | local operation or command exists |
| `integration` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationAsksConnectChannel` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationDiscord` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationFigma` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationFront` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGitHubEnterpriseServerConnect` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGitHubPersonal` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubConnect` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubImportConnect` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubImportRefresh` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubRemoveCodeAccess` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `integrationGitlabConnect` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGitlabTestConnection` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGong` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGoogleSheets` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationHasScopes` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationIntercom` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationJiraPersonal` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationLoom` | getter | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationMicrosoftPersonalConnect` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationMicrosoftTeams` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationRequest` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSalesforce` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSentryConnect` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlack` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackAsks` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackCustomViewNotifications` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackCustomerChannelLink` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackImportEmojis` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackOrAsksUpdateSlackTeamName` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackOrgProjectUpdatesPost` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackPersonal` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackPost` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackProjectPost` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationTemplate` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationTemplates` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationZendesk` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrations` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationsSettings` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `issue` | method | implemented | local operation or command exists |
| `issueAddLabel` | method | blocked_needs_design | issue label mutation needs issue target pinning and target-mismatch tests |
| `issueExternalSyncDisable` | method | blocked_needs_design | issue external-sync disable changes integration state and needs explicit integration guard semantics |
| `issueFigmaFileKeySearch` | method | implemented | local operation or command exists |
| `issueFilterSuggestion` | method | implemented | local operation or command exists |
| `issueImportCheckCSV` | method | blocked_needs_design | CSV import validation can expose imported row payloads and needs an explicit redaction/output model |
| `issueImportCheckSync` | method | blocked_needs_design | sync import validation can expose external tracker payloads and needs an explicit redaction/output model |
| `issueImportCreateAsana` | method | blocked_needs_design | Asana issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateCSVJira` | method | blocked_needs_design | CSV/Jira issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateClubhouse` | method | blocked_needs_design | Clubhouse issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateGithub` | method | blocked_needs_design | GitHub issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateJira` | method | blocked_needs_design | Jira issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportJqlCheck` | method | blocked_needs_design | JQL import validation can expose external tracker payloads and needs an explicit redaction/output model |
| `issueImportProcess` | method | blocked_needs_design | issue import processing advances external import workflow state and needs explicit integration guard semantics |
| `issueLabel` | method | implemented | local operation or command exists |
| `issueLabelRestore` | method | blocked_needs_design | issue label lifecycle restore needs explicit organization/admin safety semantics |
| `issueLabelRetire` | method | blocked_needs_design | issue label lifecycle retire needs explicit organization/admin safety semantics |
| `issueLabels` | method | implemented | local operation or command exists |
| `issuePriorityValues` | getter | implemented | local operation or command exists |
| `issueRelation` | method | implemented | local operation or command exists |
| `issueRelations` | method | implemented | local operation or command exists |
| `issueReminder` | method | blocked_needs_design | issue reminder mutation changes notification state and needs target-pinned guard semantics |
| `issueRemoveLabel` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueRepositorySuggestions` | method | intentionally_excluded | repository suggestion reads expose VCS integration metadata outside the default Linear work CLI surface |
| `issueSearch` | method | implemented | local operation or command exists |
| `issueShare` | method | blocked_needs_design | issue sharing changes access state and needs target-pinned guard semantics |
| `issueSubscribe` | method | blocked_needs_design | issue subscription changes notification state and needs target-pinned guard semantics |
| `issueTitleSuggestionFromCustomerRequest` | method | implemented | local operation or command exists |
| `issueToRelease` | method | implemented | local operation or command exists |
| `issueToReleaseDeleteByIssueAndRelease` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueToReleases` | method | implemented | local operation or command exists |
| `issueUnshare` | method | blocked_needs_design | issue unsharing changes access state and needs target-pinned guard semantics |
| `issueUnsubscribe` | method | blocked_needs_design | issue unsubscribe changes notification state and needs target-pinned guard semantics |
| `issueVcsBranchSearch` | method | implemented | local operation or command exists |
| `issues` | method | implemented | local operation or command exists |
| `latestReleaseByAccessKey` | getter | intentionally_excluded | access-key release reads are unauthenticated sharing surfaces outside the token-scoped agent CLI |
| `logout` | method | blocked_needs_design | mutation needs product and safety design |
| `logoutAllSessions` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `logoutOtherSessions` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `logoutSession` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `notification` | method | implemented | local operation or command exists |
| `notificationArchiveAll` | method | blocked_needs_design | mutation needs product and safety design |
| `notificationMarkReadAll` | method | blocked_needs_design | mutation needs product and safety design |
| `notificationMarkUnreadAll` | method | blocked_needs_design | mutation needs product and safety design |
| `notificationSnoozeAll` | method | blocked_needs_design | mutation needs product and safety design |
| `notificationSubscription` | method | implemented | local operation or command exists |
| `notificationSubscriptions` | method | implemented | local operation or command exists |
| `notificationUnsnoozeAll` | method | blocked_needs_design | mutation needs product and safety design |
| `notifications` | method | implemented | local operation or command exists |
| `organization` | getter | implemented | local operation or command exists |
| `organizationDeleteChallenge` | getter | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationExists` | method | implemented | local operation or command exists |
| `organizationInvite` | method | intentionally_excluded | organization invite reads can expose invitee and admin metadata outside an agent-safe CLI surface |
| `organizationInvites` | method | intentionally_excluded | organization invite reads can expose invitee and admin metadata outside an agent-safe CLI surface |
| `organizationStartTrial` | getter | blocked_needs_design | mutation needs product and safety design |
| `organizationStartTrialForPlan` | method | blocked_needs_design | mutation needs product and safety design |
| `project` | method | implemented | local operation or command exists |
| `projectAddLabel` | method | blocked_needs_design | project label mutation needs project target pinning and target-mismatch tests |
| `projectExternalSyncDisable` | method | blocked_needs_design | project external-sync disable changes integration state and needs explicit integration guard semantics |
| `projectFilterSuggestion` | method | implemented | local operation or command exists |
| `projectLabel` | method | implemented | local operation or command exists |
| `projectLabelRestore` | method | blocked_needs_design | project label lifecycle restore needs explicit organization/admin safety semantics |
| `projectLabelRetire` | method | blocked_needs_design | project label lifecycle retire needs explicit organization/admin safety semantics |
| `projectLabels` | method | implemented | local operation or command exists |
| `projectMilestone` | method | implemented | local operation or command exists |
| `projectMilestones` | method | implemented | local operation or command exists |
| `projectRelation` | method | implemented | local operation or command exists |
| `projectRelations` | method | implemented | local operation or command exists |
| `projectRemoveLabel` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectStatus` | method | implemented | local operation or command exists |
| `projectStatuses` | method | implemented | local operation or command exists |
| `projectUpdate` | method | implemented | local operation or command exists |
| `projectUpdates` | method | implemented | local operation or command exists |
| `projects` | method | implemented | local operation or command exists |
| `pushSubscriptionTest` | method | intentionally_excluded | push subscription diagnostics are notification-device integration plumbing outside the CLI surface |
| `rateLimitStatus` | getter | implemented | local operation or command exists |
| `recentReleasesByAccessKey` | method | intentionally_excluded | access-key release reads are unauthenticated sharing surfaces outside the token-scoped agent CLI |
| `refreshGoogleSheetsData` | method | blocked_needs_design | mutation needs product and safety design |
| `release` | method | implemented | local operation or command exists |
| `releaseComplete` | method | blocked_needs_design | mutation needs product and safety design |
| `releaseCompleteByAccessKey` | method | blocked_needs_design | mutation needs product and safety design |
| `releaseNote` | method | implemented | local operation or command exists |
| `releaseNotes` | method | implemented | local operation or command exists |
| `releasePipeline` | method | implemented | local operation or command exists |
| `releasePipelineByAccessKey` | getter | intentionally_excluded | access-key release reads are unauthenticated sharing surfaces outside the token-scoped agent CLI |
| `releasePipelines` | method | implemented | local operation or command exists |
| `releaseSearch` | method | implemented | local operation or command exists |
| `releaseStage` | method | implemented | local operation or command exists |
| `releaseStages` | method | implemented | local operation or command exists |
| `releaseSync` | method | blocked_needs_design | mutation needs product and safety design |
| `releaseSyncByAccessKey` | method | blocked_needs_design | mutation needs product and safety design |
| `releaseUpdateByPipeline` | method | blocked_needs_design | mutation needs product and safety design |
| `releaseUpdateByPipelineByAccessKey` | method | blocked_needs_design | mutation needs product and safety design |
| `releases` | method | implemented | local operation or command exists |
| `resendOrganizationInvite` | method | blocked_needs_design | mutation needs product and safety design |
| `resendOrganizationInviteByEmail` | method | blocked_needs_design | mutation needs product and safety design |
| `roadmap` | method | implemented | local operation or command exists |
| `roadmapToProject` | method | implemented | local operation or command exists |
| `roadmapToProjects` | method | implemented | local operation or command exists |
| `roadmaps` | method | implemented | local operation or command exists |
| `rotateSecretWebhook` | method | blocked_needs_design | SDK method is not matched to a GraphQL root field; explicit classification required |
| `samlTokenUserAccountAuth` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `searchDocuments` | method | implemented | local operation or command exists |
| `searchIssues` | method | implemented | local operation or command exists |
| `searchProjects` | method | implemented | local operation or command exists |
| `semanticSearch` | method | implemented | local operation or command exists |
| `slaConfigurations` | method | implemented | local operation or command exists |
| `ssoUrlFromEmail` | method | intentionally_excluded | SSO discovery from email belongs to auth flow tooling, not the Linear work CLI |
| `suspendUser` | method | blocked_needs_design | SDK method is not matched to a GraphQL root field; explicit classification required |
| `team` | method | implemented | local operation or command exists |
| `teamMembership` | method | implemented | local operation or command exists |
| `teamMemberships` | method | implemented | local operation or command exists |
| `teams` | method | implemented | local operation or command exists |
| `template` | method | implemented | local operation or command exists |
| `templates` | getter | implemented | local operation or command exists |
| `templatesForIntegration` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeSchedule` | method | implemented | local operation or command exists |
| `timeScheduleRefreshIntegrationSchedule` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeScheduleUpsertExternal` | method | blocked_needs_design | mutation needs product and safety design |
| `timeSchedules` | method | implemented | local operation or command exists |
| `trackAnonymousEvent` | method | blocked_needs_design | mutation needs product and safety design |
| `triageResponsibilities` | method | implemented | local operation or command exists |
| `triageResponsibility` | method | implemented | local operation or command exists |
| `unarchiveCustomerNeed` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveDocument` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveInitiative` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveInitiativeUpdate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveIssue` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveNotification` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveProject` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveProjectStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveProjectUpdate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveRelease` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveReleasePipeline` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveReleaseStage` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveRoadmap` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unarchiveTeam` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `unsuspendUser` | method | blocked_needs_design | SDK method is not matched to a GraphQL root field; explicit classification required |
| `updateAgentSession` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `updateAgentSkill` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateAttachment` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateComment` | method | implemented | local operation or command exists |
| `updateCustomView` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomer` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomerNeed` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomerStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomerTier` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCycle` | method | implemented | local operation or command exists |
| `updateDocument` | method | implemented | local operation or command exists |
| `updateEmailIntakeAddress` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateEntityExternalLink` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateFavorite` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateGitAutomationState` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateGitAutomationTargetBranch` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateInitiative` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateInitiativeRelation` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateInitiativeToProject` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateInitiativeUpdate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateIntegrationIntercomSettings` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `updateIntegrationsSettings` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `updateIssue` | method | implemented | local operation or command exists |
| `updateIssueBatch` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateIssueImport` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateIssueLabel` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateIssueRelation` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateNotification` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateNotificationCategoryChannelSubscription` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateNotificationSubscription` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateOrganization` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateOrganizationInvite` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateProject` | method | implemented | local operation or command exists |
| `updateProjectLabel` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateProjectMilestone` | method | implemented | local operation or command exists |
| `updateProjectRelation` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateProjectStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateProjectUpdate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateRelease` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateReleaseNote` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateReleasePipeline` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateReleaseStage` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateRoadmap` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateRoadmapToProject` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateTeam` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateTeamMembership` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateTemplate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateTimeSchedule` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateTriageResponsibility` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateUser` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateUserFlag` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateUserSettings` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateViewPreferences` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateWebhook` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `updateWorkflowState` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `user` | method | implemented | local operation or command exists |
| `userChangeRole` | method | intentionally_excluded | user role changes are organization administration outside the ordinary agent CLI surface |
| `userDiscordConnect` | method | intentionally_excluded | Discord account connection belongs to user auth/integration setup, not work CLI reads |
| `userExternalUserDisconnect` | method | intentionally_excluded | external-user disconnection is identity integration administration outside the ordinary agent CLI surface |
| `userRevokeAllSessions` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userRevokeSession` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userSessions` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `userSettings` | getter | implemented | local operation or command exists |
| `userSettingsFlagsReset` | method | intentionally_excluded | user settings flag reset is internal preference administration outside the ordinary agent CLI surface |
| `userUnlinkFromIdentityProvider` | method | intentionally_excluded | identity-provider unlinking is auth administration outside the ordinary agent CLI surface |
| `users` | method | implemented | local operation or command exists |
| `verifyGitHubEnterpriseServerInstallation` | method | intentionally_excluded | GitHub Enterprise installation verification is integration administration outside the CLI surface |
| `viewer` | getter | implemented | local operation or command exists |
| `webhook` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `webhooks` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `workflowState` | method | implemented | local operation or command exists |
| `workflowStates` | method | implemented | local operation or command exists |

## Upstream Query Root Fields

| Field | Return type | Status | Evidence |
| --- | --- | --- | --- |
| `_dummy` | `String!` | safe_candidate | read operation may fit future CLI coverage |
| `administrableTeams` | `TeamConnection!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentActivities` | `AgentActivityConnection!` | implemented | root field used by local GraphQL operation |
| `agentActivity` | `AgentActivity!` | implemented | root field used by local GraphQL operation |
| `agentSession` | `AgentSession!` | implemented | root field used by local GraphQL operation |
| `agentSessionSandbox` | `CodingAgentSandboxPayload` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessions` | `AgentSessionConnection!` | implemented | root field used by local GraphQL operation |
| `agentSkill` | `AgentSkill!` | implemented | root field used by local GraphQL operation |
| `agentSkills` | `AgentSkillConnection!` | implemented | root field used by local GraphQL operation |
| `applicationInfo` | `Application!` | implemented | root field used by local GraphQL operation |
| `archivedIntegrations` | `[ArchivedIntegrationPayload!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `archivedTeams` | `[Team!]!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `attachment` | `Attachment!` | implemented | root field used by local GraphQL operation |
| `attachmentIssue` | `Issue!` | implemented | root field used by local GraphQL operation |
| `attachmentSources` | `AttachmentSourcesPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `attachments` | `AttachmentConnection!` | implemented | root field used by local GraphQL operation |
| `attachmentsForURL` | `AttachmentConnection!` | implemented | root field used by local GraphQL operation |
| `auditEntries` | `AuditEntryConnection!` | blocked_needs_design | audit logs can expose actor, IP, country, and request metadata; needs explicit admin/security output model |
| `auditEntryTypes` | `[AuditEntryType!]!` | implemented | root field used by local GraphQL operation |
| `authenticationSessions` | `[AuthenticationSessionResponse!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `availableUsers` | `AuthResolverResponse!` | intentionally_excluded | available-user picker enumeration is a specialized product resolver; `user list` is the supported user read surface |
| `comment` | `Comment!` | implemented | root field used by local GraphQL operation |
| `comments` | `CommentConnection!` | implemented | root field used by local GraphQL operation |
| `customView` | `CustomView!` | implemented | root field used by local GraphQL operation |
| `customViewDetailsSuggestion` | `CustomViewSuggestionPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `customViewHasSubscribers` | `CustomViewHasSubscribersPayload!` | implemented | root field used by local GraphQL operation |
| `customViews` | `CustomViewConnection!` | implemented | root field used by local GraphQL operation |
| `customer` | `Customer!` | implemented | root field used by local GraphQL operation |
| `customerNeed` | `CustomerNeed!` | implemented | root field used by local GraphQL operation |
| `customerNeeds` | `CustomerNeedConnection!` | implemented | root field used by local GraphQL operation |
| `customerStatus` | `CustomerStatus!` | implemented | root field used by local GraphQL operation |
| `customerStatuses` | `CustomerStatusConnection!` | implemented | root field used by local GraphQL operation |
| `customerTier` | `CustomerTier!` | implemented | root field used by local GraphQL operation |
| `customerTiers` | `CustomerTierConnection!` | implemented | root field used by local GraphQL operation |
| `customers` | `CustomerConnection!` | implemented | root field used by local GraphQL operation |
| `cycle` | `Cycle!` | implemented | root field used by local GraphQL operation |
| `cycles` | `CycleConnection!` | implemented | root field used by local GraphQL operation |
| `document` | `Document!` | implemented | root field used by local GraphQL operation |
| `documentContentHistory` | `DocumentContentHistoryPayload!` | blocked_needs_design | content, thread, and archive payload reads can expose body/blob data; needs explicit opt-in projection before CLI exposure |
| `documentContentHistoryEntries` | `DocumentContentHistoryPayload!` | blocked_needs_design | content, thread, and archive payload reads can expose body/blob data; needs explicit opt-in projection before CLI exposure |
| `documentContentHistoryTimeline` | `DocumentContentHistoryTimelinePayload!` | blocked_needs_design | content, thread, and archive payload reads can expose body/blob data; needs explicit opt-in projection before CLI exposure |
| `documents` | `DocumentConnection!` | implemented | root field used by local GraphQL operation |
| `emailIntakeAddress` | `EmailIntakeAddress!` | intentionally_excluded | email intake administration sits outside the ordinary agent CLI read surface |
| `emoji` | `Emoji!` | implemented | root field used by local GraphQL operation |
| `emojis` | `EmojiConnection!` | implemented | root field used by local GraphQL operation |
| `entityExternalLink` | `EntityExternalLink!` | implemented | root field used by local GraphQL operation |
| `externalUser` | `ExternalUser!` | implemented | root field used by local GraphQL operation |
| `externalUsers` | `ExternalUserConnection!` | implemented | root field used by local GraphQL operation |
| `failuresForOauthWebhooks` | `[WebhookFailureEvent!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `favorite` | `Favorite!` | implemented | root field used by local GraphQL operation |
| `favorites` | `FavoriteConnection!` | implemented | root field used by local GraphQL operation |
| `fetchData` | `FetchDataPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `initiative` | `Initiative!` | implemented | root field used by local GraphQL operation |
| `initiativeRelation` | `InitiativeRelation!` | implemented | root field used by local GraphQL operation |
| `initiativeRelations` | `InitiativeRelationConnection!` | implemented | root field used by local GraphQL operation |
| `initiativeToProject` | `InitiativeToProject!` | implemented | root field used by local GraphQL operation |
| `initiativeToProjects` | `InitiativeToProjectConnection!` | implemented | root field used by local GraphQL operation |
| `initiativeUpdate` | `InitiativeUpdate!` | implemented | root field used by local GraphQL operation |
| `initiativeUpdates` | `InitiativeUpdateConnection!` | implemented | root field used by local GraphQL operation |
| `initiatives` | `InitiativeConnection!` | implemented | root field used by local GraphQL operation |
| `integration` | `Integration!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationHasScopes` | `IntegrationHasScopesPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationTemplate` | `IntegrationTemplate!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationTemplates` | `IntegrationTemplateConnection!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrations` | `IntegrationConnection!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationsSettings` | `IntegrationsSettings!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `issue` | `Issue!` | implemented | root field used by local GraphQL operation |
| `issueFigmaFileKeySearch` | `IssueConnection!` | implemented | root field used by local GraphQL operation |
| `issueFilterSuggestion` | `IssueFilterSuggestionPayload!` | implemented | root field used by local GraphQL operation |
| `issueImportCheckCSV` | `IssueImportCheckPayload!` | blocked_needs_design | CSV import validation can expose imported row payloads and needs an explicit redaction/output model |
| `issueImportCheckSync` | `IssueImportSyncCheckPayload!` | blocked_needs_design | sync import validation can expose external tracker payloads and needs an explicit redaction/output model |
| `issueImportJqlCheck` | `IssueImportJqlCheckPayload!` | blocked_needs_design | JQL import validation can expose external tracker payloads and needs an explicit redaction/output model |
| `issueLabel` | `IssueLabel!` | implemented | root field used by local GraphQL operation |
| `issueLabels` | `IssueLabelConnection!` | implemented | root field used by local GraphQL operation |
| `issuePriorityValues` | `[IssuePriorityValue!]!` | implemented | root field used by local GraphQL operation |
| `issueRelation` | `IssueRelation!` | implemented | root field used by local GraphQL operation |
| `issueRelations` | `IssueRelationConnection!` | implemented | root field used by local GraphQL operation |
| `issueRepositorySuggestions` | `RepositorySuggestionsPayload!` | intentionally_excluded | repository suggestion reads expose VCS integration metadata outside the default Linear work CLI surface |
| `issueSearch` | `IssueConnection!` | implemented | root field used by local GraphQL operation |
| `issueTitleSuggestionFromCustomerRequest` | `IssueTitleSuggestionFromCustomerRequestPayload!` | implemented | root field used by local GraphQL operation |
| `issueToRelease` | `IssueToRelease!` | implemented | root field used by local GraphQL operation |
| `issueToReleases` | `IssueToReleaseConnection!` | implemented | root field used by local GraphQL operation |
| `issueVcsBranchSearch` | `Issue` | implemented | root field used by local GraphQL operation |
| `issues` | `IssueConnection!` | implemented | root field used by local GraphQL operation |
| `latestReleaseByAccessKey` | `Release` | intentionally_excluded | access-key release reads are unauthenticated sharing surfaces outside the token-scoped agent CLI |
| `microsoftTeamsChannels` | `MicrosoftTeamsChannelsPayload!` | intentionally_excluded | Microsoft Teams channel enumeration exposes chat integration metadata outside the default Linear work CLI surface |
| `notification` | `Notification!` | implemented | root field used by local GraphQL operation |
| `notificationSubscription` | `NotificationSubscription!` | implemented | root field used by local GraphQL operation |
| `notificationSubscriptions` | `NotificationSubscriptionConnection!` | implemented | root field used by local GraphQL operation |
| `notifications` | `NotificationConnection!` | implemented | root field used by local GraphQL operation |
| `notificationsUnreadCount` | `Int!` | safe_candidate | read operation may fit future CLI coverage |
| `oauthApplication` | `OAuthApplication!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplications` | `[OAuthApplication!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `organization` | `Organization!` | implemented | root field used by local GraphQL operation |
| `organizationDomainClaimRequest` | `OrganizationDomainClaimPayload!` | intentionally_excluded | organization domain claim requests expose org-admin domain-verification metadata outside the ordinary agent CLI surface |
| `organizationExists` | `OrganizationExistsPayload!` | implemented | root field used by local GraphQL operation |
| `organizationInvite` | `OrganizationInvite!` | intentionally_excluded | organization invite reads can expose invitee and admin metadata outside an agent-safe CLI surface |
| `organizationInviteDetails` | `OrganizationInviteDetailsPayload!` | intentionally_excluded | organization invite reads can expose invitee and admin metadata outside an agent-safe CLI surface |
| `organizationInvites` | `OrganizationInviteConnection!` | intentionally_excluded | organization invite reads can expose invitee and admin metadata outside an agent-safe CLI surface |
| `organizationMeta` | `OrganizationMeta` | safe_candidate | read operation may fit future CLI coverage |
| `project` | `Project!` | implemented | root field used by local GraphQL operation |
| `projectFilterSuggestion` | `ProjectFilterSuggestionPayload!` | implemented | root field used by local GraphQL operation |
| `projectLabel` | `ProjectLabel!` | implemented | root field used by local GraphQL operation |
| `projectLabels` | `ProjectLabelConnection!` | implemented | root field used by local GraphQL operation |
| `projectMilestone` | `ProjectMilestone!` | implemented | root field used by local GraphQL operation |
| `projectMilestones` | `ProjectMilestoneConnection!` | implemented | root field used by local GraphQL operation |
| `projectRelation` | `ProjectRelation!` | implemented | root field used by local GraphQL operation |
| `projectRelations` | `ProjectRelationConnection!` | implemented | root field used by local GraphQL operation |
| `projectStatus` | `ProjectStatus!` | implemented | root field used by local GraphQL operation |
| `projectStatusProjectCount` | `ProjectStatusCountPayload!` | implemented | root field used by local GraphQL operation |
| `projectStatuses` | `ProjectStatusConnection!` | implemented | root field used by local GraphQL operation |
| `projectUpdate` | `ProjectUpdate!` | implemented | root field used by local GraphQL operation |
| `projectUpdates` | `ProjectUpdateConnection!` | implemented | root field used by local GraphQL operation |
| `projects` | `ProjectConnection!` | implemented | root field used by local GraphQL operation |
| `pushSubscriptionTest` | `PushSubscriptionTestPayload!` | intentionally_excluded | push subscription diagnostics are notification-device integration plumbing outside the CLI surface |
| `rateLimitStatus` | `RateLimitPayload!` | implemented | root field used by local GraphQL operation |
| `recentReleasesByAccessKey` | `[Release!]!` | intentionally_excluded | access-key release reads are unauthenticated sharing surfaces outside the token-scoped agent CLI |
| `release` | `Release!` | implemented | root field used by local GraphQL operation |
| `releaseNote` | `ReleaseNote!` | implemented | root field used by local GraphQL operation |
| `releaseNotes` | `ReleaseNoteConnection!` | implemented | root field used by local GraphQL operation |
| `releasePipeline` | `ReleasePipeline!` | implemented | root field used by local GraphQL operation |
| `releasePipelineByAccessKey` | `ReleasePipeline!` | intentionally_excluded | access-key release reads are unauthenticated sharing surfaces outside the token-scoped agent CLI |
| `releasePipelines` | `ReleasePipelineConnection!` | implemented | root field used by local GraphQL operation |
| `releaseSearch` | `[Release!]!` | implemented | root field used by local GraphQL operation |
| `releaseStage` | `ReleaseStage!` | implemented | root field used by local GraphQL operation |
| `releaseStages` | `ReleaseStageConnection!` | implemented | root field used by local GraphQL operation |
| `releases` | `ReleaseConnection!` | implemented | root field used by local GraphQL operation |
| `roadmap` | `Roadmap!` | implemented | root field used by local GraphQL operation |
| `roadmapToProject` | `RoadmapToProject!` | implemented | root field used by local GraphQL operation |
| `roadmapToProjects` | `RoadmapToProjectConnection!` | implemented | root field used by local GraphQL operation |
| `roadmaps` | `RoadmapConnection!` | implemented | root field used by local GraphQL operation |
| `searchDocuments` | `DocumentSearchPayload!` | implemented | root field used by local GraphQL operation |
| `searchIssues` | `IssueSearchPayload!` | implemented | root field used by local GraphQL operation |
| `searchProjects` | `ProjectSearchPayload!` | implemented | root field used by local GraphQL operation |
| `semanticSearch` | `SemanticSearchPayload!` | implemented | root field used by local GraphQL operation |
| `slaConfigurations` | `[SlaConfiguration!]!` | implemented | root field used by local GraphQL operation |
| `ssoUrlFromEmail` | `SsoUrlFromEmailResponse!` | intentionally_excluded | SSO discovery from email belongs to auth flow tooling, not the Linear work CLI |
| `team` | `Team!` | implemented | root field used by local GraphQL operation |
| `teamMembership` | `TeamMembership!` | implemented | root field used by local GraphQL operation |
| `teamMemberships` | `TeamMembershipConnection!` | implemented | root field used by local GraphQL operation |
| `teams` | `TeamConnection!` | implemented | root field used by local GraphQL operation |
| `template` | `Template!` | implemented | root field used by local GraphQL operation |
| `templates` | `[Template!]!` | implemented | root field used by local GraphQL operation |
| `templatesForIntegration` | `[Template!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeSchedule` | `TimeSchedule!` | implemented | root field used by local GraphQL operation |
| `timeSchedules` | `TimeScheduleConnection!` | implemented | root field used by local GraphQL operation |
| `triageResponsibilities` | `TriageResponsibilityConnection!` | implemented | root field used by local GraphQL operation |
| `triageResponsibility` | `TriageResponsibility!` | implemented | root field used by local GraphQL operation |
| `user` | `User!` | implemented | root field used by local GraphQL operation |
| `userSessions` | `[AuthenticationSessionResponse!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `userSettings` | `UserSettings!` | implemented | root field used by local GraphQL operation |
| `users` | `UserConnection!` | implemented | root field used by local GraphQL operation |
| `verifyGitHubEnterpriseServerInstallation` | `GitHubEnterpriseServerInstallVerificationPayload!` | intentionally_excluded | GitHub Enterprise installation verification is integration administration outside the CLI surface |
| `viewer` | `User!` | implemented | root field used by local GraphQL operation |
| `webhook` | `Webhook!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `webhooks` | `WebhookConnection!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `workflowState` | `WorkflowState!` | implemented | root field used by local GraphQL operation |
| `workflowStates` | `WorkflowStateConnection!` | implemented | root field used by local GraphQL operation |

## Upstream Mutation Root Fields

| Field | Return type | Status | Evidence |
| --- | --- | --- | --- |
| `agentActivityCreate` | `AgentActivityPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `agentActivityCreatePrompt` | `AgentActivityPayload!` | blocked_needs_design | mutation needs product and safety design |
| `agentActivityDeleteQueued` | `AgentActivityPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `agentActivitySendQueued` | `AgentActivityPayload!` | blocked_needs_design | mutation needs product and safety design |
| `agentSessionCreate` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnComment` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnIssue` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionUpdate` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionUpdateExternalUrl` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSkillCreate` | `AgentSkillPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `agentSkillDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `agentSkillUpdate` | `AgentSkillPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `airbyteIntegrationConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `attachmentCreate` | `AttachmentPayload!` | implemented | root field used by local GraphQL operation |
| `attachmentDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `attachmentLinkDiscord` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkFront` | `FrontAttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkGitHubIssue` | `AttachmentPayload!` | blocked_needs_design | attachment-to-GitHub linking mutates third-party integration state; needs explicit integration guard semantics |
| `attachmentLinkGitHubPR` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkGitLabMR` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkIntercom` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkJiraIssue` | `AttachmentPayload!` | blocked_needs_design | attachment-to-Jira linking mutates third-party integration state; needs explicit integration guard semantics |
| `attachmentLinkSalesforce` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkSlack` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkURL` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkZendesk` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentSyncToSlack` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentUpdate` | `AttachmentPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `commentCreate` | `CommentPayload!` | implemented | root field used by local GraphQL operation |
| `commentDelete` | `DeletePayload!` | implemented | root field used by local GraphQL operation |
| `commentResolve` | `CommentPayload!` | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `commentUnresolve` | `CommentPayload!` | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `commentUpdate` | `CommentPayload!` | implemented | root field used by local GraphQL operation |
| `contactCreate` | `ContactPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `contactSalesCreate` | `ContactPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCsvExportReport` | `CreateCsvExportReportPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createInitiativeUpdateReminder` | `InitiativeUpdateReminderPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createOrganizationFromOnboarding` | `CreateOrJoinOrganizationResponse!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProjectUpdateReminder` | `ProjectUpdateReminderPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customViewCreate` | `CustomViewPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customViewDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customViewUpdate` | `CustomViewPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerCreate` | `CustomerPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customerMerge` | `CustomerPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerNeedArchive` | `CustomerNeedArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerNeedCreate` | `CustomerNeedPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerNeedCreateFromAttachment` | `CustomerNeedPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerNeedDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customerNeedUnarchive` | `CustomerNeedArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerNeedUpdate` | `CustomerNeedUpdatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerStatusCreate` | `CustomerStatusPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerStatusDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customerStatusUpdate` | `CustomerStatusPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerTierCreate` | `CustomerTierPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerTierDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customerTierUpdate` | `CustomerTierPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerUnsync` | `CustomerPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerUpdate` | `CustomerPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customerUpsert` | `CustomerPayload!` | blocked_needs_design | mutation needs product and safety design |
| `cycleArchive` | `CycleArchivePayload!` | implemented | root field used by local GraphQL operation |
| `cycleCreate` | `CyclePayload!` | implemented | root field used by local GraphQL operation |
| `cycleShiftAll` | `CyclePayload!` | blocked_needs_design | bulk Cycle date shifting is a state-changing organization operation that needs target-pinned guard semantics |
| `cycleStartUpcomingCycleToday` | `CyclePayload!` | blocked_needs_design | starting an upcoming Cycle changes team planning state and needs target-pinned guard semantics |
| `cycleUpdate` | `CyclePayload!` | implemented | root field used by local GraphQL operation |
| `documentCreate` | `DocumentPayload!` | implemented | root field used by local GraphQL operation |
| `documentDelete` | `DocumentArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `documentUnarchive` | `DocumentArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `documentUpdate` | `DocumentPayload!` | implemented | root field used by local GraphQL operation |
| `emailIntakeAddressCreate` | `EmailIntakeAddressPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `emailIntakeAddressDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `emailIntakeAddressRefreshSesDomainStatus` | `EmailIntakeAddressRefreshSesDomainStatusPayload!` | blocked_needs_design | mutation needs product and safety design |
| `emailIntakeAddressRotate` | `EmailIntakeAddressPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `emailIntakeAddressUpdate` | `EmailIntakeAddressPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `emailTokenUserAccountAuth` | `AuthResolverResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emailUnsubscribe` | `EmailUnsubscribePayload!` | blocked_needs_design | mutation needs product and safety design |
| `emailUserAccountAuthChallenge` | `EmailUserAccountAuthChallengeResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emojiCreate` | `EmojiPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `emojiDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `entityExternalLinkCreate` | `EntityExternalLinkPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `entityExternalLinkDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `entityExternalLinkUpdate` | `EntityExternalLinkPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `favoriteCreate` | `FavoritePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `favoriteDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `favoriteUpdate` | `FavoritePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `fileUpload` | `UploadPayload!` | implemented | root field used by local GraphQL operation |
| `fileUploadDangerouslyDelete` | `FileUploadDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `gitAutomationStateCreate` | `GitAutomationStatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `gitAutomationStateDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `gitAutomationStateUpdate` | `GitAutomationStatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `gitAutomationTargetBranchCreate` | `GitAutomationTargetBranchPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `gitAutomationTargetBranchDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `gitAutomationTargetBranchUpdate` | `GitAutomationTargetBranchPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `googleUserAccountAuth` | `AuthResolverResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `imageUploadFromUrl` | `ImageUploadFromUrlPayload!` | blocked_needs_design | mutation needs product and safety design |
| `importFileUpload` | `UploadPayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeAddLabel` | `InitiativePayload!` | blocked_needs_design | initiative label mutation needs initiative target pinning and target-mismatch tests |
| `initiativeArchive` | `InitiativeArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeCreate` | `InitiativePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `initiativeRelationCreate` | `InitiativeRelationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeRelationDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `initiativeRelationUpdate` | `InitiativeRelationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeRemoveLabel` | `InitiativePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `initiativeToProjectCreate` | `InitiativeToProjectPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeToProjectDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `initiativeToProjectUpdate` | `InitiativeToProjectPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeUnarchive` | `InitiativeArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeUpdate` | `InitiativePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeUpdateArchive` | `InitiativeUpdateArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeUpdateCreate` | `InitiativeUpdatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeUpdateUnarchive` | `InitiativeUpdateArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `initiativeUpdateUpdate` | `InitiativeUpdatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `integrationArchive` | `DeletePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationAsksConnectChannel` | `AsksChannelConnectPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationCustomerDataAttributesRefresh` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `integrationDiscord` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationFigma` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationFront` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGitHubEnterpriseServerConnect` | `GitHubEnterpriseServerPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGitHubPersonal` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubCommitCreate` | `GitHubCommitIntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubImportConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubImportRefresh` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGithubRemoveCodeAccess` | `IntegrationGithubRemoveCodeAccessPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `integrationGitlabConnect` | `GitLabIntegrationCreatePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGitlabTestConnection` | `GitLabTestConnectionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGong` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGoogleCalendarPersonalConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationGoogleSheets` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationIntercom` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationIntercomDelete` | `IntegrationPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `integrationIntercomSettingsUpdate` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationJiraFetchProjectStatuses` | `JiraFetchProjectStatusesPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationJiraPersonal` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationJiraUpdate` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationLaunchDarklyConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationLaunchDarklyPersonalConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationLoom` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationMcpServerConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationMcpServerPersonalConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationMicrosoftPersonalConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationMicrosoftTeams` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationMicrosoftTeamsProjectPost` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationOpsgenieConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationOpsgenieRefreshScheduleMappings` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationPagerDutyConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationPagerDutyRefreshScheduleMappings` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationRequest` | `IntegrationRequestPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSalesforce` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSalesforceMetadataRefresh` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSentryConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSettingsUpdate` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlack` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackAsks` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackCustomViewNotifications` | `SlackChannelConnectPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackCustomerChannelLink` | `SuccessPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackImportEmojis` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackInitiativePost` | `SlackChannelConnectPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackOrAsksUpdateSlackTeamName` | `IntegrationSlackWorkspaceNamePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackOrgInitiativeUpdatesPost` | `SlackChannelConnectPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackOrgProjectUpdatesPost` | `SlackChannelConnectPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackPersonal` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackPost` | `SlackChannelConnectPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackProjectPost` | `SlackChannelConnectPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationSlackWorkflowAccessUpdate` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationTemplateCreate` | `IntegrationTemplatePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationTemplateDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `integrationUpdate` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationZendesk` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationsSettingsCreate` | `IntegrationsSettingsPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `integrationsSettingsUpdate` | `IntegrationsSettingsPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `issueAddLabel` | `IssuePayload!` | blocked_needs_design | issue label mutation needs issue target pinning and target-mismatch tests |
| `issueArchive` | `IssueArchivePayload!` | implemented | root field used by local GraphQL operation |
| `issueBatchCreate` | `IssueBatchPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueBatchUpdate` | `IssueBatchPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueCreate` | `IssuePayload!` | implemented | root field used by local GraphQL operation |
| `issueDelete` | `IssueArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueDescriptionUpdateFromFront` | `IssuePayload!` | blocked_needs_design | Front-origin description updates mutate issue content through integration state; needs explicit integration guard semantics |
| `issueExternalSyncDisable` | `IssuePayload!` | blocked_needs_design | issue external-sync disable changes integration state and needs explicit integration guard semantics |
| `issueImportCreateAsana` | `IssueImportPayload!` | blocked_needs_design | Asana issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateCSVJira` | `IssueImportPayload!` | blocked_needs_design | CSV/Jira issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateClubhouse` | `IssueImportPayload!` | blocked_needs_design | Clubhouse issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateGithub` | `IssueImportPayload!` | blocked_needs_design | GitHub issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateJira` | `IssueImportPayload!` | blocked_needs_design | Jira issue import creation starts external import workflow state and needs explicit integration guard semantics |
| `issueImportCreateLinearV2` | `IssueImportPayload!` | blocked_needs_design | Linear v2 issue import creation starts import workflow state and needs explicit import guard semantics |
| `issueImportDelete` | `IssueImportDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueImportProcess` | `IssueImportPayload!` | blocked_needs_design | issue import processing advances external import workflow state and needs explicit integration guard semantics |
| `issueImportUpdate` | `IssueImportPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueLabelCreate` | `IssueLabelPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueLabelDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueLabelRestore` | `IssueLabelPayload!` | blocked_needs_design | issue label lifecycle restore needs explicit organization/admin safety semantics |
| `issueLabelRetire` | `IssueLabelPayload!` | blocked_needs_design | issue label lifecycle retire needs explicit organization/admin safety semantics |
| `issueLabelUpdate` | `IssueLabelPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueRelationCreate` | `IssueRelationPayload!` | implemented | root field used by local GraphQL operation |
| `issueRelationDelete` | `DeletePayload!` | implemented | root field used by local GraphQL operation |
| `issueRelationUpdate` | `IssueRelationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueReminder` | `IssuePayload!` | blocked_needs_design | issue reminder mutation changes notification state and needs target-pinned guard semantics |
| `issueRemoveLabel` | `IssuePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueShare` | `IssuePayload!` | blocked_needs_design | issue sharing changes access state and needs target-pinned guard semantics |
| `issueSubscribe` | `IssuePayload!` | blocked_needs_design | issue subscription changes notification state and needs target-pinned guard semantics |
| `issueToReleaseCreate` | `IssueToReleasePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueToReleaseDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueToReleaseDeleteByIssueAndRelease` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueUnarchive` | `IssueArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueUnshare` | `IssuePayload!` | blocked_needs_design | issue unsharing changes access state and needs target-pinned guard semantics |
| `issueUnsubscribe` | `IssuePayload!` | blocked_needs_design | issue unsubscribe changes notification state and needs target-pinned guard semantics |
| `issueUpdate` | `IssuePayload!` | implemented | root field used by local GraphQL operation |
| `jiraIntegrationConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `joinOrganizationFromOnboarding` | `CreateOrJoinOrganizationResponse!` | blocked_needs_design | mutation needs product and safety design |
| `leaveOrganization` | `CreateOrJoinOrganizationResponse!` | blocked_needs_design | mutation needs product and safety design |
| `logout` | `LogoutResponse!` | blocked_needs_design | mutation needs product and safety design |
| `logoutAllSessions` | `LogoutResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `logoutOtherSessions` | `LogoutResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `logoutSession` | `LogoutResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `notificationArchive` | `NotificationArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `notificationArchiveAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationCategoryChannelSubscriptionUpdate` | `UserSettingsPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `notificationMarkReadAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationMarkUnreadAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationSnoozeAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationSubscriptionCreate` | `NotificationSubscriptionPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `notificationSubscriptionDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `notificationSubscriptionUpdate` | `NotificationSubscriptionPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `notificationUnarchive` | `NotificationArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `notificationUnsnoozeAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationUpdate` | `NotificationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `oauthApplicationArchive` | `OAuthApplicationArchivePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplicationCreate` | `OAuthApplicationCreatePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplicationRotateSecret` | `OAuthApplicationRotateSecretPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplicationRotateWebhookSecret` | `OAuthApplicationRotateWebhookSecretPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplicationUpdate` | `OAuthApplicationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `organizationCancelDelete` | `OrganizationCancelDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationDelete` | `OrganizationDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationDeleteChallenge` | `OrganizationDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationDomainClaim` | `OrganizationDomainSimplePayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationDomainCreate` | `OrganizationDomainPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `organizationDomainDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationDomainUpdate` | `OrganizationDomainPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `organizationDomainVerify` | `OrganizationDomainPayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationInviteCreate` | `OrganizationInvitePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `organizationInviteDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationInviteUpdate` | `OrganizationInvitePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `organizationStartTrial` | `OrganizationStartTrialPayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationStartTrialForPlan` | `OrganizationStartTrialPayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationUpdate` | `OrganizationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `passkeyLoginFinish` | `AuthResolverResponse!` | blocked_needs_design | mutation needs product and safety design |
| `passkeyLoginStart` | `PasskeyLoginStartResponse!` | blocked_needs_design | mutation needs product and safety design |
| `projectAddLabel` | `ProjectPayload!` | blocked_needs_design | project label mutation needs project target pinning and target-mismatch tests |
| `projectArchive` | `ProjectArchivePayload!` | implemented | root field used by local GraphQL operation |
| `projectCreate` | `ProjectPayload!` | implemented | root field used by local GraphQL operation |
| `projectCreateSlackChannel` | `ProjectPayload!` | blocked_needs_design | project Slack channel creation mutates chat integration state and needs explicit integration guard semantics |
| `projectDelete` | `ProjectArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectExternalSyncDisable` | `ProjectPayload!` | blocked_needs_design | project external-sync disable changes integration state and needs explicit integration guard semantics |
| `projectLabelCreate` | `ProjectLabelPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectLabelDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectLabelRestore` | `ProjectLabelPayload!` | blocked_needs_design | project label lifecycle restore needs explicit organization/admin safety semantics |
| `projectLabelRetire` | `ProjectLabelPayload!` | blocked_needs_design | project label lifecycle retire needs explicit organization/admin safety semantics |
| `projectLabelUpdate` | `ProjectLabelPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectMilestoneCreate` | `ProjectMilestonePayload!` | implemented | root field used by local GraphQL operation |
| `projectMilestoneDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectMilestoneMove` | `ProjectMilestoneMovePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectMilestoneUpdate` | `ProjectMilestonePayload!` | implemented | root field used by local GraphQL operation |
| `projectReassignStatus` | `SuccessPayload!` | blocked_needs_design | project status reassignment mutates project workflow state and needs target-pinned guard semantics |
| `projectRelationCreate` | `ProjectRelationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectRelationDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectRelationUpdate` | `ProjectRelationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectRemoveLabel` | `ProjectPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectStatusArchive` | `ProjectStatusArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectStatusCreate` | `ProjectStatusPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectStatusUnarchive` | `ProjectStatusArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectStatusUpdate` | `ProjectStatusPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectUnarchive` | `ProjectArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectUpdate` | `ProjectPayload!` | implemented | root field used by local GraphQL operation |
| `projectUpdateArchive` | `ProjectUpdateArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectUpdateCreate` | `ProjectUpdatePayload!` | implemented | root field used by local GraphQL operation |
| `projectUpdateDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectUpdateUnarchive` | `ProjectUpdateArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectUpdateUpdate` | `ProjectUpdatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `pushSubscriptionCreate` | `PushSubscriptionPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `pushSubscriptionDelete` | `PushSubscriptionPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `reactionCreate` | `ReactionPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `reactionDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `refreshGoogleSheetsData` | `IntegrationPayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseArchive` | `ReleaseArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseComplete` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseCompleteByAccessKey` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseCreate` | `ReleasePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseDelete` | `ReleaseArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `releaseNoteCreate` | `ReleaseNotePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseNoteDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `releaseNoteUpdate` | `ReleaseNotePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releasePipelineArchive` | `ReleasePipelineArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releasePipelineCreate` | `ReleasePipelinePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releasePipelineDelete` | `ReleasePipelineArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `releasePipelineUnarchive` | `ReleasePipelineArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releasePipelineUpdate` | `ReleasePipelinePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseStageArchive` | `ReleaseStageArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseStageCreate` | `ReleaseStagePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseStageUnarchive` | `ReleaseStageArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseStageUpdate` | `ReleaseStagePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseSync` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseSyncByAccessKey` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseUnarchive` | `ReleaseArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseUpdate` | `ReleasePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `releaseUpdateByPipeline` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseUpdateByPipelineByAccessKey` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `resendOrganizationInvite` | `DeletePayload!` | blocked_needs_design | mutation needs product and safety design |
| `resendOrganizationInviteByEmail` | `DeletePayload!` | blocked_needs_design | mutation needs product and safety design |
| `roadmapArchive` | `RoadmapArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `roadmapCreate` | `RoadmapPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `roadmapDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `roadmapToProjectCreate` | `RoadmapToProjectPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `roadmapToProjectDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `roadmapToProjectUpdate` | `RoadmapToProjectPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `roadmapUnarchive` | `RoadmapArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `roadmapUpdate` | `RoadmapPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `samlTokenUserAccountAuth` | `AuthResolverResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `teamCreate` | `TeamPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `teamCyclesDelete` | `TeamPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `teamDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `teamKeyDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `teamMembershipCreate` | `TeamMembershipPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `teamMembershipDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `teamMembershipUpdate` | `TeamMembershipPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `teamUnarchive` | `TeamArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `teamUpdate` | `TeamPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `templateCreate` | `TemplatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `templateDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `templateUpdate` | `TemplatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `timeScheduleCreate` | `TimeSchedulePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `timeScheduleDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `timeScheduleRefreshIntegrationSchedule` | `TimeSchedulePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeScheduleUpdate` | `TimeSchedulePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `timeScheduleUpsertExternal` | `TimeSchedulePayload!` | blocked_needs_design | mutation needs product and safety design |
| `trackAnonymousEvent` | `EventTrackingPayload!` | blocked_needs_design | mutation needs product and safety design |
| `triageResponsibilityCreate` | `TriageResponsibilityPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `triageResponsibilityDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `triageResponsibilityUpdate` | `TriageResponsibilityPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateIntegrationSlackScopes` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `userChangeRole` | `UserAdminPayload!` | intentionally_excluded | user role changes are organization administration outside the ordinary agent CLI surface |
| `userDiscordConnect` | `UserPayload!` | intentionally_excluded | Discord account connection belongs to user auth/integration setup, not work CLI reads |
| `userExternalUserDisconnect` | `UserPayload!` | intentionally_excluded | external-user disconnection is identity integration administration outside the ordinary agent CLI surface |
| `userFlagUpdate` | `UserSettingsFlagPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `userRevokeAllSessions` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userRevokeSession` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userSettingsFlagsReset` | `UserSettingsFlagsResetPayload!` | intentionally_excluded | user settings flag reset is internal preference administration outside the ordinary agent CLI surface |
| `userSettingsUpdate` | `UserSettingsPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `userSuspend` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userUnlinkFromIdentityProvider` | `UserAdminPayload!` | intentionally_excluded | identity-provider unlinking is auth administration outside the ordinary agent CLI surface |
| `userUnsuspend` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userUpdate` | `UserPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `viewPreferencesCreate` | `ViewPreferencesPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `viewPreferencesDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `viewPreferencesUpdate` | `ViewPreferencesPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `webhookCreate` | `WebhookPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `webhookDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `webhookRotateSecret` | `WebhookRotateSecretPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `webhookUpdate` | `WebhookPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `workflowStateArchive` | `WorkflowStateArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `workflowStateCreate` | `WorkflowStatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `workflowStateUpdate` | `WorkflowStatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |

## Local Generated Go Operations

| Operation | Kind | Root fields | Status | Evidence |
| --- | --- | --- | --- | --- |
| `AttachmentLinkURL` | mutation | `attachmentCreate` | implemented | `internal/client/generated.go` |
| `CommentDelete` | mutation | `commentDelete` | implemented | `internal/client/generated.go` |
| `CommentUpdate` | mutation | `commentUpdate` | implemented | `internal/client/generated.go` |
| `CompletedWorkflowStates` | query | `workflowStates` | implemented | `internal/client/generated.go` |
| `CycleArchive` | mutation | `cycleArchive` | implemented | `internal/client/generated.go` |
| `CycleCreate` | mutation | `cycleCreate` | implemented | `internal/client/generated.go` |
| `CycleReport` | query | `cycle` | implemented | `internal/client/generated.go` |
| `CycleUpdate` | mutation | `cycleUpdate` | implemented | `internal/client/generated.go` |
| `DocumentCreate` | mutation | `documentCreate` | implemented | `internal/client/generated.go` |
| `DocumentUpdate` | mutation | `documentUpdate` | implemented | `internal/client/generated.go` |
| `Documents` | query | `documents` | implemented | `internal/client/generated.go` |
| `IssueArchive` | mutation | `issueArchive` | implemented | `internal/client/generated.go` |
| `IssueBlockedIssues` | query | `issue` | implemented | `internal/client/generated.go` |
| `IssueClose` | mutation | `issueUpdate` | implemented | `internal/client/generated.go` |
| `IssueCommentCreate` | mutation | `commentCreate` | implemented | `internal/client/generated.go` |
| `IssueCreate` | mutation | `issueCreate` | implemented | `internal/client/generated.go` |
| `IssueDependencies` | query | `issue` | implemented | `internal/client/generated.go` |
| `IssueLabels` | query | `issueLabels` | implemented | `internal/client/generated.go` |
| `IssueRelationCreate` | mutation | `issueRelationCreate` | implemented | `internal/client/generated.go` |
| `IssueRelationDelete` | mutation | `issueRelationDelete` | implemented | `internal/client/generated.go` |
| `IssueUpdate` | mutation | `issueUpdate` | implemented | `internal/client/generated.go` |
| `IssuesByTeam` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamAssignee` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamBlocks` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamCreatedAfter` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamCreatedBefore` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamCycle` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamHasBlockers` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamLabel` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamProject` | query | `issues` | implemented | `internal/client/generated.go` |
| `IssuesByTeamState` | query | `issues` | implemented | `internal/client/generated.go` |
| `NextIssuesByTeam` | query | `issues` | implemented | `internal/client/generated.go` |
| `Organization` | query | `organization` | implemented | `internal/client/generated.go` |
| `ProjectArchive` | mutation | `projectArchive` | implemented | `internal/client/generated.go` |
| `ProjectCreate` | mutation | `projectCreate` | implemented | `internal/client/generated.go` |
| `ProjectMilestoneCreate` | mutation | `projectMilestoneCreate` | implemented | `internal/client/generated.go` |
| `ProjectMilestoneUpdate` | mutation | `projectMilestoneUpdate` | implemented | `internal/client/generated.go` |
| `ProjectUpdate` | mutation | `projectUpdate` | implemented | `internal/client/generated.go` |
| `ProjectUpdateCreate` | mutation | `projectUpdateCreate` | implemented | `internal/client/generated.go` |
| `Projects` | query | `team` | implemented | `internal/client/generated.go` |
| `StartedWorkflowStates` | query | `workflowStates` | implemented | `internal/client/generated.go` |
| `TargetProject` | query | `project` | implemented | `internal/client/generated.go` |
| `Teams` | query | `teams` | implemented | `internal/client/generated.go` |
| `Viewer` | query | `viewer` | implemented | `internal/client/generated.go` |
| `WorkflowStatesByType` | query | `workflowStates` | implemented | `internal/client/generated.go` |
| `agentActivities` | query | `agentActivities` | implemented | `internal/client/generated.go` |
| `agentActivity` | query | `agentActivity` | implemented | `internal/client/generated.go` |
| `agentSession` | query | `agentSession` | implemented | `internal/client/generated.go` |
| `agentSessions` | query | `agentSessions` | implemented | `internal/client/generated.go` |
| `agentSkill` | query | `agentSkill` | implemented | `internal/client/generated.go` |
| `agentSkills` | query | `agentSkills` | implemented | `internal/client/generated.go` |
| `applicationInfo` | query | `applicationInfo` | implemented | `internal/client/generated.go` |
| `attachment` | query | `attachment` | implemented | `internal/client/generated.go` |
| `attachmentIssue` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_attachments` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_botActor` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_children` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_comments` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_documents` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_formerAttachments` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_formerNeeds` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_history` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_inverseRelations` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_labels` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_needs` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_relations` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_releases` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_sharedAccess` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_stateHistory` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachmentIssue_subscribers` | query | `attachmentIssue` | implemented | `internal/client/generated.go` |
| `attachments` | query | `attachments` | implemented | `internal/client/generated.go` |
| `attachmentsForURL` | query | `attachmentsForURL` | implemented | `internal/client/generated.go` |
| `auditEntryTypes` | query | `auditEntryTypes` | implemented | `internal/client/generated.go` |
| `comment` | query | `comment` | implemented | `internal/client/generated.go` |
| `comment_botActor` | query | `comment` | implemented | `internal/client/generated.go` |
| `comment_children` | query | `comment` | implemented | `internal/client/generated.go` |
| `comment_createdIssues` | query | `comment` | implemented | `internal/client/generated.go` |
| `comments` | query | `comments` | implemented | `internal/client/generated.go` |
| `customView` | query | `customView` | implemented | `internal/client/generated.go` |
| `customViewHasSubscribers` | query | `customViewHasSubscribers` | implemented | `internal/client/generated.go` |
| `customView_initiatives` | query | `customView` | implemented | `internal/client/generated.go` |
| `customView_issues` | query | `customView` | implemented | `internal/client/generated.go` |
| `customView_organizationViewPreferences` | query | `customView` | implemented | `internal/client/generated.go` |
| `customView_organizationViewPreferences_preferences` | query | `customView` | implemented | `internal/client/generated.go` |
| `customView_projects` | query | `customView` | implemented | `internal/client/generated.go` |
| `customView_userViewPreferences` | query | `customView` | implemented | `internal/client/generated.go` |
| `customView_userViewPreferences_preferences` | query | `customView` | implemented | `internal/client/generated.go` |
| `customView_viewPreferencesValues` | query | `customView` | implemented | `internal/client/generated.go` |
| `customViews` | query | `customViews` | implemented | `internal/client/generated.go` |
| `customer` | query | `customer` | implemented | `internal/client/generated.go` |
| `customerNeed` | query | `customerNeed` | implemented | `internal/client/generated.go` |
| `customerNeed_projectAttachment` | query | `customerNeed` | implemented | `internal/client/generated.go` |
| `customerNeeds` | query | `customerNeeds` | implemented | `internal/client/generated.go` |
| `customerStatus` | query | `customerStatus` | implemented | `internal/client/generated.go` |
| `customerStatuses` | query | `customerStatuses` | implemented | `internal/client/generated.go` |
| `customerTier` | query | `customerTier` | implemented | `internal/client/generated.go` |
| `customerTiers` | query | `customerTiers` | implemented | `internal/client/generated.go` |
| `customers` | query | `customers` | implemented | `internal/client/generated.go` |
| `cycle` | query | `cycle` | implemented | `internal/client/generated.go` |
| `cycle_issues` | query | `cycle` | implemented | `internal/client/generated.go` |
| `cycle_uncompletedIssuesUponClose` | query | `cycle` | implemented | `internal/client/generated.go` |
| `cycles` | query | `cycles` | implemented | `internal/client/generated.go` |
| `document` | query | `document` | implemented | `internal/client/generated.go` |
| `document_comments` | query | `document` | implemented | `internal/client/generated.go` |
| `emoji` | query | `emoji` | implemented | `internal/client/generated.go` |
| `emojis` | query | `emojis` | implemented | `internal/client/generated.go` |
| `entityExternalLink` | query | `entityExternalLink` | implemented | `internal/client/generated.go` |
| `externalUser` | query | `externalUser` | implemented | `internal/client/generated.go` |
| `externalUsers` | query | `externalUsers` | implemented | `internal/client/generated.go` |
| `favorite` | query | `favorite` | implemented | `internal/client/generated.go` |
| `favorite_children` | query | `favorite` | implemented | `internal/client/generated.go` |
| `favorites` | query | `favorites` | implemented | `internal/client/generated.go` |
| `fileUpload` | mutation | `fileUpload` | implemented | `internal/client/generated.go` |
| `initiative` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiativeRelation` | query | `initiativeRelation` | implemented | `internal/client/generated.go` |
| `initiativeRelations` | query | `initiativeRelations` | implemented | `internal/client/generated.go` |
| `initiativeToProject` | query | `initiativeToProject` | implemented | `internal/client/generated.go` |
| `initiativeToProjects` | query | `initiativeToProjects` | implemented | `internal/client/generated.go` |
| `initiativeUpdate` | query | `initiativeUpdate` | implemented | `internal/client/generated.go` |
| `initiativeUpdate_comments` | query | `initiativeUpdate` | implemented | `internal/client/generated.go` |
| `initiativeUpdates` | query | `initiativeUpdates` | implemented | `internal/client/generated.go` |
| `initiative_documents` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiative_history` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiative_initiativeUpdates` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiative_links` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiative_projects` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiative_subInitiatives` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiatives` | query | `initiatives` | implemented | `internal/client/generated.go` |
| `issue` | query | `issue` | implemented | `internal/client/generated.go` |
| `issueFigmaFileKeySearch` | query | `issueFigmaFileKeySearch` | implemented | `internal/client/generated.go` |
| `issueFilterSuggestion` | query | `issueFilterSuggestion` | implemented | `internal/client/generated.go` |
| `issueLabel` | query | `issueLabel` | implemented | `internal/client/generated.go` |
| `issueLabel_children` | query | `issueLabel` | implemented | `internal/client/generated.go` |
| `issueLabel_issues` | query | `issueLabel` | implemented | `internal/client/generated.go` |
| `issuePriorityValues` | query | `issuePriorityValues` | implemented | `internal/client/generated.go` |
| `issueRelation` | query | `issueRelation` | implemented | `internal/client/generated.go` |
| `issueRelations` | query | `issueRelations` | implemented | `internal/client/generated.go` |
| `issueSearch` | query | `issueSearch` | implemented | `internal/client/generated.go` |
| `issueTitleSuggestionFromCustomerRequest` | query | `issueTitleSuggestionFromCustomerRequest` | implemented | `internal/client/generated.go` |
| `issueToRelease` | query | `issueToRelease` | implemented | `internal/client/generated.go` |
| `issueToReleases` | query | `issueToReleases` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_attachments` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_botActor` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_children` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_comments` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_documents` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_formerAttachments` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_formerNeeds` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_history` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_inverseRelations` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_labels` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_needs` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_relations` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_releases` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_sharedAccess` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_stateHistory` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issueVcsBranchSearch_subscribers` | query | `issueVcsBranchSearch` | implemented | `internal/client/generated.go` |
| `issue_attachments` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_botActor` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_children` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_comments` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_documents` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_formerAttachments` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_formerNeeds` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_history` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_inverseRelations` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_labels` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_needs` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_relations` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_releases` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_sharedAccess` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_stateHistory` | query | `issue` | implemented | `internal/client/generated.go` |
| `issue_subscribers` | query | `issue` | implemented | `internal/client/generated.go` |
| `issues` | query | `issues` | implemented | `internal/client/generated.go` |
| `notification` | query | `notification` | implemented | `internal/client/generated.go` |
| `notificationSubscription` | query | `notificationSubscription` | implemented | `internal/client/generated.go` |
| `notificationSubscriptions` | query | `notificationSubscriptions` | implemented | `internal/client/generated.go` |
| `notifications` | query | `notifications` | implemented | `internal/client/generated.go` |
| `organizationExists` | query | `organizationExists` | implemented | `internal/client/generated.go` |
| `organization_labels` | query | `organization` | implemented | `internal/client/generated.go` |
| `organization_projectLabels` | query | `organization` | implemented | `internal/client/generated.go` |
| `organization_teams` | query | `organization` | implemented | `internal/client/generated.go` |
| `organization_templates` | query | `organization` | implemented | `internal/client/generated.go` |
| `organization_users` | query | `organization` | implemented | `internal/client/generated.go` |
| `project` | query | `project` | implemented | `internal/client/generated.go` |
| `projectFilterSuggestion` | query | `projectFilterSuggestion` | implemented | `internal/client/generated.go` |
| `projectLabel` | query | `projectLabel` | implemented | `internal/client/generated.go` |
| `projectLabel_children` | query | `projectLabel` | implemented | `internal/client/generated.go` |
| `projectLabel_projects` | query | `projectLabel` | implemented | `internal/client/generated.go` |
| `projectLabels` | query | `projectLabels` | implemented | `internal/client/generated.go` |
| `projectMilestone` | query | `projectMilestone` | implemented | `internal/client/generated.go` |
| `projectMilestone_issues` | query | `projectMilestone` | implemented | `internal/client/generated.go` |
| `projectMilestones` | query | `projectMilestones` | implemented | `internal/client/generated.go` |
| `projectRelation` | query | `projectRelation` | implemented | `internal/client/generated.go` |
| `projectRelations` | query | `projectRelations` | implemented | `internal/client/generated.go` |
| `projectStatus` | query | `projectStatus` | implemented | `internal/client/generated.go` |
| `projectStatusProjectCount` | query | `projectStatusProjectCount` | implemented | `internal/client/generated.go` |
| `projectStatuses` | query | `projectStatuses` | implemented | `internal/client/generated.go` |
| `projectUpdate` | query | `projectUpdate` | implemented | `internal/client/generated.go` |
| `projectUpdate_comments` | query | `projectUpdate` | implemented | `internal/client/generated.go` |
| `projectUpdates` | query | `projectUpdates` | implemented | `internal/client/generated.go` |
| `project_attachments` | query | `project` | implemented | `internal/client/generated.go` |
| `project_comments` | query | `project` | implemented | `internal/client/generated.go` |
| `project_documents` | query | `project` | implemented | `internal/client/generated.go` |
| `project_externalLinks` | query | `project` | implemented | `internal/client/generated.go` |
| `project_history` | query | `project` | implemented | `internal/client/generated.go` |
| `project_initiativeToProjects` | query | `project` | implemented | `internal/client/generated.go` |
| `project_initiatives` | query | `project` | implemented | `internal/client/generated.go` |
| `project_inverseRelations` | query | `project` | implemented | `internal/client/generated.go` |
| `project_issues` | query | `project` | implemented | `internal/client/generated.go` |
| `project_labels` | query | `project` | implemented | `internal/client/generated.go` |
| `project_members` | query | `project` | implemented | `internal/client/generated.go` |
| `project_needs` | query | `project` | implemented | `internal/client/generated.go` |
| `project_projectMilestones` | query | `project` | implemented | `internal/client/generated.go` |
| `project_projectUpdates` | query | `project` | implemented | `internal/client/generated.go` |
| `project_relations` | query | `project` | implemented | `internal/client/generated.go` |
| `project_teams` | query | `project` | implemented | `internal/client/generated.go` |
| `projects` | query | `projects` | implemented | `internal/client/generated.go` |
| `rateLimitStatus` | query | `rateLimitStatus` | implemented | `internal/client/generated.go` |
| `release` | query | `release` | implemented | `internal/client/generated.go` |
| `releaseNote` | query | `releaseNote` | implemented | `internal/client/generated.go` |
| `releaseNotes` | query | `releaseNotes` | implemented | `internal/client/generated.go` |
| `releasePipeline` | query | `releasePipeline` | implemented | `internal/client/generated.go` |
| `releasePipeline_releases` | query | `releasePipeline` | implemented | `internal/client/generated.go` |
| `releasePipeline_stages` | query | `releasePipeline` | implemented | `internal/client/generated.go` |
| `releasePipeline_teams` | query | `releasePipeline` | implemented | `internal/client/generated.go` |
| `releasePipelines` | query | `releasePipelines` | implemented | `internal/client/generated.go` |
| `releaseSearch` | query | `releaseSearch` | implemented | `internal/client/generated.go` |
| `releaseStage` | query | `releaseStage` | implemented | `internal/client/generated.go` |
| `releaseStage_releases` | query | `releaseStage` | implemented | `internal/client/generated.go` |
| `releaseStages` | query | `releaseStages` | implemented | `internal/client/generated.go` |
| `release_documents` | query | `release` | implemented | `internal/client/generated.go` |
| `release_history` | query | `release` | implemented | `internal/client/generated.go` |
| `release_issues` | query | `release` | implemented | `internal/client/generated.go` |
| `release_links` | query | `release` | implemented | `internal/client/generated.go` |
| `releases` | query | `releases` | implemented | `internal/client/generated.go` |
| `roadmap` | query | `roadmap` | implemented | `internal/client/generated.go` |
| `roadmapToProject` | query | `roadmapToProject` | implemented | `internal/client/generated.go` |
| `roadmapToProjects` | query | `roadmapToProjects` | implemented | `internal/client/generated.go` |
| `roadmap_projects` | query | `roadmap` | implemented | `internal/client/generated.go` |
| `roadmaps` | query | `roadmaps` | implemented | `internal/client/generated.go` |
| `searchDocuments` | query | `searchDocuments` | implemented | `internal/client/generated.go` |
| `searchIssues` | query | `searchIssues` | implemented | `internal/client/generated.go` |
| `searchProjects` | query | `searchProjects` | implemented | `internal/client/generated.go` |
| `semanticSearch` | query | `semanticSearch` | implemented | `internal/client/generated.go` |
| `slaConfigurations` | query | `slaConfigurations` | implemented | `internal/client/generated.go` |
| `team` | query | `team` | implemented | `internal/client/generated.go` |
| `teamEstimateConfig` | query | `team` | implemented | `internal/client/generated.go` |
| `teamMembership` | query | `teamMembership` | implemented | `internal/client/generated.go` |
| `teamMemberships` | query | `teamMemberships` | implemented | `internal/client/generated.go` |
| `team_cycles` | query | `team` | implemented | `internal/client/generated.go` |
| `team_gitAutomationStates` | query | `team` | implemented | `internal/client/generated.go` |
| `team_issues` | query | `team` | implemented | `internal/client/generated.go` |
| `team_labels` | query | `team` | implemented | `internal/client/generated.go` |
| `team_members` | query | `team` | implemented | `internal/client/generated.go` |
| `team_memberships` | query | `team` | implemented | `internal/client/generated.go` |
| `team_projects` | query | `team` | implemented | `internal/client/generated.go` |
| `team_releasePipelines` | query | `team` | implemented | `internal/client/generated.go` |
| `team_states` | query | `team` | implemented | `internal/client/generated.go` |
| `team_templates` | query | `team` | implemented | `internal/client/generated.go` |
| `template` | query | `template` | implemented | `internal/client/generated.go` |
| `templateContent` | query | `template` | implemented | `internal/client/generated.go` |
| `templates` | query | `templates` | implemented | `internal/client/generated.go` |
| `timeSchedule` | query | `timeSchedule` | implemented | `internal/client/generated.go` |
| `timeSchedules` | query | `timeSchedules` | implemented | `internal/client/generated.go` |
| `triageResponsibilities` | query | `triageResponsibilities` | implemented | `internal/client/generated.go` |
| `triageResponsibility` | query | `triageResponsibility` | implemented | `internal/client/generated.go` |
| `triageResponsibility_manualSelection` | query | `triageResponsibility` | implemented | `internal/client/generated.go` |
| `user` | query | `user` | implemented | `internal/client/generated.go` |
| `userSettings` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_appsAndIntegrations` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_assignments` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_billing` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_commentsAndReplies` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_customers` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_documentChanges` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_feed` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_mentions` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_postsAndUpdates` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_reactions` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_reminders` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_reviews` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_statusChanges` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_subscriptions` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_system` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationCategoryPreferences_triage` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationChannelPreferences` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile_schedule` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile_schedule_friday` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile_schedule_monday` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile_schedule_saturday` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile_schedule_sunday` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile_schedule_thursday` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile_schedule_tuesday` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_notificationDeliveryPreferences_mobile_schedule_wednesday` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_theme` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_theme_custom` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `userSettings_theme_custom_sidebar` | query | `userSettings` | implemented | `internal/client/generated.go` |
| `user_assignedIssues` | query | `user` | implemented | `internal/client/generated.go` |
| `user_createdIssues` | query | `user` | implemented | `internal/client/generated.go` |
| `user_delegatedIssues` | query | `user` | implemented | `internal/client/generated.go` |
| `user_teamMemberships` | query | `user` | implemented | `internal/client/generated.go` |
| `user_teams` | query | `user` | implemented | `internal/client/generated.go` |
| `users` | query | `users` | implemented | `internal/client/generated.go` |
| `viewer` | query | `viewer` | implemented | `internal/client/generated.go` |
| `viewer_assignedIssues` | query | `viewer` | implemented | `internal/client/generated.go` |
| `viewer_createdIssues` | query | `viewer` | implemented | `internal/client/generated.go` |
| `viewer_delegatedIssues` | query | `viewer` | implemented | `internal/client/generated.go` |
| `viewer_drafts` | query | `viewer` | implemented | `internal/client/generated.go` |
| `viewer_teamMemberships` | query | `viewer` | implemented | `internal/client/generated.go` |
| `viewer_teams` | query | `viewer` | implemented | `internal/client/generated.go` |
| `workflowState` | query | `workflowState` | implemented | `internal/client/generated.go` |
| `workflowState_issues` | query | `workflowState` | implemented | `internal/client/generated.go` |
| `workflowStates` | query | `workflowStates` | implemented | `internal/client/generated.go` |

## Repo Domain-Map Commands

| Domain | Command | Backing | Scope | Status | Evidence |
| --- | --- | --- | --- | --- | --- |
| Core target | `whoami` | `Query.viewer`, `User` | Reads the authenticated user. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `target` | `Query.organization`, `Query.teams`, `Query.team`, `Query.projects`, `Query.project` | Resolves the active token's organization, team, and optional project. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `doctor` | `Query.viewer`, `Query.teams`, `TargetProject` (`Query.project`) when `project_id` is pinned | Read-only health check for config load, token presence, and pinned-target confirmation. Does not print token values. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `application info` | `Query.applicationInfo` | Read-only public OAuth application metadata by client id. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `organization exists` | `Query.organizationExists` | Read-only URL-key existence check for organization lookup. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `organization labels` | `Organization.labels` via `Query.organization` | Read-only organization-level issue labels. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `organization project-labels` | `Organization.projectLabels` via `Query.organization` | Read-only organization-level project labels. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `organization teams` | `Organization.teams` via `Query.organization` | Read-only teams visible to the authenticated user. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `organization templates` | `Organization.templates` via `Query.organization` | Read-only organization-level templates. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `organization users` | `Organization.users` via `Query.organization` | Read-only active users visible to the authenticated user. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Core target | `rate-limit status` | `Query.rateLimitStatus` | Read-only quota status for the authenticated Linear client. | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| AgentActivity | `agent-activity list` | `Query.agentActivities` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| AgentActivity | `agent-activity get` | `Query.agentActivity` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| AgentActivity | `agent-activity create` | `Mutation.agentActivityCreate` | Blocked: create writes into an agent session and needs explicit session/comment guard semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| AgentActivity | `agent-activity update` | `Mutation.agentActivityUpdate` | Blocked: update must resolve the agent session and activity scope before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| AgentActivity | `agent-activity archive` | `Mutation.agentActivityArchive` | Blocked: destructive command needs explicit AgentActivity safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| AgentSkill | `agent-skill list` | `Query.agentSkills` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| AgentSkill | `agent-skill get` | `Query.agentSkill` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| AgentSkill | `agent-skill create` | `Mutation.agentSkillCreate` | Blocked: create can expose reusable agent instructions and needs explicit team/owner guard semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| AgentSkill | `agent-skill update` | `Mutation.agentSkillUpdate` | Blocked: update must resolve the AgentSkill's team and ownership scope before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| AgentSkill | `agent-skill archive` | `Mutation.agentSkillArchive` | Blocked: destructive command needs explicit AgentSkill safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| AgentSession | `agent-session list` | `Query.agentSessions` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| AgentSession | `agent-session get` | `Query.agentSession` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ExternalUser | `external-user list` | `Query.externalUsers` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ExternalUser | `external-user get` | `Query.externalUser` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| AuditEntry | `audit-entry types` | `Query.auditEntryTypes` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| AuditEntry | `audit-entry list` | `Query.auditEntries` | Blocked: audit log entries can expose actor, IP, country, and request metadata; needs an explicit admin/security output model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification list` | `Query.notifications` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Notification | `notification get` | `Query.notification` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Notification | `notification subscription list` | `Query.notificationSubscriptions` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Notification | `notification subscription get` | `Query.notificationSubscription` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Notification | `notification archive` | `Mutation.notificationArchive` | Blocked: mutates the authenticated user's inbox state; needs an explicit viewer-state safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification archive all` | `Mutation.notificationArchiveAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification update` | `Mutation.notificationUpdate` | Blocked: direct inbox-state mutation needs an explicit viewer-state safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification mark read all` | `Mutation.notificationMarkReadAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification mark unread all` | `Mutation.notificationMarkUnreadAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification snooze all` | `Mutation.notificationSnoozeAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification unsnooze all` | `Mutation.notificationUnsnoozeAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification category channel subscription update` | `Mutation.notificationCategoryChannelSubscriptionUpdate` | Blocked: viewer notification preference mutation needs an explicit consent model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification subscription create` | `Mutation.notificationSubscriptionCreate` | Blocked: subscription writes can target several entity types and need explicit target-resolution semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification subscription update` | `Mutation.notificationSubscriptionUpdate` | Blocked: update must resolve the subscription target before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Notification | `notification subscription delete` | `Mutation.notificationSubscriptionDelete` | Blocked: destructive viewer preference command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Release | `release-pipeline list` | `Query.releasePipelines` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-pipeline get` | `Query.releasePipeline` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-pipeline releases` | `ReleasePipeline.releases` via `Query.releasePipeline` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-pipeline stages` | `ReleasePipeline.stages` via `Query.releasePipeline` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-pipeline teams` | `ReleasePipeline.teams` via `Query.releasePipeline` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-stage list` | `Query.releaseStages` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-stage get` | `Query.releaseStage` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-stage releases` | `ReleaseStage.releases` via `Query.releaseStage` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release list` | `Query.releases` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release search` | `Query.releaseSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release get` | `Query.release` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release history` | `Release.history` via `Query.release` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release documents` | `Release.documents` via `Query.release` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release issues` | `Release.issues` via `Query.release` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release links` | `Release.links` via `Query.release` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `external-link get` | `Query.entityExternalLink` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-note list` | `Query.releaseNotes` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-note get` | `Query.releaseNote` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `issue-to-release list` | `Query.issueToReleases` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `issue-to-release get` | `Query.issueToRelease` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Release | `release-pipeline create` | `Mutation.releasePipelineCreate` | Blocked: pipeline configuration is team/admin release surface and needs explicit guard semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-pipeline update` | `Mutation.releasePipelineUpdate` | Blocked: update must resolve and compare associated teams before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-pipeline archive` | `Mutation.releasePipelineArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-pipeline unarchive` | `Mutation.releasePipelineUnarchive` | Blocked: restore command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-pipeline delete` | `Mutation.releasePipelineDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Release | `release-stage create` | `Mutation.releaseStageCreate` | Blocked: release workflow configuration needs explicit pipeline/team guard semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-stage update` | `Mutation.releaseStageUpdate` | Blocked: update must resolve the stage's pipeline and teams before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-stage archive` | `Mutation.releaseStageArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-stage unarchive` | `Mutation.releaseStageUnarchive` | Blocked: restore command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release create` | `Mutation.releaseCreate` | Blocked: create must resolve pipeline/team guard semantics before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release update` | `Mutation.releaseUpdate` | Blocked: update must resolve the release pipeline/stage and associated teams before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release archive` | `Mutation.releaseArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release unarchive` | `Mutation.releaseUnarchive` | Blocked: restore command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release delete` | `Mutation.releaseDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Release | `release complete` | `Mutation.releaseComplete`, `Mutation.releaseCompleteByAccessKey` | Blocked: lifecycle transition and access-key behavior need explicit guard semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release sync` | `Mutation.releaseSync`, `Mutation.releaseSyncByAccessKey` | Blocked: sync mutates release associations and needs explicit guard semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-note create` | `Mutation.releaseNoteCreate` | Blocked: create must resolve release pipeline and release range semantics before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-note update` | `Mutation.releaseNoteUpdate` | Blocked: update must resolve covered releases and pipeline before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-note archive` | `Mutation.releaseNoteArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `release-note delete` | `Mutation.releaseNoteDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Release | `issue-to-release create` | `Mutation.issueToReleaseCreate` | Blocked: association write must compare issue and release scope before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `issue-to-release update` | `Mutation.issueToReleaseUpdate` | Blocked: association update must compare issue and release scope before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Release | `issue-to-release delete` | `Mutation.issueToReleaseDelete` | Blocked: destructive association command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Issue | `issue list` | `Query.issues`, optionally filtered by `Issue.team.id`, `Issue.state.type` (`--state`, with `--status` as an alias; human state names are normalized to the schema state type before filtering), `Issue.project.id`, `Issue.assignee.id`, `Issue.labels.some.id`, `Issue.cycle.id`, `Issue.createdAt.gte` (`--created-after` / `--created-since`), `Issue.createdAt.lte`, `Issue.hasBlockedByRelations.eq`, or `Issue.hasBlockingRelations.eq`; `--blocked-by ISSUE` traverses `Issue.relations` with `IssueRelation.type == "blocks"` and returns matching `IssueRelation.relatedIssue`; `--all-teams` omits the team filter | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue search` | `Query.issues`, filtered by `Issue.searchableContent` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue figma-file-key-search` | `Query.issueFigmaFileKeySearch`; returns compact issue summaries for a Figma file key | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue priority-values` | `Query.issuePriorityValues` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue filter-suggestion` | `Query.issueFilterSuggestion`; returns the suggested filter JSON plus log id only | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue title-suggestion` | `Query.issueTitleSuggestionFromCustomerRequest`; returns the suggested title plus log id only | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue get` | `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue deps` | `Query.issue`, `Issue.parent`, `Issue.children`, `Issue.relations`, `Issue.inverseRelations`; `IssueRelation.type == "blocks"` separates blocked issues from blockers | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue attachments` | `Issue.attachments` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue bot-actor` | `Issue.botActor` via `Query.issue` | Read-only, bot metadata only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue children` | `Issue.children` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue documents` | `Issue.documents` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue former-attachments` | `Issue.formerAttachments` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue former-needs` | `Issue.formerNeeds` via `Query.issue`; returns customer-need metadata without body/content | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue history` | `Issue.history` via `Query.issue`; returns compact metadata only, not raw change payloads or content fields | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue inverse-relations` | `Issue.inverseRelations` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue labels` | `Issue.labels` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue needs` | `Issue.needs` via `Query.issue`; returns customer-need metadata without body/content | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue relations` | `Issue.relations` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue releases` | `Issue.releases` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue shared-access` | `Issue.sharedAccess` via `Query.issue`; omits shared user details and exposes only flags/counts/disallowed fields | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue state-history` | `Issue.stateHistory` via `Query.issue` | Read-only, workflow-state span metadata | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue subscribers` | `Issue.subscribers` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search get` | `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search attachments` | `Issue.attachments` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search bot-actor` | `Issue.botActor` via `Query.issueVcsBranchSearch` | Read-only, bot metadata only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search children` | `Issue.children` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search documents` | `Issue.documents` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search former-attachments` | `Issue.formerAttachments` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search comments` | `Issue.comments` via `Query.issueVcsBranchSearch`; returns comment metadata without body | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search former-needs` | `Issue.formerNeeds` via `Query.issueVcsBranchSearch`; returns customer-need metadata without body/content | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search history` | `Issue.history` via `Query.issueVcsBranchSearch`; returns compact metadata only, not raw change payloads or content fields | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search inverse-relations` | `Issue.inverseRelations` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search labels` | `Issue.labels` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search needs` | `Issue.needs` via `Query.issueVcsBranchSearch`; returns customer-need metadata without body/content | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search relations` | `Issue.relations` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search releases` | `Issue.releases` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search shared-access` | `Issue.sharedAccess` via `Query.issueVcsBranchSearch`; omits shared user details and exposes only flags/counts/disallowed fields | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search state-history` | `Issue.stateHistory` via `Query.issueVcsBranchSearch` | Read-only, workflow-state span metadata | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue vcs-branch-search subscribers` | `Issue.subscribers` via `Query.issueVcsBranchSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue id` | Current checkout issue identifier from git/jj context | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Issue | `issue title` | `Query.issue` after current checkout or explicit issue resolution | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue url` | `Query.issue` after current checkout or explicit issue resolution | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue open` | `Query.issue` resolves `Issue.url`, then the platform opener (`xdg-open`/`open`/`rundll32`) launches it with the URL as a discrete argv argument | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue export` | `Query.issue` (`GetIssueDetail`), `Issue.comments`, and `Issue.attachments` are assembled into a single markdown file (`<DIR>/<identifier>.md`) holding the metadata header, description, comments, and attachment URLs; capped at 250 comments/attachments with a stderr note when more pages exist | Read-only, writes only local files | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue import` | Reads a CSV or JSON file (format from extension), normalizes each row's state/priority, rejects any row whose `team` key ≠ the pinned `team_key`, then creates each issue via guarded `Mutation.issueCreate` (`CreateIssue`); `--dry-run` renders the normalized rows locally and performs no mutation | Team-scoped per row; each create re-runs the write guard; `--dry-run` writes nothing | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue bulk-export` | `Query.team`/`Team.issues` (`ListIssuesByTeam`) for the resolved team are written to a CSV or JSON file (format from extension), capped by `--limit` (default 250) | Read-only, writes only the local file | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue branch` | `Query.issue`, `Issue.branchName` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue pr` | `Query.issue`; emits a local `gh pr create` title/body plan without calling GitHub | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `next` | `Query.issues`, filtered by `Issue.team.id`, `Issue.state.type == "unstarted"`, and `Issue.hasBlockedByRelations.eq == false`; fetches `Issue.relations`, `Issue.priority`, and `Issue.createdAt`, then ranks by active unblock count, priority, and age. `--dry-run` prints the top candidate and writes nothing; without it the top candidate is started via guarded `Mutation.issueUpdate` (`StartIssue`); `--checkout` runs `git checkout -b <Issue.branchName>` before starting | `--dry-run` read-only; otherwise resource-scoped start of the picked issue | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `done` | Current checkout issue identifier, then `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue create` | `Mutation.issueCreate` with `IssueCreateInput.teamId`, optional `projectId`; `--description-file` is resolved locally before mutation; `--template` reads `Template.templateData` via `Query.template` (free read) and fills title/description defaults that explicit flags override; `--section NAME=VALUE` fills a markdown section locally before mutation; `--dry-run` renders the assembled draft locally and performs no mutation; `--state` (alias `--status`) normalizes a human state name to a schema state type and resolves `IssueCreateInput.stateId` via `Query.workflowStates` filtered by team + type; `--priority` normalizes human words (`urgent`/`high`/`medium`/`low`/`none`) or `0-4` to `IssueCreateInput.priority` | Team-scoped unless `projectId` is set; `--dry-run` writes nothing | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue update` | `Mutation.issueUpdate` with `IssueUpdateInput`; `--description-file` replaces description, while `--append` or `--append-file` first reads `Issue.description` and appends text before sending `description`; `--state` (alias `--status`) and `--priority` are normalized the same way as on `issue create`, with `stateId` resolved via `Query.workflowStates` filtered by the issue's team + type | Resource-scoped when a project target is involved | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue start` | `Query.viewer`, `Query.workflowStates` filtered to `started`, then `Mutation.issueUpdate` with `IssueUpdateInput.assigneeId` and `stateId` | Resource-scoped when a project target is involved | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue comment` | `Mutation.commentCreate`; `--body -` reads stdin and `--body-file` reads a local file before mutation | Resource-scoped to the issue's resolved team/project | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue reply` | `Mutation.commentCreate` with `CommentCreateInput.parentId`; `--body-file` reads a local file before mutation | Resource-scoped to the issue's resolved team/project | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue close` | `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue link` | `Mutation.attachmentCreate` with `AttachmentCreateInput.issueId` and `url` | Resource-scoped: resolve the issue through `requireIssue` and compare the pinned team/project before attaching | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Issue | `issue comments` | `Issue.comments` via `Query.issue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| IssueRelation | `issue-relation list` | `Query.issueRelations` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| IssueRelation | `issue-relation get` | `Query.issueRelation` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| IssueRelation | `issue relate` | `Mutation.issueRelationCreate` with `IssueRelationCreateInput` | Team-scoped on both endpoints: resolve each issue and compare the pinned team before linking; `--type blocks` is refused when it would close a direct cycle | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| IssueRelation | `issue unrelate` | `Mutation.issueRelationDelete` | Resolve the relation, then compare the pinned team for both linked issues before deleting | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| IssueRelation | `issue-relation update` | `Mutation.issueRelationUpdate` | Blocked: update must resolve and compare both issue endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Comment | `comment list` | `Query.comments` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Comment | `comment get` | `Query.comment` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Comment | `comment bot-actor` | `Comment.botActor` via `Query.comment` | Read-only, bot metadata only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Comment | `comment children` | `Comment.children` via `Query.comment` | Read-only, body-free metadata | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Comment | `comment created-issues` | `Comment.createdIssues` via `Query.comment` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Comment | `comment update` | `Mutation.commentUpdate` with `CommentUpdateInput` | Resolve the comment, then compare the pinned team through its parent issue; non-issue comments are refused | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Comment | `comment delete` | `Mutation.commentDelete` | Resolve the comment, then compare the pinned team through its parent issue before deleting; non-issue comments are refused | blocked_needs_design | destructive command needs explicit safety semantics |
| Comment | `comment resolve` | `Mutation.commentResolve` | Blocked: resolving must first identify and compare the parent issue/project/update/document scope | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Comment | `comment unresolve` | `Mutation.commentUnresolve` | Blocked: unresolving must first identify and compare the parent issue/project/update/document scope | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Project | `project list` | `Query.team`, `Team.projects` | Read-only, resolved-team scoped | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Project | `project all` | `Query.projects` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Project | `project get` | `Query.project` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Project | `project open` | `Query.project` resolves `Project.url`, then the platform opener (`xdg-open`/`open`/`rundll32`) launches it with the URL as a discrete argv argument | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Project | `project attachments` | `Project.attachments` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project documents` | `Project.documents` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project external-links` | `Project.externalLinks` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project history` | `Project.history` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project initiative-links` | `Project.initiativeToProjects` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project initiatives` | `Project.initiatives` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project inverse-relations` | `Project.inverseRelations` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project issues` | `Project.issues` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project comments` | `Project.comments` | Read-only, body-free metadata | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project labels` | `Project.labels` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project create` | `Mutation.projectCreate` with `ProjectCreateInput.teamIds` | Team-scoped | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Project | `project update` | `Mutation.projectUpdate` with `ProjectUpdateInput` | Resource-scoped, compare `project_id` | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Project | `project archive` | `Mutation.projectArchive` | Resource-scoped, compare `project_id` | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Project | `project members` | `Project.members` plus `Mutation.projectUpdate` with `ProjectUpdateInput.memberIds` | Read-only for list, resource-scoped for writes | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Project | `project needs` | `Project.needs` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project relations` | `Project.relations` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project teams` | `Project.teams` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project updates` | `Project.projectUpdates` | Read-only, body-free metadata | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Project | `project filter-suggestion` | `Query.projectFilterSuggestion` | Read-only suggestion payload | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectUpdate | `project-update list` | `Query.projectUpdates` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectUpdate | `project-update get` | `Query.projectUpdate` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectUpdate | `project-update comments` | `ProjectUpdate.comments` | Read-only, body-free metadata | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| ProjectUpdate | `project-update create` | `Mutation.projectUpdateCreate` with `ProjectUpdateCreateInput` | Resource-scoped, compare `project_id` (pinned project) and team ownership | blocked_needs_design | write command needs explicit target and safety semantics |
| ProjectUpdate | `project-update update` | `Mutation.projectUpdateUpdate` | Blocked: update must resolve and compare the owning project before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectUpdate | `project-update archive` | `Mutation.projectUpdateArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectStatus | `project-status list` | `Query.projectStatuses` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectStatus | `project-status get` | `Query.projectStatus` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectStatus | `project-status project-count` | `Query.projectStatusProjectCount` | Read-only count payload | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectStatus | `project-status create` | `Mutation.projectStatusCreate` | Blocked: organization project status configuration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectStatus | `project-status update` | `Mutation.projectStatusUpdate` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectStatus | `project-status archive` | `Mutation.projectStatusArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectStatus | `project-status unarchive` | `Mutation.projectStatusUnarchive` | Blocked: restore semantics need an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectLabel | `project-label list` | `Query.projectLabels` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectLabel | `project-label get` | `Query.projectLabel` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectLabel | `project-label children` | `ProjectLabel.children` via `Query.projectLabel` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectLabel | `project-label projects` | `ProjectLabel.projects` via `Query.projectLabel` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectLabel | `project-label create` | `Mutation.projectLabelCreate` | Blocked: organization label configuration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectLabel | `project-label update` | `Mutation.projectLabelUpdate` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectLabel | `project-label delete` | `Mutation.projectLabelDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| ProjectLabel | `project-label retire` | `Mutation.projectLabelRetire` | Blocked: lifecycle command needs explicit admin safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectLabel | `project-label restore` | `Mutation.projectLabelRestore` | Blocked: restore semantics need an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectRelation | `project-relation list` | `Query.projectRelations` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectRelation | `project-relation get` | `Query.projectRelation` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectRelation | `project-relation create` | `Mutation.projectRelationCreate` | Blocked: create must resolve and compare both project dependency endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectRelation | `project-relation update` | `Mutation.projectRelationUpdate` | Blocked: update must resolve and compare both project dependency endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| ProjectRelation | `project-relation delete` | `Mutation.projectRelationDelete` | Blocked: destructive command needs explicit project dependency safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Cycle | `cycle list` | `Query.cycles` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Cycle | `cycle get` | `Query.cycle` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Cycle | `cycle issues` | `Cycle.issues` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Cycle | `cycle uncompleted-issues` | `Cycle.uncompletedIssuesUponClose` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Cycle | `cycle create` | `Mutation.cycleCreate` | Team-scoped | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Cycle | `cycle update` | `Mutation.cycleUpdate` | Team-scoped | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Cycle | `cycle archive` | `Mutation.cycleArchive` | Team-scoped | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Sprint | `sprint current` | `Query.cycles` filtered to active/current cycles | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Sprint | `sprint report` | `Query.cycle` plus `Cycle.issues` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectMilestone | `project-milestone all` | `Query.projectMilestones` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectMilestone | `project-milestone list` | `Project.projectMilestones` via `Query.project` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectMilestone | `project-milestone get` | `Query.projectMilestone` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectMilestone | `project-milestone issues` | `ProjectMilestone.issues` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| ProjectMilestone | `project-milestone create` | `Mutation.projectMilestoneCreate` with `projectId` | Resource-scoped, compare `project_id` | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectMilestone | `project-milestone update` | `Mutation.projectMilestoneUpdate` | Resource-scoped, compare resolved project | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| ProjectMilestone | `project-milestone delete` | `Mutation.projectMilestoneDelete` | Resource-scoped, compare resolved project | blocked_needs_design | destructive command needs explicit safety semantics |
| Document | `document list` | `Query.documents` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Document | `document get` | `Query.document` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Document | `document comments` | `Document.comments` | Read-only, body-free metadata | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Document | `document create` | `Mutation.documentCreate` with `DocumentCreateInput.teamId` from the resolved team and optional `projectId` from the pinned project; `--content` (or `--content-file`, or `--content -` for stdin) | Team-scoped unless a `project_id` is pinned | blocked_needs_design | write command needs explicit target and safety semantics |
| Document | `document update` | `Mutation.documentUpdate`; resolves the existing document via `Query.document` and compares its `team` (and pinned `project`) before mutating | Resource-scoped, compare team and pinned project | blocked_needs_design | write command needs explicit target and safety semantics |
| Document | `document delete` | `Mutation.documentDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Label | `label list` | `Query.issueLabels` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Label | `label get` | `Query.issueLabel` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Label | `label children` | `IssueLabel.children` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Label | `label issues` | `IssueLabel.issues` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Label | `label create` | `Mutation.issueLabelCreate` with optional `teamId` | Blocked: optional team scope needs explicit org/team target behavior before writes | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Label | `label update` | `Mutation.issueLabelUpdate` | Blocked: update must resolve and compare the label's owning team before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Label | `label delete` | `Mutation.issueLabelDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Team | `team list` | `Query.teams` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Team | `team get` | `Query.team` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Team | `team create` | `Mutation.teamCreate` | Blocked: organization administration surface needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Team | `team update` | `Mutation.teamUpdate` | Blocked: team metadata writes need stronger authority checks than ordinary target comparison | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Team | `team delete` | `Mutation.teamDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Team | `team cycles` | `Team.cycles` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team issues` | `Team.issues` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team labels` | `Team.labels` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team members` | `Team.members` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team memberships` | `Team.memberships` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team projects` | `Team.projects` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team release-pipelines` | `Team.releasePipelines` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team states` | `Team.states` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team git-automation-states` | `Team.gitAutomationStates` | Read-only, rule/state/target-branch metadata only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team templates` | `Team.templates` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Team | `team-membership list` | `Query.teamMemberships` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Team | `team-membership get` | `Query.teamMembership` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Team | `team-membership create` | `Mutation.teamMembershipCreate` | Blocked: organization membership administration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Team | `team-membership update` | `Mutation.teamMembershipUpdate` | Blocked: update must resolve and compare the membership's team and organization before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Team | `team-membership delete` | `Mutation.teamMembershipDelete` | Blocked: destructive membership command needs explicit admin safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| User | `user list` | `Query.users` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user get` | `Query.user` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user me` | `Query.viewer` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user drafts` | `User.drafts` via `Query.viewer` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings get` | `Query.userSettings` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings notification-categories` | `Query.userSettings.notificationCategoryPreferences` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings notification-category CATEGORY` | `Query.userSettings.notificationCategoryPreferences.<category>` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings notification-channels` | `Query.userSettings.notificationChannelPreferences` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings notification-delivery` | `Query.userSettings.notificationDeliveryPreferences` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings mobile-delivery` | `Query.userSettings.notificationDeliveryPreferences.mobile` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings mobile-schedule` | `Query.userSettings.notificationDeliveryPreferences.mobile.schedule` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings mobile-schedule-day DAY` | `Query.userSettings.notificationDeliveryPreferences.mobile.schedule.<day>` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings theme` | `Query.userSettings.theme` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings custom-theme` | `Query.userSettings.theme.custom` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user settings custom-sidebar-theme` | `Query.userSettings.theme.custom.sidebar` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user assigned-issues` | `User.assignedIssues` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| User | `user created-issues` | `User.createdIssues` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| User | `user delegated-issues` | `User.delegatedIssues` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| User | `user team-memberships` | `User.teamMemberships` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| User | `user teams` | `User.teams` | Read-only | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| User | `user my-assigned-issues` | `User.assignedIssues` via `Query.viewer` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user my-created-issues` | `User.createdIssues` via `Query.viewer` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user my-delegated-issues` | `User.delegatedIssues` via `Query.viewer` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user my-team-memberships` | `User.teamMemberships` via `Query.viewer` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| User | `user my-teams` | `User.teams` via `Query.viewer` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| WorkflowState | `workflow-state list` | `Query.workflowStates` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| WorkflowState | `workflow-state get` | `Query.workflowState` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| WorkflowState | `workflow-state issues` | `WorkflowState.issues` via `Query.workflowState` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| WorkflowState | `workflow-state create` | `Mutation.workflowStateCreate` | Blocked: team workflow configuration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| WorkflowState | `workflow-state update` | `Mutation.workflowStateUpdate` | Blocked: update must resolve and compare the owning team before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| WorkflowState | `workflow-state archive` | `Mutation.workflowStateArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| TimeSchedule | `time-schedule list` | `Query.timeSchedules` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| TimeSchedule | `time-schedule get` | `Query.timeSchedule` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| TimeSchedule | `time-schedule create` | `Mutation.timeScheduleCreate` | Blocked: schedule create needs explicit owner/admin safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| TimeSchedule | `time-schedule update` | `Mutation.timeScheduleUpdate` | Blocked: update must resolve schedule scope before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| TimeSchedule | `time-schedule delete` | `Mutation.timeScheduleDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| TimeSchedule | `time-schedule upsert-external` | `Mutation.timeScheduleUpsertExternal` | Blocked: external integration sync surface is not an ordinary agent workflow | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| TriageResponsibility | `triage-responsibility list` | `Query.triageResponsibilities` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| TriageResponsibility | `triage-responsibility get` | `Query.triageResponsibility` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| TriageResponsibility | `triage-responsibility manual-selection` | `TriageResponsibility.manualSelection` via `Query.triageResponsibility` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| TriageResponsibility | `triage-responsibility create` | `Mutation.triageResponsibilityCreate` | Blocked: team triage configuration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| TriageResponsibility | `triage-responsibility update` | `Mutation.triageResponsibilityUpdate` | Blocked: update must resolve and compare the owning team before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| TriageResponsibility | `triage-responsibility delete` | `Mutation.triageResponsibilityDelete` | Blocked: destructive team triage configuration command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| SLA Configuration | `sla-configuration list` | `Query.slaConfigurations` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| SemanticSearch | `semantic-search` | `Query.semanticSearch` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Search | `search documents` | `Query.searchDocuments` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Search | `search issues` | `Query.searchIssues` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Search | `search projects` | `Query.searchProjects` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Template | `template list` | `Query.templates` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Template | `template get` | `Query.template` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Template | `template create` | `Mutation.templateCreate` | Blocked: create can be organization-, team-, or pipeline-scoped and needs explicit guard semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Template | `template update` | `Mutation.templateUpdate` | Blocked: update must resolve and compare the template's organization, team, or pipeline scope before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Template | `template delete` | `Mutation.templateDelete` | Blocked: destructive command needs explicit template-scope safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Initiative | `initiative list` | `Query.initiatives` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Initiative | `initiative get` | `Query.initiative` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Initiative | `initiative history` | `Initiative.history` via `Query.initiative` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Initiative | `initiative links` | `Initiative.links` via `Query.initiative` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Initiative | `initiative sub-initiatives` | `Initiative.subInitiatives` via `Query.initiative` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Initiative | `initiative updates` | `Initiative.initiativeUpdates` via `Query.initiative` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Initiative | `initiative documents` | `Initiative.documents` via `Query.initiative` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Initiative | `initiative projects` | `Initiative.projects` via `Query.initiative` | Read-only direct projects | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Initiative | `initiative create` | `Mutation.createInitiative` | Blocked: initiative create needs an explicit organization-scoped safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Initiative | `initiative update` | `Mutation.updateInitiative` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Initiative | `initiative archive` | `Mutation.archiveInitiative` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| InitiativeRelation | `initiative-relation list` | `Query.initiativeRelations` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| InitiativeRelation | `initiative-relation get` | `Query.initiativeRelation` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| InitiativeRelation | `initiative-relation create` | `Mutation.initiativeRelationCreate` | Blocked: create must resolve and compare both Initiative hierarchy endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| InitiativeRelation | `initiative-relation update` | `Mutation.initiativeRelationUpdate` | Blocked: update must resolve and compare both Initiative hierarchy endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| InitiativeRelation | `initiative-relation delete` | `Mutation.initiativeRelationDelete` | Blocked: destructive command needs explicit hierarchy safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| InitiativeToProject | `initiative-to-project list` | `Query.initiativeToProjects` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| InitiativeToProject | `initiative-to-project get` | `Query.initiativeToProject` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| InitiativeToProject | `initiative-to-project create` | `Mutation.initiativeToProjectCreate` | Blocked: create must resolve and compare both Initiative and Project endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| InitiativeToProject | `initiative-to-project update` | `Mutation.initiativeToProjectUpdate` | Blocked: update must resolve and compare both Initiative and Project endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| InitiativeToProject | `initiative-to-project delete` | `Mutation.initiativeToProjectDelete` | Blocked: destructive command needs explicit association safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| RoadmapToProject | `roadmap-to-project list` | `Query.roadmapToProjects` | Legacy read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| RoadmapToProject | `roadmap-to-project get` | `Query.roadmapToProject` | Legacy read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| RoadmapToProject | `roadmap-to-project create` | `Mutation.roadmapToProjectCreate` | Blocked: deprecated create must resolve and compare both Roadmap and Project endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| RoadmapToProject | `roadmap-to-project update` | `Mutation.roadmapToProjectUpdate` | Blocked: deprecated update must resolve and compare both Roadmap and Project endpoints before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| RoadmapToProject | `roadmap-to-project delete` | `Mutation.roadmapToProjectDelete` | Blocked: destructive deprecated association command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| InitiativeUpdate | `initiative-update list` | `Query.initiativeUpdates` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| InitiativeUpdate | `initiative-update get` | `Query.initiativeUpdate` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| InitiativeUpdate | `initiative-update comments` | `InitiativeUpdate.comments` | Read-only, body-free metadata | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| InitiativeUpdate | `initiative-update create` | `Mutation.initiativeUpdateCreate` | Blocked: create must resolve and compare the owning Initiative before posting | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| InitiativeUpdate | `initiative-update update` | `Mutation.initiativeUpdateUpdate` | Blocked: update must resolve and compare the owning Initiative before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| InitiativeUpdate | `initiative-update archive` | `Mutation.initiativeUpdateArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| InitiativeUpdate | `initiative-update unarchive` | `Mutation.initiativeUpdateUnarchive` | Blocked: unarchive needs explicit lifecycle and target semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Roadmap | `roadmap list` | `Query.roadmaps` | Legacy read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Roadmap | `roadmap get` | `Query.roadmap` | Legacy read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Roadmap | `roadmap projects` | `Roadmap.projects` via `Query.roadmap` | Legacy read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Roadmap | `roadmap create` | `Mutation.roadmapCreate` | Blocked: deprecated organization-scoped planning surface needs an explicit safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Roadmap | `roadmap update` | `Mutation.roadmapUpdate` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Roadmap | `roadmap archive` | `Mutation.roadmapArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Roadmap | `roadmap delete` | `Mutation.roadmapDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| CustomView | `custom-view list` | `Query.customViews` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| CustomView | `custom-view subscribers` | `Query.customViewHasSubscribers` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| CustomView | `custom-view get` | `Query.customView` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| CustomView | `custom-view initiatives` | `Query.customView_initiatives` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root |
| CustomView | `custom-view issues` | `Query.customView_issues` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root |
| CustomView | `custom-view organization-preferences` | `Query.customView_organizationViewPreferences` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root |
| CustomView | `custom-view organization-preferences values` | `Query.customView_organizationViewPreferences_preferences` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root |
| CustomView | `custom-view projects` | `Query.customView_projects` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root |
| CustomView | `custom-view user-preferences` | `Query.customView_userViewPreferences` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root |
| CustomView | `custom-view user-preferences values` | `Query.customView_userViewPreferences_preferences` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root |
| CustomView | `custom-view preference-values` | `Query.customView_viewPreferencesValues` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root |
| CustomView | `custom-view create` | `Mutation.createCustomView` | Blocked: custom view create needs an explicit organization-scoped safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| CustomView | `custom-view update` | `Mutation.updateCustomView` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| CustomView | `custom-view delete` | `Mutation.deleteCustomView` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Customer | `customer list` | `Query.customers` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer get` | `Query.customer` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer-need list` | `Query.customerNeeds` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer-need get` | `Query.customerNeed` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer-need project-attachment` | `CustomerNeed.projectAttachment` via `Query.customerNeed` | Read-only, metadata-only projection | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer-status list` | `Query.customerStatuses` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer-status get` | `Query.customerStatus` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer-tier list` | `Query.customerTiers` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer-tier get` | `Query.customerTier` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Customer | `customer create` | `Mutation.customerCreate` | Blocked: customer create needs an explicit organization-scoped safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer update` | `Mutation.customerUpdate` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer archive` | `Mutation.customerArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer-need create` | `Mutation.customerNeedCreate` | Blocked: need creation must prove the linked issue, project, or customer target before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer-need update` | `Mutation.customerNeedUpdate` | Blocked: update must resolve the need and compare the linked issue or project target before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer-need archive` | `Mutation.customerNeedArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer-need delete` | `Mutation.customerNeedDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Customer | `customer-status create` | `Mutation.customerStatusCreate` | Blocked: organization lifecycle configuration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer-status update` | `Mutation.customerStatusUpdate` | Blocked: organization lifecycle configuration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer-status delete` | `Mutation.customerStatusDelete` | Blocked: destructive admin command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Customer | `customer-tier create` | `Mutation.customerTierCreate` | Blocked: organization tier configuration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer-tier update` | `Mutation.customerTierUpdate` | Blocked: organization tier configuration needs an explicit admin safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Customer | `customer-tier delete` | `Mutation.customerTierDelete` | Blocked: destructive admin command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Favorite | `favorite list` | `Query.favorites` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Favorite | `favorite children` | `Favorite.children` via `Query.favorite` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Favorite | `favorite get` | `Query.favorite` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Favorite | `favorite create` | `Mutation.createFavorite` | Blocked: favorite create needs an explicit viewer-scoped safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Favorite | `favorite update` | `Mutation.updateFavorite` | Blocked: update must resolve and compare the owning viewer before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Favorite | `favorite delete` | `Mutation.deleteFavorite` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Emoji | `emoji list` | `Query.emojis` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Emoji | `emoji get` | `Query.emoji` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Emoji | `emoji create` | `Mutation.createEmoji` | Blocked: emoji create needs an explicit organization-scoped safety model | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Emoji | `emoji delete` | `Mutation.deleteEmoji` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| File | `files upload` | `Mutation.fileUpload` then an HTTP PUT of the bytes to the pre-signed URL | Raw Linear asset, not target-pinned; prints the asset URL for a later guarded attachment write | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| File | `files download` | Plain HTTP GET of the asset URL to a local path | Read-only, no API; no auth header is attached so a user-supplied URL never receives the Linear token | implemented | `linctl --help` / public CLI tests; no direct GraphQL root in backing |
| Attachment | `attachment list` | `Query.attachments` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment url` | `Query.attachmentsForURL` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment get` | `Query.attachment` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue get` | `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue attachments` | `Issue.attachments` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue bot-actor` | `Issue.botActor` via `Query.attachmentIssue` | Read-only, bot metadata only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue children` | `Issue.children` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue comments` | `Issue.comments` via `Query.attachmentIssue`; returns comment metadata without body | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue documents` | `Issue.documents` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue former-attachments` | `Issue.formerAttachments` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue former-needs` | `Issue.formerNeeds` via `Query.attachmentIssue`; returns customer-need metadata without body/content | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue history` | `Issue.history` via `Query.attachmentIssue` | Read-only, compact metadata only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue inverse-relations` | `Issue.inverseRelations` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue labels` | `Issue.labels` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue needs` | `Issue.needs` via `Query.attachmentIssue`; returns customer-need metadata without body/content | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue relations` | `Issue.relations` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue releases` | `Issue.releases` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue shared-access` | `Issue.sharedAccess` via `Query.attachmentIssue`; omits shared user details and exposes only flags/counts/disallowed fields | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue state-history` | `Issue.stateHistory` via `Query.attachmentIssue` | Read-only, workflow-state span metadata | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment issue subscribers` | `Issue.subscribers` via `Query.attachmentIssue` | Read-only | implemented | `linctl --help`, `docs/domain-map.md`, and local GraphQL root |
| Attachment | `attachment create` | `Mutation.attachmentCreate` | Blocked: attachment create must resolve and compare the owning issue's team before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Attachment | `attachment update` | `Mutation.attachmentUpdate` | Blocked: update must resolve and compare the owning issue before mutation | blocked_needs_design | blocked in `docs/domain-map.md` pending explicit safety semantics |
| Attachment | `attachment delete` | `Mutation.attachmentDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
