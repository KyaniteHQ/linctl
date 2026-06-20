# Linear API coverage ledger

Generated from current local sources and upstream Linear SDK commit `df20561`.

Sources:

- Upstream SDK methods: `/tmp/linear-sdk-source/packages/sdk/src/_generated_sdk.ts`
- Upstream schema roots: `/tmp/linear-sdk-source/packages/sdk/src/schema.graphql`
- Local generated operations: `internal/client/generated.go`
- Local GraphQL operations: `internal/client/operations/*.graphql`
- Repo domain map: `docs/domain-map.md`

Statuses: `implemented`, `accepted_gap`, `safe_candidate`, `blocked_needs_design`, `intentionally_excluded`.

## Summary

| Surface | Total | Implemented/root-backed | Classified |
| --- | ---: | ---: | ---: |
| Upstream SDK root methods | 458 | 91 | 458 |
| Upstream Query root fields | 158 | 79 | 158 |
| Upstream Mutation root fields | 364 | 12 | 364 |
| Local generated Go operations | 142 | 142 | 142 |
| Domain-map commands | 233 | 122 | 233 |

## Upstream SDK Root Methods

| Method | Kind | Status | Evidence |
| --- | --- | --- | --- |
| `administrableTeams` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentActivities` | method | implemented | local operation or command exists |
| `agentActivity` | method | implemented | local operation or command exists |
| `agentSession` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnComment` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnIssue` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionUpdateExternalUrl` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessions` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
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
| `attachmentIssue` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `attachmentLinkDiscord` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentLinkFront` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentLinkGitHubIssue` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `attachmentLinkGitHubPR` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentLinkGitLabMR` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentLinkIntercom` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentLinkJiraIssue` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `attachmentLinkSalesforce` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentLinkSlack` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentLinkURL` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentLinkZendesk` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentSyncToSlack` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachments` | method | implemented | local operation or command exists |
| `attachmentsForURL` | method | implemented | local operation or command exists |
| `auditEntries` | method | safe_candidate | read operation may fit future CLI coverage |
| `auditEntryTypes` | getter | implemented | local operation or command exists |
| `authenticationSessions` | getter | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `availableUsers` | getter | accepted_gap | repo-planned or likely useful CLI domain |
| `comment` | method | implemented | local operation or command exists |
| `commentResolve` | method | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `commentUnresolve` | method | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `comments` | method | implemented | local operation or command exists |
| `constructor` | method | safe_candidate | read operation may fit future CLI coverage |
| `createAgentActivity` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createAgentSkill` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createAttachment` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createComment` | method | implemented | local operation or command exists |
| `createContact` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCsvExportReport` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomView` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomer` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomerNeed` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomerStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCustomerTier` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createCycle` | method | implemented | local operation or command exists |
| `createDocument` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
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
| `createIssueRelation` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createIssueToRelease` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createNotificationSubscription` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createOrganizationInvite` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProject` | method | implemented | local operation or command exists |
| `createProjectLabel` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProjectMilestone` | method | implemented | local operation or command exists |
| `createProjectRelation` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProjectStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProjectUpdate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
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
| `customerMerge` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerNeed` | method | implemented | local operation or command exists |
| `customerNeedCreateFromAttachment` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerNeeds` | method | implemented | local operation or command exists |
| `customerStatus` | method | implemented | local operation or command exists |
| `customerStatuses` | method | implemented | local operation or command exists |
| `customerTier` | method | implemented | local operation or command exists |
| `customerTiers` | method | implemented | local operation or command exists |
| `customerUnsync` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerUpsert` | method | safe_candidate | read operation may fit future CLI coverage |
| `customers` | method | implemented | local operation or command exists |
| `cycle` | method | implemented | local operation or command exists |
| `cycleShiftAll` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `cycleStartUpcomingCycleToday` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `cycles` | method | implemented | local operation or command exists |
| `deleteAgentSkill` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteAttachment` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `deleteComment` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
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
| `deleteIssueRelation` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
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
| `documentContentHistory` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `documents` | method | implemented | local operation or command exists |
| `emailIntakeAddress` | method | safe_candidate | read operation may fit future CLI coverage |
| `emailIntakeAddressRefreshSesDomainStatus` | method | safe_candidate | read operation may fit future CLI coverage |
| `emailIntakeAddressRotate` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `emailTokenUserAccountAuth` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emailUnsubscribe` | method | safe_candidate | read operation may fit future CLI coverage |
| `emailUserAccountAuthChallenge` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emoji` | method | implemented | local operation or command exists |
| `emojis` | method | implemented | local operation or command exists |
| `entityExternalLink` | method | implemented | local operation or command exists |
| `externalUser` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `externalUsers` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `favorite` | method | implemented | local operation or command exists |
| `favorites` | method | implemented | local operation or command exists |
| `fileUpload` | method | safe_candidate | read operation may fit future CLI coverage |
| `googleUserAccountAuth` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `imageUploadFromUrl` | method | safe_candidate | read operation may fit future CLI coverage |
| `importFileUpload` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `issueAddLabel` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueExternalSyncDisable` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueFigmaFileKeySearch` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueFilterSuggestion` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCheckCSV` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCheckSync` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateAsana` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateCSVJira` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateClubhouse` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateGithub` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateJira` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportJqlCheck` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportProcess` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabel` | method | implemented | local operation or command exists |
| `issueLabelRestore` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabelRetire` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabels` | method | implemented | local operation or command exists |
| `issuePriorityValues` | getter | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRelation` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRelations` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueReminder` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRemoveLabel` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueRepositorySuggestions` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueSearch` | method | implemented | local operation or command exists |
| `issueShare` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueSubscribe` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueTitleSuggestionFromCustomerRequest` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToRelease` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToReleaseDeleteByIssueAndRelease` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueToReleases` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueUnshare` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueUnsubscribe` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issueVcsBranchSearch` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `issues` | method | implemented | local operation or command exists |
| `latestReleaseByAccessKey` | getter | safe_candidate | read operation may fit future CLI coverage |
| `logout` | method | safe_candidate | read operation may fit future CLI coverage |
| `logoutAllSessions` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `logoutOtherSessions` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `logoutSession` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `notification` | method | implemented | local operation or command exists |
| `notificationArchiveAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationMarkReadAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationMarkUnreadAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationSnoozeAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationSubscription` | method | implemented | local operation or command exists |
| `notificationSubscriptions` | method | implemented | local operation or command exists |
| `notificationUnsnoozeAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notifications` | method | implemented | local operation or command exists |
| `organization` | getter | implemented | local operation or command exists |
| `organizationDeleteChallenge` | getter | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationExists` | method | implemented | local operation or command exists |
| `organizationInvite` | method | safe_candidate | read operation may fit future CLI coverage |
| `organizationInvites` | method | safe_candidate | read operation may fit future CLI coverage |
| `organizationStartTrial` | getter | safe_candidate | read operation may fit future CLI coverage |
| `organizationStartTrialForPlan` | method | safe_candidate | read operation may fit future CLI coverage |
| `project` | method | implemented | local operation or command exists |
| `projectAddLabel` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectExternalSyncDisable` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectFilterSuggestion` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabel` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelRestore` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelRetire` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabels` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectMilestone` | method | implemented | local operation or command exists |
| `projectMilestones` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectRelation` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectRelations` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectRemoveLabel` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectStatus` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectStatuses` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `projectUpdate` | method | implemented | local operation or command exists |
| `projectUpdates` | method | implemented | local operation or command exists |
| `projects` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `pushSubscriptionTest` | method | safe_candidate | read operation may fit future CLI coverage |
| `rateLimitStatus` | getter | implemented | local operation or command exists |
| `recentReleasesByAccessKey` | method | safe_candidate | read operation may fit future CLI coverage |
| `refreshGoogleSheetsData` | method | safe_candidate | read operation may fit future CLI coverage |
| `release` | method | implemented | local operation or command exists |
| `releaseComplete` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseCompleteByAccessKey` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseNote` | method | implemented | local operation or command exists |
| `releaseNotes` | method | implemented | local operation or command exists |
| `releasePipeline` | method | implemented | local operation or command exists |
| `releasePipelineByAccessKey` | getter | safe_candidate | read operation may fit future CLI coverage |
| `releasePipelines` | method | implemented | local operation or command exists |
| `releaseSearch` | method | implemented | local operation or command exists |
| `releaseStage` | method | implemented | local operation or command exists |
| `releaseStages` | method | implemented | local operation or command exists |
| `releaseSync` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseSyncByAccessKey` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseUpdateByPipeline` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseUpdateByPipelineByAccessKey` | method | safe_candidate | read operation may fit future CLI coverage |
| `releases` | method | implemented | local operation or command exists |
| `resendOrganizationInvite` | method | safe_candidate | read operation may fit future CLI coverage |
| `resendOrganizationInviteByEmail` | method | safe_candidate | read operation may fit future CLI coverage |
| `roadmap` | method | implemented | local operation or command exists |
| `roadmapToProject` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmapToProjects` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmaps` | method | implemented | local operation or command exists |
| `rotateSecretWebhook` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `samlTokenUserAccountAuth` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `searchDocuments` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `searchIssues` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `searchProjects` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `semanticSearch` | method | safe_candidate | read operation may fit future CLI coverage |
| `slaConfigurations` | method | safe_candidate | read operation may fit future CLI coverage |
| `ssoUrlFromEmail` | method | safe_candidate | read operation may fit future CLI coverage |
| `suspendUser` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `team` | method | implemented | local operation or command exists |
| `teamMembership` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `teamMemberships` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `teams` | method | implemented | local operation or command exists |
| `template` | method | implemented | local operation or command exists |
| `templates` | getter | implemented | local operation or command exists |
| `templatesForIntegration` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeSchedule` | method | implemented | local operation or command exists |
| `timeScheduleRefreshIntegrationSchedule` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeScheduleUpsertExternal` | method | safe_candidate | read operation may fit future CLI coverage |
| `timeSchedules` | method | implemented | local operation or command exists |
| `trackAnonymousEvent` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `unsuspendUser` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `updateAgentSession` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `updateAgentSkill` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateAttachment` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateComment` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomView` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomer` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomerNeed` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomerStatus` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCustomerTier` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `updateCycle` | method | implemented | local operation or command exists |
| `updateDocument` | method | blocked_needs_design | write operation needs guarded target semantics before exposure |
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
| `userChangeRole` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `userDiscordConnect` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `userExternalUserDisconnect` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `userRevokeAllSessions` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userRevokeSession` | method | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userSessions` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `userSettings` | getter | accepted_gap | repo-planned or likely useful CLI domain |
| `userSettingsFlagsReset` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `userUnlinkFromIdentityProvider` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `users` | method | implemented | local operation or command exists |
| `verifyGitHubEnterpriseServerInstallation` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `agentSession` | `AgentSession!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionSandbox` | `CodingAgentSandboxPayload` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessions` | `AgentSessionConnection!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSkill` | `AgentSkill!` | implemented | root field used by local GraphQL operation |
| `agentSkills` | `AgentSkillConnection!` | implemented | root field used by local GraphQL operation |
| `applicationInfo` | `Application!` | implemented | root field used by local GraphQL operation |
| `archivedIntegrations` | `[ArchivedIntegrationPayload!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `archivedTeams` | `[Team!]!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `attachment` | `Attachment!` | implemented | root field used by local GraphQL operation |
| `attachmentIssue` | `Issue!` | accepted_gap | repo-planned or likely useful CLI domain |
| `attachmentSources` | `AttachmentSourcesPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `attachments` | `AttachmentConnection!` | implemented | root field used by local GraphQL operation |
| `attachmentsForURL` | `AttachmentConnection!` | implemented | root field used by local GraphQL operation |
| `auditEntries` | `AuditEntryConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `auditEntryTypes` | `[AuditEntryType!]!` | implemented | root field used by local GraphQL operation |
| `authenticationSessions` | `[AuthenticationSessionResponse!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `availableUsers` | `AuthResolverResponse!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `documentContentHistory` | `DocumentContentHistoryPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `documentContentHistoryEntries` | `DocumentContentHistoryPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `documentContentHistoryTimeline` | `DocumentContentHistoryTimelinePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `documents` | `DocumentConnection!` | implemented | root field used by local GraphQL operation |
| `emailIntakeAddress` | `EmailIntakeAddress!` | safe_candidate | read operation may fit future CLI coverage |
| `emoji` | `Emoji!` | implemented | root field used by local GraphQL operation |
| `emojis` | `EmojiConnection!` | implemented | root field used by local GraphQL operation |
| `entityExternalLink` | `EntityExternalLink!` | implemented | root field used by local GraphQL operation |
| `externalUser` | `ExternalUser!` | accepted_gap | repo-planned or likely useful CLI domain |
| `externalUsers` | `ExternalUserConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `issueFigmaFileKeySearch` | `IssueConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueFilterSuggestion` | `IssueFilterSuggestionPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCheckCSV` | `IssueImportCheckPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCheckSync` | `IssueImportSyncCheckPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportJqlCheck` | `IssueImportJqlCheckPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabel` | `IssueLabel!` | implemented | root field used by local GraphQL operation |
| `issueLabels` | `IssueLabelConnection!` | implemented | root field used by local GraphQL operation |
| `issuePriorityValues` | `[IssuePriorityValue!]!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRelation` | `IssueRelation!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRelations` | `IssueRelationConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRepositorySuggestions` | `RepositorySuggestionsPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueSearch` | `IssueConnection!` | implemented | root field used by local GraphQL operation |
| `issueTitleSuggestionFromCustomerRequest` | `IssueTitleSuggestionFromCustomerRequestPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToRelease` | `IssueToRelease!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToReleases` | `IssueToReleaseConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueVcsBranchSearch` | `Issue` | accepted_gap | repo-planned or likely useful CLI domain |
| `issues` | `IssueConnection!` | implemented | root field used by local GraphQL operation |
| `latestReleaseByAccessKey` | `Release` | safe_candidate | read operation may fit future CLI coverage |
| `microsoftTeamsChannels` | `MicrosoftTeamsChannelsPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `notification` | `Notification!` | implemented | root field used by local GraphQL operation |
| `notificationSubscription` | `NotificationSubscription!` | implemented | root field used by local GraphQL operation |
| `notificationSubscriptions` | `NotificationSubscriptionConnection!` | implemented | root field used by local GraphQL operation |
| `notifications` | `NotificationConnection!` | implemented | root field used by local GraphQL operation |
| `notificationsUnreadCount` | `Int!` | safe_candidate | read operation may fit future CLI coverage |
| `oauthApplication` | `OAuthApplication!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplications` | `[OAuthApplication!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `organization` | `Organization!` | implemented | root field used by local GraphQL operation |
| `organizationDomainClaimRequest` | `OrganizationDomainClaimPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `organizationExists` | `OrganizationExistsPayload!` | implemented | root field used by local GraphQL operation |
| `organizationInvite` | `OrganizationInvite!` | safe_candidate | read operation may fit future CLI coverage |
| `organizationInviteDetails` | `OrganizationInviteDetailsPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `organizationInvites` | `OrganizationInviteConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `organizationMeta` | `OrganizationMeta` | safe_candidate | read operation may fit future CLI coverage |
| `project` | `Project!` | implemented | root field used by local GraphQL operation |
| `projectFilterSuggestion` | `ProjectFilterSuggestionPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabel` | `ProjectLabel!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabels` | `ProjectLabelConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectMilestone` | `ProjectMilestone!` | implemented | root field used by local GraphQL operation |
| `projectMilestones` | `ProjectMilestoneConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectRelation` | `ProjectRelation!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectRelations` | `ProjectRelationConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectStatus` | `ProjectStatus!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectStatusProjectCount` | `ProjectStatusCountPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectStatuses` | `ProjectStatusConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectUpdate` | `ProjectUpdate!` | implemented | root field used by local GraphQL operation |
| `projectUpdates` | `ProjectUpdateConnection!` | implemented | root field used by local GraphQL operation |
| `projects` | `ProjectConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `pushSubscriptionTest` | `PushSubscriptionTestPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `rateLimitStatus` | `RateLimitPayload!` | implemented | root field used by local GraphQL operation |
| `recentReleasesByAccessKey` | `[Release!]!` | safe_candidate | read operation may fit future CLI coverage |
| `release` | `Release!` | implemented | root field used by local GraphQL operation |
| `releaseNote` | `ReleaseNote!` | implemented | root field used by local GraphQL operation |
| `releaseNotes` | `ReleaseNoteConnection!` | implemented | root field used by local GraphQL operation |
| `releasePipeline` | `ReleasePipeline!` | implemented | root field used by local GraphQL operation |
| `releasePipelineByAccessKey` | `ReleasePipeline!` | safe_candidate | read operation may fit future CLI coverage |
| `releasePipelines` | `ReleasePipelineConnection!` | implemented | root field used by local GraphQL operation |
| `releaseSearch` | `[Release!]!` | implemented | root field used by local GraphQL operation |
| `releaseStage` | `ReleaseStage!` | implemented | root field used by local GraphQL operation |
| `releaseStages` | `ReleaseStageConnection!` | implemented | root field used by local GraphQL operation |
| `releases` | `ReleaseConnection!` | implemented | root field used by local GraphQL operation |
| `roadmap` | `Roadmap!` | implemented | root field used by local GraphQL operation |
| `roadmapToProject` | `RoadmapToProject!` | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmapToProjects` | `RoadmapToProjectConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmaps` | `RoadmapConnection!` | implemented | root field used by local GraphQL operation |
| `searchDocuments` | `DocumentSearchPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `searchIssues` | `IssueSearchPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `searchProjects` | `ProjectSearchPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `semanticSearch` | `SemanticSearchPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `slaConfigurations` | `[SlaConfiguration!]!` | safe_candidate | read operation may fit future CLI coverage |
| `ssoUrlFromEmail` | `SsoUrlFromEmailResponse!` | safe_candidate | read operation may fit future CLI coverage |
| `team` | `Team!` | implemented | root field used by local GraphQL operation |
| `teamMembership` | `TeamMembership!` | accepted_gap | repo-planned or likely useful CLI domain |
| `teamMemberships` | `TeamMembershipConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `userSettings` | `UserSettings!` | accepted_gap | repo-planned or likely useful CLI domain |
| `users` | `UserConnection!` | implemented | root field used by local GraphQL operation |
| `verifyGitHubEnterpriseServerInstallation` | `GitHubEnterpriseServerInstallVerificationPayload!` | safe_candidate | read operation may fit future CLI coverage |
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
| `attachmentCreate` | `AttachmentPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `attachmentDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `attachmentLinkDiscord` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkFront` | `FrontAttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkGitHubIssue` | `AttachmentPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `attachmentLinkGitHubPR` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkGitLabMR` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkIntercom` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkJiraIssue` | `AttachmentPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `attachmentLinkSalesforce` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkSlack` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkURL` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentLinkZendesk` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentSyncToSlack` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `attachmentUpdate` | `AttachmentPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `commentCreate` | `CommentPayload!` | implemented | root field used by local GraphQL operation |
| `commentDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `commentResolve` | `CommentPayload!` | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `commentUnresolve` | `CommentPayload!` | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `commentUpdate` | `CommentPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
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
| `cycleShiftAll` | `CyclePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `cycleStartUpcomingCycleToday` | `CyclePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `cycleUpdate` | `CyclePayload!` | implemented | root field used by local GraphQL operation |
| `documentCreate` | `DocumentPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `documentDelete` | `DocumentArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `documentUnarchive` | `DocumentArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `documentUpdate` | `DocumentPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
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
| `fileUpload` | `UploadPayload!` | blocked_needs_design | mutation needs product and safety design |
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
| `initiativeAddLabel` | `InitiativePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `issueAddLabel` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueArchive` | `IssueArchivePayload!` | implemented | root field used by local GraphQL operation |
| `issueBatchCreate` | `IssueBatchPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueBatchUpdate` | `IssueBatchPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueCreate` | `IssuePayload!` | implemented | root field used by local GraphQL operation |
| `issueDelete` | `IssueArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueDescriptionUpdateFromFront` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueExternalSyncDisable` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateAsana` | `IssueImportPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateCSVJira` | `IssueImportPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateClubhouse` | `IssueImportPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateGithub` | `IssueImportPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateJira` | `IssueImportPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportCreateLinearV2` | `IssueImportPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportDelete` | `IssueImportDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueImportProcess` | `IssueImportPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueImportUpdate` | `IssueImportPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueLabelCreate` | `IssueLabelPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueLabelDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueLabelRestore` | `IssueLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabelRetire` | `IssueLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabelUpdate` | `IssueLabelPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueRelationCreate` | `IssueRelationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueRelationDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueRelationUpdate` | `IssueRelationPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueReminder` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRemoveLabel` | `IssuePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueShare` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueSubscribe` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToReleaseCreate` | `IssueToReleasePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueToReleaseDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueToReleaseDeleteByIssueAndRelease` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueUnarchive` | `IssueArchivePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `issueUnshare` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueUnsubscribe` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `projectAddLabel` | `ProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectArchive` | `ProjectArchivePayload!` | implemented | root field used by local GraphQL operation |
| `projectCreate` | `ProjectPayload!` | implemented | root field used by local GraphQL operation |
| `projectCreateSlackChannel` | `ProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectDelete` | `ProjectArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectExternalSyncDisable` | `ProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelCreate` | `ProjectLabelPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectLabelDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectLabelRestore` | `ProjectLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelRetire` | `ProjectLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelUpdate` | `ProjectLabelPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectMilestoneCreate` | `ProjectMilestonePayload!` | implemented | root field used by local GraphQL operation |
| `projectMilestoneDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectMilestoneMove` | `ProjectMilestoneMovePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `projectMilestoneUpdate` | `ProjectMilestonePayload!` | implemented | root field used by local GraphQL operation |
| `projectReassignStatus` | `SuccessPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `projectUpdateCreate` | `ProjectUpdatePayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
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
| `userChangeRole` | `UserAdminPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userDiscordConnect` | `UserPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userExternalUserDisconnect` | `UserPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userFlagUpdate` | `UserSettingsFlagPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `userRevokeAllSessions` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userRevokeSession` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userSettingsFlagsReset` | `UserSettingsFlagsResetPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userSettingsUpdate` | `UserSettingsPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `userSuspend` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userUnlinkFromIdentityProvider` | `UserAdminPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `CompletedWorkflowStates` | query | `workflowStates` | implemented | `internal/client/generated.go` |
| `CycleArchive` | mutation | `cycleArchive` | implemented | `internal/client/generated.go` |
| `CycleCreate` | mutation | `cycleCreate` | implemented | `internal/client/generated.go` |
| `CycleReport` | query | `cycle` | implemented | `internal/client/generated.go` |
| `CycleUpdate` | mutation | `cycleUpdate` | implemented | `internal/client/generated.go` |
| `Documents` | query | `documents` | implemented | `internal/client/generated.go` |
| `IssueArchive` | mutation | `issueArchive` | implemented | `internal/client/generated.go` |
| `IssueBlockedIssues` | query | `issue` | implemented | `internal/client/generated.go` |
| `IssueClose` | mutation | `issueUpdate` | implemented | `internal/client/generated.go` |
| `IssueCommentCreate` | mutation | `commentCreate` | implemented | `internal/client/generated.go` |
| `IssueCreate` | mutation | `issueCreate` | implemented | `internal/client/generated.go` |
| `IssueDependencies` | query | `issue` | implemented | `internal/client/generated.go` |
| `IssueLabels` | query | `issueLabels` | implemented | `internal/client/generated.go` |
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
| `ProjectMilestones` | query | `project` | implemented | `internal/client/generated.go` |
| `ProjectUpdate` | mutation | `projectUpdate` | implemented | `internal/client/generated.go` |
| `ProjectUpdates` | query | `project` | implemented | `internal/client/generated.go` |
| `Projects` | query | `team` | implemented | `internal/client/generated.go` |
| `StartedWorkflowStates` | query | `workflowStates` | implemented | `internal/client/generated.go` |
| `TargetProject` | query | `project` | implemented | `internal/client/generated.go` |
| `Teams` | query | `teams` | implemented | `internal/client/generated.go` |
| `Viewer` | query | `viewer` | implemented | `internal/client/generated.go` |
| `agentActivities` | query | `agentActivities` | implemented | `internal/client/generated.go` |
| `agentActivity` | query | `agentActivity` | implemented | `internal/client/generated.go` |
| `agentSkill` | query | `agentSkill` | implemented | `internal/client/generated.go` |
| `agentSkills` | query | `agentSkills` | implemented | `internal/client/generated.go` |
| `applicationInfo` | query | `applicationInfo` | implemented | `internal/client/generated.go` |
| `attachment` | query | `attachment` | implemented | `internal/client/generated.go` |
| `attachments` | query | `attachments` | implemented | `internal/client/generated.go` |
| `attachmentsForURL` | query | `attachmentsForURL` | implemented | `internal/client/generated.go` |
| `auditEntryTypes` | query | `auditEntryTypes` | implemented | `internal/client/generated.go` |
| `comment` | query | `comment` | implemented | `internal/client/generated.go` |
| `comments` | query | `comments` | implemented | `internal/client/generated.go` |
| `customView` | query | `customView` | implemented | `internal/client/generated.go` |
| `customViewHasSubscribers` | query | `customViewHasSubscribers` | implemented | `internal/client/generated.go` |
| `customViews` | query | `customViews` | implemented | `internal/client/generated.go` |
| `customer` | query | `customer` | implemented | `internal/client/generated.go` |
| `customerNeed` | query | `customerNeed` | implemented | `internal/client/generated.go` |
| `customerNeeds` | query | `customerNeeds` | implemented | `internal/client/generated.go` |
| `customerStatus` | query | `customerStatus` | implemented | `internal/client/generated.go` |
| `customerStatuses` | query | `customerStatuses` | implemented | `internal/client/generated.go` |
| `customerTier` | query | `customerTier` | implemented | `internal/client/generated.go` |
| `customerTiers` | query | `customerTiers` | implemented | `internal/client/generated.go` |
| `customers` | query | `customers` | implemented | `internal/client/generated.go` |
| `cycle` | query | `cycle` | implemented | `internal/client/generated.go` |
| `cycles` | query | `cycles` | implemented | `internal/client/generated.go` |
| `document` | query | `document` | implemented | `internal/client/generated.go` |
| `emoji` | query | `emoji` | implemented | `internal/client/generated.go` |
| `emojis` | query | `emojis` | implemented | `internal/client/generated.go` |
| `entityExternalLink` | query | `entityExternalLink` | implemented | `internal/client/generated.go` |
| `favorite` | query | `favorite` | implemented | `internal/client/generated.go` |
| `favorite_children` | query | `favorite` | implemented | `internal/client/generated.go` |
| `favorites` | query | `favorites` | implemented | `internal/client/generated.go` |
| `initiative` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiativeRelation` | query | `initiativeRelation` | implemented | `internal/client/generated.go` |
| `initiativeRelations` | query | `initiativeRelations` | implemented | `internal/client/generated.go` |
| `initiativeToProject` | query | `initiativeToProject` | implemented | `internal/client/generated.go` |
| `initiativeToProjects` | query | `initiativeToProjects` | implemented | `internal/client/generated.go` |
| `initiativeUpdate` | query | `initiativeUpdate` | implemented | `internal/client/generated.go` |
| `initiativeUpdates` | query | `initiativeUpdates` | implemented | `internal/client/generated.go` |
| `initiative_history` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiative_initiativeUpdates` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiative_links` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiative_subInitiatives` | query | `initiative` | implemented | `internal/client/generated.go` |
| `initiatives` | query | `initiatives` | implemented | `internal/client/generated.go` |
| `issue` | query | `issue` | implemented | `internal/client/generated.go` |
| `issueLabel` | query | `issueLabel` | implemented | `internal/client/generated.go` |
| `issueSearch` | query | `issueSearch` | implemented | `internal/client/generated.go` |
| `issue_comments` | query | `issue` | implemented | `internal/client/generated.go` |
| `issues` | query | `issues` | implemented | `internal/client/generated.go` |
| `notification` | query | `notification` | implemented | `internal/client/generated.go` |
| `notificationSubscription` | query | `notificationSubscription` | implemented | `internal/client/generated.go` |
| `notificationSubscriptions` | query | `notificationSubscriptions` | implemented | `internal/client/generated.go` |
| `notifications` | query | `notifications` | implemented | `internal/client/generated.go` |
| `organizationExists` | query | `organizationExists` | implemented | `internal/client/generated.go` |
| `organization_templates` | query | `organization` | implemented | `internal/client/generated.go` |
| `project` | query | `project` | implemented | `internal/client/generated.go` |
| `projectMilestone` | query | `projectMilestone` | implemented | `internal/client/generated.go` |
| `projectUpdate` | query | `projectUpdate` | implemented | `internal/client/generated.go` |
| `projectUpdates` | query | `projectUpdates` | implemented | `internal/client/generated.go` |
| `project_members` | query | `project` | implemented | `internal/client/generated.go` |
| `rateLimitStatus` | query | `rateLimitStatus` | implemented | `internal/client/generated.go` |
| `release` | query | `release` | implemented | `internal/client/generated.go` |
| `releaseNote` | query | `releaseNote` | implemented | `internal/client/generated.go` |
| `releaseNotes` | query | `releaseNotes` | implemented | `internal/client/generated.go` |
| `releasePipeline` | query | `releasePipeline` | implemented | `internal/client/generated.go` |
| `releasePipeline_releases` | query | `releasePipeline` | implemented | `internal/client/generated.go` |
| `releasePipeline_stages` | query | `releasePipeline` | implemented | `internal/client/generated.go` |
| `releasePipelines` | query | `releasePipelines` | implemented | `internal/client/generated.go` |
| `releaseSearch` | query | `releaseSearch` | implemented | `internal/client/generated.go` |
| `releaseStage` | query | `releaseStage` | implemented | `internal/client/generated.go` |
| `releaseStage_releases` | query | `releaseStage` | implemented | `internal/client/generated.go` |
| `releaseStages` | query | `releaseStages` | implemented | `internal/client/generated.go` |
| `release_history` | query | `release` | implemented | `internal/client/generated.go` |
| `release_links` | query | `release` | implemented | `internal/client/generated.go` |
| `releases` | query | `releases` | implemented | `internal/client/generated.go` |
| `roadmap` | query | `roadmap` | implemented | `internal/client/generated.go` |
| `roadmaps` | query | `roadmaps` | implemented | `internal/client/generated.go` |
| `team` | query | `team` | implemented | `internal/client/generated.go` |
| `team_members` | query | `team` | implemented | `internal/client/generated.go` |
| `template` | query | `template` | implemented | `internal/client/generated.go` |
| `templates` | query | `templates` | implemented | `internal/client/generated.go` |
| `timeSchedule` | query | `timeSchedule` | implemented | `internal/client/generated.go` |
| `timeSchedules` | query | `timeSchedules` | implemented | `internal/client/generated.go` |
| `triageResponsibilities` | query | `triageResponsibilities` | implemented | `internal/client/generated.go` |
| `triageResponsibility` | query | `triageResponsibility` | implemented | `internal/client/generated.go` |
| `triageResponsibility_manualSelection` | query | `triageResponsibility` | implemented | `internal/client/generated.go` |
| `user` | query | `user` | implemented | `internal/client/generated.go` |
| `users` | query | `users` | implemented | `internal/client/generated.go` |
| `viewer` | query | `viewer` | implemented | `internal/client/generated.go` |
| `workflowState` | query | `workflowState` | implemented | `internal/client/generated.go` |
| `workflowStates` | query | `workflowStates` | implemented | `internal/client/generated.go` |

## Repo Domain-Map Commands

| Domain | Command | Backing | Scope | Status | Evidence |
| --- | --- | --- | --- | --- | --- |
| Core target | `whoami` | `Query.viewer`, `User` | Reads the authenticated user. | implemented | `linctl --help` / public CLI tests |
| Core target | `target` | `Query.organization`, `Query.teams`, `Query.team`, `Query.projects`, `Query.project` | Resolves the active token's organization, team, and optional project. | implemented | `linctl --help` / public CLI tests |
| Core target | `doctor` | `Query.viewer`, `Query.teams`, optional `Query.project` | Read-only health check for config load, token presence, and pinned-target confirmation. Does not print token values. | accepted_gap | planned in `docs/domain-map.md` |
| Core target | `application info` | `Query.applicationInfo` | Read-only public OAuth application metadata by client id. | implemented | `linctl --help` / public CLI tests |
| Core target | `organization exists` | `Query.organizationExists` | Read-only URL-key existence check for workspace lookup. | implemented | `linctl --help` / public CLI tests |
| Core target | `organization templates` | `Organization.templates` via `Query.organization` | Read-only workspace-level templates. | implemented | `linctl --help` / public CLI tests |
| Core target | `rate-limit status` | `Query.rateLimitStatus` | Read-only quota status for the authenticated Linear client. | implemented | `linctl --help` / public CLI tests |
| AgentActivity | `agent-activity list` | `Query.agentActivities` | Read-only | implemented | `linctl --help` / public CLI tests |
| AgentActivity | `agent-activity get` | `Query.agentActivity` | Read-only | implemented | `linctl --help` / public CLI tests |
| AgentActivity | `agent-activity create` | `Mutation.agentActivityCreate` | Blocked: create writes into an agent session and needs explicit session/comment guard semantics | accepted_gap | planned in `docs/domain-map.md` |
| AgentActivity | `agent-activity update` | `Mutation.agentActivityUpdate` | Blocked: update must resolve the agent session and activity scope before mutation | accepted_gap | planned in `docs/domain-map.md` |
| AgentActivity | `agent-activity archive` | `Mutation.agentActivityArchive` | Blocked: destructive command needs explicit AgentActivity safety semantics | accepted_gap | planned in `docs/domain-map.md` |
| AgentSkill | `agent-skill list` | `Query.agentSkills` | Read-only | implemented | `linctl --help` / public CLI tests |
| AgentSkill | `agent-skill get` | `Query.agentSkill` | Read-only | implemented | `linctl --help` / public CLI tests |
| AgentSkill | `agent-skill create` | `Mutation.agentSkillCreate` | Blocked: create can expose reusable agent instructions and needs explicit team/owner guard semantics | accepted_gap | planned in `docs/domain-map.md` |
| AgentSkill | `agent-skill update` | `Mutation.agentSkillUpdate` | Blocked: update must resolve the AgentSkill's team and ownership scope before mutation | accepted_gap | planned in `docs/domain-map.md` |
| AgentSkill | `agent-skill archive` | `Mutation.agentSkillArchive` | Blocked: destructive command needs explicit AgentSkill safety semantics | accepted_gap | planned in `docs/domain-map.md` |
| AuditEntry | `audit-entry types` | `Query.auditEntryTypes` | Read-only | implemented | `linctl --help` / public CLI tests |
| AuditEntry | `audit-entry list` | `Query.auditEntries` | Blocked: audit log entries can expose actor, IP, country, and request metadata; needs an explicit admin/security output model | accepted_gap | planned in `docs/domain-map.md` |
| Notification | `notification list` | `Query.notifications` | Read-only | implemented | `linctl --help` / public CLI tests |
| Notification | `notification get` | `Query.notification` | Read-only | implemented | `linctl --help` / public CLI tests |
| Notification | `notification subscription list` | `Query.notificationSubscriptions` | Read-only | implemented | `linctl --help` / public CLI tests |
| Notification | `notification subscription get` | `Query.notificationSubscription` | Read-only | implemented | `linctl --help` / public CLI tests |
| Notification | `notification archive` | `Mutation.notificationArchive` | Blocked: mutates the authenticated user's inbox state; needs an explicit viewer-state safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification archive all` | `Mutation.notificationArchiveAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification update` | `Mutation.notificationUpdate` | Blocked: direct inbox-state mutation needs an explicit viewer-state safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification mark read all` | `Mutation.notificationMarkReadAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification mark unread all` | `Mutation.notificationMarkUnreadAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification snooze all` | `Mutation.notificationSnoozeAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification unsnooze all` | `Mutation.notificationUnsnoozeAll` | Blocked: bulk inbox mutation needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification category channel subscription update` | `Mutation.notificationCategoryChannelSubscriptionUpdate` | Blocked: viewer notification preference mutation needs an explicit consent model | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification subscription create` | `Mutation.notificationSubscriptionCreate` | Blocked: subscription writes can target several entity types and need explicit target-resolution semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification subscription update` | `Mutation.notificationSubscriptionUpdate` | Blocked: update must resolve the subscription target before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Notification | `notification subscription delete` | `Mutation.notificationSubscriptionDelete` | Blocked: destructive viewer preference command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Release | `release-pipeline list` | `Query.releasePipelines` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-pipeline get` | `Query.releasePipeline` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-pipeline releases` | `ReleasePipeline.releases` via `Query.releasePipeline` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-pipeline stages` | `ReleasePipeline.stages` via `Query.releasePipeline` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-stage list` | `Query.releaseStages` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-stage get` | `Query.releaseStage` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-stage releases` | `ReleaseStage.releases` via `Query.releaseStage` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release list` | `Query.releases` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release search` | `Query.releaseSearch` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release get` | `Query.release` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release history` | `Release.history` via `Query.release` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release links` | `Release.links` via `Query.release` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `external-link get` | `Query.entityExternalLink` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-note list` | `Query.releaseNotes` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-note get` | `Query.releaseNote` | Read-only | implemented | `linctl --help` / public CLI tests |
| Release | `release-pipeline create` | `Mutation.releasePipelineCreate` | Blocked: pipeline configuration is team/admin release surface and needs explicit guard semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-pipeline update` | `Mutation.releasePipelineUpdate` | Blocked: update must resolve and compare associated teams before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-pipeline archive` | `Mutation.releasePipelineArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-pipeline unarchive` | `Mutation.releasePipelineUnarchive` | Blocked: restore command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-pipeline delete` | `Mutation.releasePipelineDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Release | `release-stage create` | `Mutation.releaseStageCreate` | Blocked: release workflow configuration needs explicit pipeline/team guard semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-stage update` | `Mutation.releaseStageUpdate` | Blocked: update must resolve the stage's pipeline and teams before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-stage archive` | `Mutation.releaseStageArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-stage unarchive` | `Mutation.releaseStageUnarchive` | Blocked: restore command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release create` | `Mutation.releaseCreate` | Blocked: create must resolve pipeline/team guard semantics before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release update` | `Mutation.releaseUpdate` | Blocked: update must resolve the release pipeline/stage and associated teams before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release archive` | `Mutation.releaseArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release unarchive` | `Mutation.releaseUnarchive` | Blocked: restore command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release delete` | `Mutation.releaseDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Release | `release complete` | `Mutation.releaseComplete`, `Mutation.releaseCompleteByAccessKey` | Blocked: lifecycle transition and access-key behavior need explicit guard semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release sync` | `Mutation.releaseSync`, `Mutation.releaseSyncByAccessKey` | Blocked: sync mutates release associations and needs explicit guard semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-note create` | `Mutation.releaseNoteCreate` | Blocked: create must resolve release pipeline and release range semantics before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-note update` | `Mutation.releaseNoteUpdate` | Blocked: update must resolve covered releases and pipeline before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-note archive` | `Mutation.releaseNoteArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `release-note delete` | `Mutation.releaseNoteDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Release | `issue-to-release create` | `Mutation.issueToReleaseCreate` | Blocked: association write must compare issue and release scope before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `issue-to-release update` | `Mutation.issueToReleaseUpdate` | Blocked: association update must compare issue and release scope before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Release | `issue-to-release delete` | `Mutation.issueToReleaseDelete` | Blocked: destructive association command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Issue | `issue list` | `Query.issues`, optionally filtered by `Issue.team.id`, `Issue.state.type`, `Issue.project.id`, `Issue.assignee.id`, `Issue.labels.some.id`, `Issue.cycle.id`, `Issue.createdAt.gte` (`--created-after` / `--created-since`), `Issue.createdAt.lte`, `Issue.hasBlockedByRelations.eq`, or `Issue.hasBlockingRelations.eq`; `--blocked-by ISSUE` traverses `Issue.relations` with `IssueRelation.type == "blocks"` and returns matching `IssueRelation.relatedIssue`; `--all-teams` omits the team filter | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `issue search` | `Query.issues`, filtered by `Issue.searchableContent` | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `issue get` | `Query.issue` | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `issue deps` | `Query.issue`, `Issue.parent`, `Issue.children`, `Issue.relations`, `Issue.inverseRelations`; `IssueRelation.type == "blocks"` separates blocked issues from blockers | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `issue id` | Current checkout issue identifier from git/jj context | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `issue title` | `Query.issue` after current checkout or explicit issue resolution | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `issue url` | `Query.issue` after current checkout or explicit issue resolution | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `issue branch` | `Query.issue`, `Issue.branchName` | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `issue pr` | `Query.issue`; emits a local `gh pr create` title/body plan without calling GitHub | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `next --dry-run` | `Query.issues`, filtered by `Issue.team.id`, `Issue.state.type == "unstarted"`, and `Issue.hasBlockedByRelations.eq == false`; fetches `Issue.relations`, `Issue.priority`, and `Issue.createdAt`, then ranks by active unblock count, priority, and age before printing one candidate without checkout/worktree creation | Read-only | implemented | `linctl --help` / public CLI tests |
| Issue | `done` | Current checkout issue identifier, then `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved | implemented | `linctl --help` / public CLI tests |
| Issue | `issue create` | `Mutation.issueCreate` with `IssueCreateInput.teamId`, optional `projectId`; `--description-file` is resolved locally before mutation | Team-scoped unless `projectId` is set | implemented | `linctl --help` / public CLI tests |
| Issue | `issue update` | `Mutation.issueUpdate` with `IssueUpdateInput`; `--description-file` replaces description, while `--append` or `--append-file` first reads `Issue.description` and appends text before sending `description` | Resource-scoped when a project target is involved | implemented | `linctl --help` / public CLI tests |
| Issue | `issue start` | `Query.viewer`, `Query.workflowStates` filtered to `started`, then `Mutation.issueUpdate` with `IssueUpdateInput.assigneeId` and `stateId` | Resource-scoped when a project target is involved | implemented | `linctl --help` / public CLI tests |
| Issue | `issue comment` | `Mutation.commentCreate`; `--body -` reads stdin and `--body-file` reads a local file before mutation | Resource-scoped to the issue's resolved team/project | implemented | `linctl --help` / public CLI tests |
| Issue | `issue reply` | `Mutation.commentCreate` with `CommentCreateInput.parentId`; `--body-file` reads a local file before mutation | Resource-scoped to the issue's resolved team/project | implemented | `linctl --help` / public CLI tests |
| Issue | `issue close` | `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved | implemented | `linctl --help` / public CLI tests |
| Issue | `issue comments` | `Issue.comments` via `Query.issue` | Read-only | implemented | `linctl --help` / public CLI tests |
| Comment | `comment list` | `Query.comments` | Read-only | implemented | `linctl --help` / public CLI tests |
| Comment | `comment get` | `Query.comment` | Read-only | implemented | `linctl --help` / public CLI tests |
| Comment | `comment resolve` | `Mutation.commentResolve` | Blocked: resolving must first identify and compare the parent issue/project/update/document scope | blocked_needs_design | write command needs explicit target and safety semantics |
| Comment | `comment unresolve` | `Mutation.commentUnresolve` | Blocked: unresolving must first identify and compare the parent issue/project/update/document scope | blocked_needs_design | write command needs explicit target and safety semantics |
| Project | `project list` | `Query.projects` | Read-only | implemented | `linctl --help` / public CLI tests |
| Project | `project get` | `Query.project` | Read-only | implemented | `linctl --help` / public CLI tests |
| Project | `project create` | `Mutation.projectCreate` with `ProjectCreateInput.teamIds` | Team-scoped | implemented | `linctl --help` / public CLI tests |
| Project | `project update` | `Mutation.projectUpdate` with `ProjectUpdateInput` | Resource-scoped, compare `project_id` | implemented | `linctl --help` / public CLI tests |
| Project | `project archive` | `Mutation.projectArchive` | Resource-scoped, compare `project_id` | implemented | `linctl --help` / public CLI tests |
| Project | `project members` | `Project.members` plus `Mutation.projectUpdate` with `ProjectUpdateInput.memberIds` | Read-only for list, resource-scoped for writes | implemented | `linctl --help` / public CLI tests |
| Project | `project updates` | `Project.projectUpdates` | Read-only | implemented | `linctl --help` / public CLI tests |
| ProjectUpdate | `project-update list` | `Query.projectUpdates` | Read-only | implemented | `linctl --help` / public CLI tests |
| ProjectUpdate | `project-update get` | `Query.projectUpdate` | Read-only | implemented | `linctl --help` / public CLI tests |
| ProjectUpdate | `project-update create` | `Mutation.projectUpdateCreate` | Blocked: create must resolve and compare the target project before posting | blocked_needs_design | write command needs explicit target and safety semantics |
| ProjectUpdate | `project-update update` | `Mutation.projectUpdateUpdate` | Blocked: update must resolve and compare the owning project before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| ProjectUpdate | `project-update archive` | `Mutation.projectUpdateArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Cycle | `cycle list` | `Query.cycles` | Read-only | implemented | `linctl --help` / public CLI tests |
| Cycle | `cycle get` | `Query.cycle` | Read-only | implemented | `linctl --help` / public CLI tests |
| Cycle | `cycle create` | `Mutation.cycleCreate` | Team-scoped | implemented | `linctl --help` / public CLI tests |
| Cycle | `cycle update` | `Mutation.cycleUpdate` | Team-scoped | implemented | `linctl --help` / public CLI tests |
| Cycle | `cycle archive` | `Mutation.cycleArchive` | Team-scoped | implemented | `linctl --help` / public CLI tests |
| Sprint | `sprint current` | `Query.cycles` filtered to active/current cycles | Read-only | implemented | `linctl --help` / public CLI tests |
| Sprint | `sprint report` | `Query.cycle` plus `Cycle.issues` | Read-only | implemented | `linctl --help` / public CLI tests |
| ProjectMilestone | `project-milestone list` | `Project.projectMilestones` via `Query.project` | Read-only | implemented | `linctl --help` / public CLI tests |
| ProjectMilestone | `project-milestone get` | `Query.projectMilestone` | Read-only | implemented | `linctl --help` / public CLI tests |
| ProjectMilestone | `project-milestone create` | `Mutation.projectMilestoneCreate` with `projectId` | Resource-scoped, compare `project_id` | implemented | `linctl --help` / public CLI tests |
| ProjectMilestone | `project-milestone update` | `Mutation.projectMilestoneUpdate` | Resource-scoped, compare resolved project | implemented | `linctl --help` / public CLI tests |
| ProjectMilestone | `project-milestone delete` | `Mutation.projectMilestoneDelete` | Resource-scoped, compare resolved project | blocked_needs_design | destructive command needs explicit safety semantics |
| Document | `document list` | `Query.documents` | Read-only | implemented | `linctl --help` / public CLI tests |
| Document | `document get` | `Query.document` | Read-only | implemented | `linctl --help` / public CLI tests |
| Document | `document create` | `Mutation.documentCreate` with optional `projectId`, `teamId`, `issueId`, `cycleId` | Blocked: parent can be project, team, issue, or cycle; write guard needs explicit parent-resolution semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Document | `document update` | `Mutation.documentUpdate` | Blocked: update must resolve and compare the existing parent before changing content | blocked_needs_design | write command needs explicit target and safety semantics |
| Document | `document delete` | `Mutation.documentDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Label | `label list` | `Query.issueLabels` | Read-only | implemented | `linctl --help` / public CLI tests |
| Label | `label get` | `Query.issueLabel` | Read-only | implemented | `linctl --help` / public CLI tests |
| Label | `label create` | `Mutation.issueLabelCreate` with optional `teamId` | Blocked: optional team scope needs explicit org/team target behavior before writes | blocked_needs_design | write command needs explicit target and safety semantics |
| Label | `label update` | `Mutation.issueLabelUpdate` | Blocked: update must resolve and compare the label's owning team before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Label | `label delete` | `Mutation.issueLabelDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Team | `team list` | `Query.teams` | Read-only | implemented | `linctl --help` / public CLI tests |
| Team | `team get` | `Query.team` | Read-only | implemented | `linctl --help` / public CLI tests |
| Team | `team create` | `Mutation.teamCreate` | Blocked: organization administration surface needs an explicit admin safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Team | `team update` | `Mutation.teamUpdate` | Blocked: team metadata writes need stronger authority checks than ordinary target comparison | blocked_needs_design | write command needs explicit target and safety semantics |
| Team | `team delete` | `Mutation.teamDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Team | `team members` | `Team.members` | Read-only | implemented | `linctl --help` / public CLI tests |
| User | `user list` | `Query.users` | Read-only | implemented | `linctl --help` / public CLI tests |
| User | `user get` | `Query.user` | Read-only | implemented | `linctl --help` / public CLI tests |
| User | `user me` | `Query.viewer` | Read-only | implemented | `linctl --help` / public CLI tests |
| WorkflowState | `workflow-state list` | `Query.workflowStates` | Read-only | implemented | `linctl --help` / public CLI tests |
| WorkflowState | `workflow-state get` | `Query.workflowState` | Read-only | implemented | `linctl --help` / public CLI tests |
| WorkflowState | `workflow-state create` | `Mutation.workflowStateCreate` | Blocked: team workflow configuration needs an explicit admin safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| WorkflowState | `workflow-state update` | `Mutation.workflowStateUpdate` | Blocked: update must resolve and compare the owning team before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| WorkflowState | `workflow-state archive` | `Mutation.workflowStateArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| TimeSchedule | `time-schedule list` | `Query.timeSchedules` | Read-only | implemented | `linctl --help` / public CLI tests |
| TimeSchedule | `time-schedule get` | `Query.timeSchedule` | Read-only | implemented | `linctl --help` / public CLI tests |
| TimeSchedule | `time-schedule create` | `Mutation.timeScheduleCreate` | Blocked: schedule create needs explicit owner/admin safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| TimeSchedule | `time-schedule update` | `Mutation.timeScheduleUpdate` | Blocked: update must resolve schedule scope before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| TimeSchedule | `time-schedule delete` | `Mutation.timeScheduleDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| TimeSchedule | `time-schedule upsert-external` | `Mutation.timeScheduleUpsertExternal` | Blocked: external integration sync surface is not an ordinary agent workflow | blocked_needs_design | write command needs explicit target and safety semantics |
| TriageResponsibility | `triage-responsibility list` | `Query.triageResponsibilities` | Read-only | implemented | `linctl --help` / public CLI tests |
| TriageResponsibility | `triage-responsibility get` | `Query.triageResponsibility` | Read-only | implemented | `linctl --help` / public CLI tests |
| TriageResponsibility | `triage-responsibility manual-selection` | `TriageResponsibility.manualSelection` via `Query.triageResponsibility` | Read-only | implemented | `linctl --help` / public CLI tests |
| TriageResponsibility | `triage-responsibility create` | `Mutation.triageResponsibilityCreate` | Blocked: team triage configuration needs an explicit admin safety model | accepted_gap | planned in `docs/domain-map.md` |
| TriageResponsibility | `triage-responsibility update` | `Mutation.triageResponsibilityUpdate` | Blocked: update must resolve and compare the owning team before mutation | accepted_gap | planned in `docs/domain-map.md` |
| TriageResponsibility | `triage-responsibility delete` | `Mutation.triageResponsibilityDelete` | Blocked: destructive team triage configuration command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Template | `template list` | `Query.templates` | Read-only | implemented | `linctl --help` / public CLI tests |
| Template | `template get` | `Query.template` | Read-only | implemented | `linctl --help` / public CLI tests |
| Template | `template create` | `Mutation.templateCreate` | Blocked: create can be workspace-, team-, or pipeline-scoped and needs explicit guard semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Template | `template update` | `Mutation.templateUpdate` | Blocked: update must resolve and compare the template's workspace, team, or pipeline scope before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Template | `template delete` | `Mutation.templateDelete` | Blocked: destructive command needs explicit template-scope safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Initiative | `initiative list` | `Query.initiatives` | Read-only | implemented | `linctl --help` / public CLI tests |
| Initiative | `initiative get` | `Query.initiative` | Read-only | implemented | `linctl --help` / public CLI tests |
| Initiative | `initiative history` | `Initiative.history` via `Query.initiative` | Read-only | implemented | `linctl --help` / public CLI tests |
| Initiative | `initiative links` | `Initiative.links` via `Query.initiative` | Read-only | implemented | `linctl --help` / public CLI tests |
| Initiative | `initiative sub-initiatives` | `Initiative.subInitiatives` via `Query.initiative` | Read-only | implemented | `linctl --help` / public CLI tests |
| Initiative | `initiative updates` | `Initiative.initiativeUpdates` via `Query.initiative` | Read-only | implemented | `linctl --help` / public CLI tests |
| Initiative | `initiative create` | `Mutation.createInitiative` | Blocked: initiative create needs an explicit organization-scoped safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Initiative | `initiative update` | `Mutation.updateInitiative` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Initiative | `initiative archive` | `Mutation.archiveInitiative` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| InitiativeRelation | `initiative-relation list` | `Query.initiativeRelations` | Read-only | implemented | `linctl --help` / public CLI tests |
| InitiativeRelation | `initiative-relation get` | `Query.initiativeRelation` | Read-only | implemented | `linctl --help` / public CLI tests |
| InitiativeRelation | `initiative-relation create` | `Mutation.initiativeRelationCreate` | Blocked: create must resolve and compare both Initiative hierarchy endpoints before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| InitiativeRelation | `initiative-relation update` | `Mutation.initiativeRelationUpdate` | Blocked: update must resolve and compare both Initiative hierarchy endpoints before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| InitiativeRelation | `initiative-relation delete` | `Mutation.initiativeRelationDelete` | Blocked: destructive command needs explicit hierarchy safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| InitiativeToProject | `initiative-to-project list` | `Query.initiativeToProjects` | Read-only | implemented | `linctl --help` / public CLI tests |
| InitiativeToProject | `initiative-to-project get` | `Query.initiativeToProject` | Read-only | implemented | `linctl --help` / public CLI tests |
| InitiativeToProject | `initiative-to-project create` | `Mutation.initiativeToProjectCreate` | Blocked: create must resolve and compare both Initiative and Project endpoints before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| InitiativeToProject | `initiative-to-project update` | `Mutation.initiativeToProjectUpdate` | Blocked: update must resolve and compare both Initiative and Project endpoints before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| InitiativeToProject | `initiative-to-project delete` | `Mutation.initiativeToProjectDelete` | Blocked: destructive command needs explicit association safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| InitiativeUpdate | `initiative-update list` | `Query.initiativeUpdates` | Read-only | implemented | `linctl --help` / public CLI tests |
| InitiativeUpdate | `initiative-update get` | `Query.initiativeUpdate` | Read-only | implemented | `linctl --help` / public CLI tests |
| InitiativeUpdate | `initiative-update create` | `Mutation.initiativeUpdateCreate` | Blocked: create must resolve and compare the owning Initiative before posting | blocked_needs_design | write command needs explicit target and safety semantics |
| InitiativeUpdate | `initiative-update update` | `Mutation.initiativeUpdateUpdate` | Blocked: update must resolve and compare the owning Initiative before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| InitiativeUpdate | `initiative-update archive` | `Mutation.initiativeUpdateArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| InitiativeUpdate | `initiative-update unarchive` | `Mutation.initiativeUpdateUnarchive` | Blocked: unarchive needs explicit lifecycle and target semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Roadmap | `roadmap list` | `Query.roadmaps` | Read-only | implemented | `linctl --help` / public CLI tests |
| Roadmap | `roadmap get` | `Query.roadmap` | Read-only | implemented | `linctl --help` / public CLI tests |
| Roadmap | `roadmap create` | `Mutation.roadmapCreate` | Blocked: deprecated organization-scoped planning surface needs an explicit safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Roadmap | `roadmap update` | `Mutation.roadmapUpdate` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Roadmap | `roadmap archive` | `Mutation.roadmapArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Roadmap | `roadmap delete` | `Mutation.roadmapDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| CustomView | `custom-view list` | `Query.customViews` | Read-only | implemented | `linctl --help` / public CLI tests |
| CustomView | `custom-view subscribers` | `Query.customViewHasSubscribers` | Read-only | implemented | `linctl --help` / public CLI tests |
| CustomView | `custom-view get` | `Query.customView` | Read-only | implemented | `linctl --help` / public CLI tests |
| CustomView | `custom-view create` | `Mutation.createCustomView` | Blocked: custom view create needs an explicit organization-scoped safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| CustomView | `custom-view update` | `Mutation.updateCustomView` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| CustomView | `custom-view delete` | `Mutation.deleteCustomView` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Customer | `customer list` | `Query.customers` | Read-only | implemented | `linctl --help` / public CLI tests |
| Customer | `customer get` | `Query.customer` | Read-only | implemented | `linctl --help` / public CLI tests |
| Customer | `customer-need list` | `Query.customerNeeds` | Read-only | implemented | `linctl --help` / public CLI tests |
| Customer | `customer-need get` | `Query.customerNeed` | Read-only | implemented | `linctl --help` / public CLI tests |
| Customer | `customer-status list` | `Query.customerStatuses` | Read-only | implemented | `linctl --help` / public CLI tests |
| Customer | `customer-status get` | `Query.customerStatus` | Read-only | implemented | `linctl --help` / public CLI tests |
| Customer | `customer-tier list` | `Query.customerTiers` | Read-only | implemented | `linctl --help` / public CLI tests |
| Customer | `customer-tier get` | `Query.customerTier` | Read-only | implemented | `linctl --help` / public CLI tests |
| Customer | `customer create` | `Mutation.customerCreate` | Blocked: customer create needs an explicit organization-scoped safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer update` | `Mutation.customerUpdate` | Blocked: update must resolve and compare the owning organization before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer archive` | `Mutation.customerArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer-need create` | `Mutation.customerNeedCreate` | Blocked: need creation must prove the linked issue, project, or customer target before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer-need update` | `Mutation.customerNeedUpdate` | Blocked: update must resolve the need and compare the linked issue or project target before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer-need archive` | `Mutation.customerNeedArchive` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer-need delete` | `Mutation.customerNeedDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Customer | `customer-status create` | `Mutation.customerStatusCreate` | Blocked: workspace lifecycle configuration needs an explicit admin safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer-status update` | `Mutation.customerStatusUpdate` | Blocked: workspace lifecycle configuration needs an explicit admin safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer-status delete` | `Mutation.customerStatusDelete` | Blocked: destructive admin command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Customer | `customer-tier create` | `Mutation.customerTierCreate` | Blocked: workspace tier configuration needs an explicit admin safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer-tier update` | `Mutation.customerTierUpdate` | Blocked: workspace tier configuration needs an explicit admin safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Customer | `customer-tier delete` | `Mutation.customerTierDelete` | Blocked: destructive admin command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Favorite | `favorite list` | `Query.favorites` | Read-only | implemented | `linctl --help` / public CLI tests |
| Favorite | `favorite children` | `Favorite.children` via `Query.favorite` | Read-only | implemented | `linctl --help` / public CLI tests |
| Favorite | `favorite get` | `Query.favorite` | Read-only | implemented | `linctl --help` / public CLI tests |
| Favorite | `favorite create` | `Mutation.createFavorite` | Blocked: favorite create needs an explicit viewer-scoped safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Favorite | `favorite update` | `Mutation.updateFavorite` | Blocked: update must resolve and compare the owning viewer before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Favorite | `favorite delete` | `Mutation.deleteFavorite` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Emoji | `emoji list` | `Query.emojis` | Read-only | implemented | `linctl --help` / public CLI tests |
| Emoji | `emoji get` | `Query.emoji` | Read-only | implemented | `linctl --help` / public CLI tests |
| Emoji | `emoji create` | `Mutation.createEmoji` | Blocked: emoji create needs an explicit organization-scoped safety model | blocked_needs_design | write command needs explicit target and safety semantics |
| Emoji | `emoji delete` | `Mutation.deleteEmoji` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |
| Attachment | `attachment list` | `Query.attachments` | Read-only | implemented | `linctl --help` / public CLI tests |
| Attachment | `attachment url` | `Query.attachmentsForURL` | Read-only | implemented | `linctl --help` / public CLI tests |
| Attachment | `attachment get` | `Query.attachment` | Read-only | implemented | `linctl --help` / public CLI tests |
| Attachment | `attachment create` | `Mutation.attachmentCreate` | Blocked: attachment create must resolve and compare the owning issue's team before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Attachment | `attachment update` | `Mutation.attachmentUpdate` | Blocked: update must resolve and compare the owning issue before mutation | blocked_needs_design | write command needs explicit target and safety semantics |
| Attachment | `attachment delete` | `Mutation.attachmentDelete` | Blocked: destructive command needs explicit safety semantics | blocked_needs_design | destructive command needs explicit safety semantics |

