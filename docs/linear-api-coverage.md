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
| Upstream SDK root methods | 458 | 34 | 458 |
| Upstream Query root fields | 158 | 22 | 158 |
| Upstream Mutation root fields | 364 | 12 | 364 |
| Local generated Go operations | 60 | 60 | 60 |
| Domain-map commands | 73 | 54 | 73 |

## Upstream SDK Root Methods

| Method | Kind | Status | Evidence |
| --- | --- | --- | --- |
| `administrableTeams` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentActivities` | method | safe_candidate | read operation may fit future CLI coverage |
| `agentActivity` | method | safe_candidate | read operation may fit future CLI coverage |
| `agentSession` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnComment` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnIssue` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionUpdateExternalUrl` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessions` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSkill` | method | safe_candidate | read operation may fit future CLI coverage |
| `agentSkills` | method | safe_candidate | read operation may fit future CLI coverage |
| `airbyteIntegrationConnect` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `applicationInfo` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `attachment` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `attachments` | method | safe_candidate | read operation may fit future CLI coverage |
| `attachmentsForURL` | method | safe_candidate | read operation may fit future CLI coverage |
| `auditEntries` | method | safe_candidate | read operation may fit future CLI coverage |
| `auditEntryTypes` | getter | safe_candidate | read operation may fit future CLI coverage |
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
| `customView` | method | safe_candidate | read operation may fit future CLI coverage |
| `customViewHasSubscribers` | method | safe_candidate | read operation may fit future CLI coverage |
| `customViews` | method | safe_candidate | read operation may fit future CLI coverage |
| `customer` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerMerge` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerNeed` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerNeedCreateFromAttachment` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerNeeds` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerStatus` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerStatuses` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerTier` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerTiers` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerUnsync` | method | safe_candidate | read operation may fit future CLI coverage |
| `customerUpsert` | method | safe_candidate | read operation may fit future CLI coverage |
| `customers` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `emailIntakeAddressRotate` | method | safe_candidate | read operation may fit future CLI coverage |
| `emailTokenUserAccountAuth` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emailUnsubscribe` | method | safe_candidate | read operation may fit future CLI coverage |
| `emailUserAccountAuthChallenge` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emoji` | method | safe_candidate | read operation may fit future CLI coverage |
| `emojis` | method | safe_candidate | read operation may fit future CLI coverage |
| `entityExternalLink` | method | safe_candidate | read operation may fit future CLI coverage |
| `externalUser` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `externalUsers` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `favorite` | method | safe_candidate | read operation may fit future CLI coverage |
| `favorites` | method | safe_candidate | read operation may fit future CLI coverage |
| `fileUpload` | method | safe_candidate | read operation may fit future CLI coverage |
| `googleUserAccountAuth` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `imageUploadFromUrl` | method | safe_candidate | read operation may fit future CLI coverage |
| `importFileUpload` | method | safe_candidate | read operation may fit future CLI coverage |
| `initiative` | method | safe_candidate | read operation may fit future CLI coverage |
| `initiativeRelation` | method | safe_candidate | read operation may fit future CLI coverage |
| `initiativeRelations` | method | safe_candidate | read operation may fit future CLI coverage |
| `initiativeToProject` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `initiativeToProjects` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `initiativeUpdate` | method | safe_candidate | read operation may fit future CLI coverage |
| `initiativeUpdates` | method | safe_candidate | read operation may fit future CLI coverage |
| `initiatives` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `notification` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationArchiveAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationMarkReadAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationMarkUnreadAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationSnoozeAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationSubscription` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationSubscriptions` | method | safe_candidate | read operation may fit future CLI coverage |
| `notificationUnsnoozeAll` | method | safe_candidate | read operation may fit future CLI coverage |
| `notifications` | method | safe_candidate | read operation may fit future CLI coverage |
| `organization` | getter | implemented | local operation or command exists |
| `organizationDeleteChallenge` | getter | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationExists` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `rateLimitStatus` | getter | safe_candidate | read operation may fit future CLI coverage |
| `recentReleasesByAccessKey` | method | safe_candidate | read operation may fit future CLI coverage |
| `refreshGoogleSheetsData` | method | safe_candidate | read operation may fit future CLI coverage |
| `release` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseComplete` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseCompleteByAccessKey` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseNote` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseNotes` | method | safe_candidate | read operation may fit future CLI coverage |
| `releasePipeline` | method | safe_candidate | read operation may fit future CLI coverage |
| `releasePipelineByAccessKey` | getter | safe_candidate | read operation may fit future CLI coverage |
| `releasePipelines` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseSearch` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseStage` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseStages` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseSync` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseSyncByAccessKey` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseUpdateByPipeline` | method | safe_candidate | read operation may fit future CLI coverage |
| `releaseUpdateByPipelineByAccessKey` | method | safe_candidate | read operation may fit future CLI coverage |
| `releases` | method | safe_candidate | read operation may fit future CLI coverage |
| `resendOrganizationInvite` | method | safe_candidate | read operation may fit future CLI coverage |
| `resendOrganizationInviteByEmail` | method | safe_candidate | read operation may fit future CLI coverage |
| `roadmap` | method | safe_candidate | read operation may fit future CLI coverage |
| `roadmapToProject` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmapToProjects` | method | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmaps` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `template` | method | safe_candidate | read operation may fit future CLI coverage |
| `templates` | getter | safe_candidate | read operation may fit future CLI coverage |
| `templatesForIntegration` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeSchedule` | method | safe_candidate | read operation may fit future CLI coverage |
| `timeScheduleRefreshIntegrationSchedule` | method | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeScheduleUpsertExternal` | method | safe_candidate | read operation may fit future CLI coverage |
| `timeSchedules` | method | safe_candidate | read operation may fit future CLI coverage |
| `trackAnonymousEvent` | method | safe_candidate | read operation may fit future CLI coverage |
| `triageResponsibilities` | method | safe_candidate | read operation may fit future CLI coverage |
| `triageResponsibility` | method | safe_candidate | read operation may fit future CLI coverage |
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
| `agentActivities` | `AgentActivityConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `agentActivity` | `AgentActivity!` | safe_candidate | read operation may fit future CLI coverage |
| `agentSession` | `AgentSession!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionSandbox` | `CodingAgentSandboxPayload` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessions` | `AgentSessionConnection!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSkill` | `AgentSkill!` | safe_candidate | read operation may fit future CLI coverage |
| `agentSkills` | `AgentSkillConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `applicationInfo` | `Application!` | safe_candidate | read operation may fit future CLI coverage |
| `archivedIntegrations` | `[ArchivedIntegrationPayload!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `archivedTeams` | `[Team!]!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `attachment` | `Attachment!` | safe_candidate | read operation may fit future CLI coverage |
| `attachmentIssue` | `Issue!` | accepted_gap | repo-planned or likely useful CLI domain |
| `attachmentSources` | `AttachmentSourcesPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `attachments` | `AttachmentConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `attachmentsForURL` | `AttachmentConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `auditEntries` | `AuditEntryConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `auditEntryTypes` | `[AuditEntryType!]!` | safe_candidate | read operation may fit future CLI coverage |
| `authenticationSessions` | `[AuthenticationSessionResponse!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `availableUsers` | `AuthResolverResponse!` | accepted_gap | repo-planned or likely useful CLI domain |
| `comment` | `Comment!` | implemented | root field used by local GraphQL operation |
| `comments` | `CommentConnection!` | implemented | root field used by local GraphQL operation |
| `customView` | `CustomView!` | safe_candidate | read operation may fit future CLI coverage |
| `customViewDetailsSuggestion` | `CustomViewSuggestionPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `customViewHasSubscribers` | `CustomViewHasSubscribersPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `customViews` | `CustomViewConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `customer` | `Customer!` | safe_candidate | read operation may fit future CLI coverage |
| `customerNeed` | `CustomerNeed!` | safe_candidate | read operation may fit future CLI coverage |
| `customerNeeds` | `CustomerNeedConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `customerStatus` | `CustomerStatus!` | safe_candidate | read operation may fit future CLI coverage |
| `customerStatuses` | `CustomerStatusConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `customerTier` | `CustomerTier!` | safe_candidate | read operation may fit future CLI coverage |
| `customerTiers` | `CustomerTierConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `customers` | `CustomerConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `cycle` | `Cycle!` | implemented | root field used by local GraphQL operation |
| `cycles` | `CycleConnection!` | implemented | root field used by local GraphQL operation |
| `document` | `Document!` | implemented | root field used by local GraphQL operation |
| `documentContentHistory` | `DocumentContentHistoryPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `documentContentHistoryEntries` | `DocumentContentHistoryPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `documentContentHistoryTimeline` | `DocumentContentHistoryTimelinePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `documents` | `DocumentConnection!` | implemented | root field used by local GraphQL operation |
| `emailIntakeAddress` | `EmailIntakeAddress!` | safe_candidate | read operation may fit future CLI coverage |
| `emoji` | `Emoji!` | safe_candidate | read operation may fit future CLI coverage |
| `emojis` | `EmojiConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `entityExternalLink` | `EntityExternalLink!` | safe_candidate | read operation may fit future CLI coverage |
| `externalUser` | `ExternalUser!` | accepted_gap | repo-planned or likely useful CLI domain |
| `externalUsers` | `ExternalUserConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `failuresForOauthWebhooks` | `[WebhookFailureEvent!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `favorite` | `Favorite!` | safe_candidate | read operation may fit future CLI coverage |
| `favorites` | `FavoriteConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `fetchData` | `FetchDataPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `initiative` | `Initiative!` | safe_candidate | read operation may fit future CLI coverage |
| `initiativeRelation` | `InitiativeRelation!` | safe_candidate | read operation may fit future CLI coverage |
| `initiativeRelations` | `InitiativeRelationConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `initiativeToProject` | `InitiativeToProject!` | accepted_gap | repo-planned or likely useful CLI domain |
| `initiativeToProjects` | `InitiativeToProjectConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `initiativeUpdate` | `InitiativeUpdate!` | safe_candidate | read operation may fit future CLI coverage |
| `initiativeUpdates` | `InitiativeUpdateConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `initiatives` | `InitiativeConnection!` | safe_candidate | read operation may fit future CLI coverage |
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
| `issueSearch` | `IssueConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueTitleSuggestionFromCustomerRequest` | `IssueTitleSuggestionFromCustomerRequestPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToRelease` | `IssueToRelease!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToReleases` | `IssueToReleaseConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueVcsBranchSearch` | `Issue` | accepted_gap | repo-planned or likely useful CLI domain |
| `issues` | `IssueConnection!` | implemented | root field used by local GraphQL operation |
| `latestReleaseByAccessKey` | `Release` | safe_candidate | read operation may fit future CLI coverage |
| `microsoftTeamsChannels` | `MicrosoftTeamsChannelsPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `notification` | `Notification!` | safe_candidate | read operation may fit future CLI coverage |
| `notificationSubscription` | `NotificationSubscription!` | safe_candidate | read operation may fit future CLI coverage |
| `notificationSubscriptions` | `NotificationSubscriptionConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `notifications` | `NotificationConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `notificationsUnreadCount` | `Int!` | safe_candidate | read operation may fit future CLI coverage |
| `oauthApplication` | `OAuthApplication!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplications` | `[OAuthApplication!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `organization` | `Organization!` | implemented | root field used by local GraphQL operation |
| `organizationDomainClaimRequest` | `OrganizationDomainClaimPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `organizationExists` | `OrganizationExistsPayload!` | safe_candidate | read operation may fit future CLI coverage |
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
| `rateLimitStatus` | `RateLimitPayload!` | safe_candidate | read operation may fit future CLI coverage |
| `recentReleasesByAccessKey` | `[Release!]!` | safe_candidate | read operation may fit future CLI coverage |
| `release` | `Release!` | safe_candidate | read operation may fit future CLI coverage |
| `releaseNote` | `ReleaseNote!` | safe_candidate | read operation may fit future CLI coverage |
| `releaseNotes` | `ReleaseNoteConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `releasePipeline` | `ReleasePipeline!` | safe_candidate | read operation may fit future CLI coverage |
| `releasePipelineByAccessKey` | `ReleasePipeline!` | safe_candidate | read operation may fit future CLI coverage |
| `releasePipelines` | `ReleasePipelineConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `releaseSearch` | `[Release!]!` | safe_candidate | read operation may fit future CLI coverage |
| `releaseStage` | `ReleaseStage!` | safe_candidate | read operation may fit future CLI coverage |
| `releaseStages` | `ReleaseStageConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `releases` | `ReleaseConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `roadmap` | `Roadmap!` | safe_candidate | read operation may fit future CLI coverage |
| `roadmapToProject` | `RoadmapToProject!` | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmapToProjects` | `RoadmapToProjectConnection!` | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmaps` | `RoadmapConnection!` | safe_candidate | read operation may fit future CLI coverage |
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
| `template` | `Template!` | safe_candidate | read operation may fit future CLI coverage |
| `templates` | `[Template!]!` | safe_candidate | read operation may fit future CLI coverage |
| `templatesForIntegration` | `[Template!]!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeSchedule` | `TimeSchedule!` | safe_candidate | read operation may fit future CLI coverage |
| `timeSchedules` | `TimeScheduleConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `triageResponsibilities` | `TriageResponsibilityConnection!` | safe_candidate | read operation may fit future CLI coverage |
| `triageResponsibility` | `TriageResponsibility!` | safe_candidate | read operation may fit future CLI coverage |
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
| `agentActivityCreate` | `AgentActivityPayload!` | blocked_needs_design | mutation needs product and safety design |
| `agentActivityCreatePrompt` | `AgentActivityPayload!` | blocked_needs_design | mutation needs product and safety design |
| `agentActivityDeleteQueued` | `AgentActivityPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `agentActivitySendQueued` | `AgentActivityPayload!` | blocked_needs_design | mutation needs product and safety design |
| `agentSessionCreate` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnComment` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionCreateOnIssue` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionUpdate` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSessionUpdateExternalUrl` | `AgentSessionPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `agentSkillCreate` | `AgentSkillPayload!` | blocked_needs_design | mutation needs product and safety design |
| `agentSkillDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `agentSkillUpdate` | `AgentSkillPayload!` | blocked_needs_design | mutation needs product and safety design |
| `airbyteIntegrationConnect` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `attachmentCreate` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
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
| `attachmentUpdate` | `AttachmentPayload!` | blocked_needs_design | mutation needs product and safety design |
| `commentCreate` | `CommentPayload!` | implemented | root field used by local GraphQL operation |
| `commentDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `commentResolve` | `CommentPayload!` | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `commentUnresolve` | `CommentPayload!` | blocked_needs_design | state-changing operation needs guarded target semantics before exposure |
| `commentUpdate` | `CommentPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `contactCreate` | `ContactPayload!` | blocked_needs_design | mutation needs product and safety design |
| `contactSalesCreate` | `ContactPayload!` | blocked_needs_design | mutation needs product and safety design |
| `createCsvExportReport` | `CreateCsvExportReportPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createInitiativeUpdateReminder` | `InitiativeUpdateReminderPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createOrganizationFromOnboarding` | `CreateOrJoinOrganizationResponse!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `createProjectUpdateReminder` | `ProjectUpdateReminderPayload!` | blocked_needs_design | write operation needs guarded target semantics before exposure |
| `customViewCreate` | `CustomViewPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customViewDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customViewUpdate` | `CustomViewPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerCreate` | `CustomerPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customerMerge` | `CustomerPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerNeedArchive` | `CustomerNeedArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerNeedCreate` | `CustomerNeedPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerNeedCreateFromAttachment` | `CustomerNeedPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerNeedDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customerNeedUnarchive` | `CustomerNeedArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerNeedUpdate` | `CustomerNeedUpdatePayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerStatusCreate` | `CustomerStatusPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerStatusDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customerStatusUpdate` | `CustomerStatusPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerTierCreate` | `CustomerTierPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerTierDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `customerTierUpdate` | `CustomerTierPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerUnsync` | `CustomerPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerUpdate` | `CustomerPayload!` | blocked_needs_design | mutation needs product and safety design |
| `customerUpsert` | `CustomerPayload!` | blocked_needs_design | mutation needs product and safety design |
| `cycleArchive` | `CycleArchivePayload!` | implemented | root field used by local GraphQL operation |
| `cycleCreate` | `CyclePayload!` | implemented | root field used by local GraphQL operation |
| `cycleShiftAll` | `CyclePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `cycleStartUpcomingCycleToday` | `CyclePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `cycleUpdate` | `CyclePayload!` | implemented | root field used by local GraphQL operation |
| `documentCreate` | `DocumentPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `documentDelete` | `DocumentArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `documentUnarchive` | `DocumentArchivePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `documentUpdate` | `DocumentPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `emailIntakeAddressCreate` | `EmailIntakeAddressPayload!` | blocked_needs_design | mutation needs product and safety design |
| `emailIntakeAddressDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `emailIntakeAddressRefreshSesDomainStatus` | `EmailIntakeAddressRefreshSesDomainStatusPayload!` | blocked_needs_design | mutation needs product and safety design |
| `emailIntakeAddressRotate` | `EmailIntakeAddressPayload!` | blocked_needs_design | mutation needs product and safety design |
| `emailIntakeAddressUpdate` | `EmailIntakeAddressPayload!` | blocked_needs_design | mutation needs product and safety design |
| `emailTokenUserAccountAuth` | `AuthResolverResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emailUnsubscribe` | `EmailUnsubscribePayload!` | blocked_needs_design | mutation needs product and safety design |
| `emailUserAccountAuthChallenge` | `EmailUserAccountAuthChallengeResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `emojiCreate` | `EmojiPayload!` | blocked_needs_design | mutation needs product and safety design |
| `emojiDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `entityExternalLinkCreate` | `EntityExternalLinkPayload!` | blocked_needs_design | mutation needs product and safety design |
| `entityExternalLinkDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `entityExternalLinkUpdate` | `EntityExternalLinkPayload!` | blocked_needs_design | mutation needs product and safety design |
| `favoriteCreate` | `FavoritePayload!` | blocked_needs_design | mutation needs product and safety design |
| `favoriteDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `favoriteUpdate` | `FavoritePayload!` | blocked_needs_design | mutation needs product and safety design |
| `fileUpload` | `UploadPayload!` | blocked_needs_design | mutation needs product and safety design |
| `fileUploadDangerouslyDelete` | `FileUploadDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `gitAutomationStateCreate` | `GitAutomationStatePayload!` | blocked_needs_design | mutation needs product and safety design |
| `gitAutomationStateDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `gitAutomationStateUpdate` | `GitAutomationStatePayload!` | blocked_needs_design | mutation needs product and safety design |
| `gitAutomationTargetBranchCreate` | `GitAutomationTargetBranchPayload!` | blocked_needs_design | mutation needs product and safety design |
| `gitAutomationTargetBranchDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `gitAutomationTargetBranchUpdate` | `GitAutomationTargetBranchPayload!` | blocked_needs_design | mutation needs product and safety design |
| `googleUserAccountAuth` | `AuthResolverResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `imageUploadFromUrl` | `ImageUploadFromUrlPayload!` | blocked_needs_design | mutation needs product and safety design |
| `importFileUpload` | `UploadPayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeAddLabel` | `InitiativePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `initiativeArchive` | `InitiativeArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeCreate` | `InitiativePayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `initiativeRelationCreate` | `InitiativeRelationPayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeRelationDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `initiativeRelationUpdate` | `InitiativeRelationPayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeRemoveLabel` | `InitiativePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `initiativeToProjectCreate` | `InitiativeToProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `initiativeToProjectDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `initiativeToProjectUpdate` | `InitiativeToProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `initiativeUnarchive` | `InitiativeArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeUpdate` | `InitiativePayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeUpdateArchive` | `InitiativeUpdateArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeUpdateCreate` | `InitiativeUpdatePayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeUpdateUnarchive` | `InitiativeUpdateArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `initiativeUpdateUpdate` | `InitiativeUpdatePayload!` | blocked_needs_design | mutation needs product and safety design |
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
| `issueBatchCreate` | `IssueBatchPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueBatchUpdate` | `IssueBatchPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `issueImportUpdate` | `IssueImportPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabelCreate` | `IssueLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabelDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueLabelRestore` | `IssueLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabelRetire` | `IssueLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueLabelUpdate` | `IssueLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRelationCreate` | `IssueRelationPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRelationDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueRelationUpdate` | `IssueRelationPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueReminder` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueRemoveLabel` | `IssuePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueShare` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueSubscribe` | `IssuePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToReleaseCreate` | `IssueToReleasePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `issueToReleaseDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueToReleaseDeleteByIssueAndRelease` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `issueUnarchive` | `IssueArchivePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
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
| `notificationArchive` | `NotificationArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationArchiveAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationCategoryChannelSubscriptionUpdate` | `UserSettingsPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationMarkReadAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationMarkUnreadAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationSnoozeAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationSubscriptionCreate` | `NotificationSubscriptionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationSubscriptionDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `notificationSubscriptionUpdate` | `NotificationSubscriptionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationUnarchive` | `NotificationArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationUnsnoozeAll` | `NotificationBatchActionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `notificationUpdate` | `NotificationPayload!` | blocked_needs_design | mutation needs product and safety design |
| `oauthApplicationArchive` | `OAuthApplicationArchivePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplicationCreate` | `OAuthApplicationCreatePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplicationRotateSecret` | `OAuthApplicationRotateSecretPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplicationRotateWebhookSecret` | `OAuthApplicationRotateWebhookSecretPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `oauthApplicationUpdate` | `OAuthApplicationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `organizationCancelDelete` | `OrganizationCancelDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationDelete` | `OrganizationDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationDeleteChallenge` | `OrganizationDeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationDomainClaim` | `OrganizationDomainSimplePayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationDomainCreate` | `OrganizationDomainPayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationDomainDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationDomainUpdate` | `OrganizationDomainPayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationDomainVerify` | `OrganizationDomainPayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationInviteCreate` | `OrganizationInvitePayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationInviteDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `organizationInviteUpdate` | `OrganizationInvitePayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationStartTrial` | `OrganizationStartTrialPayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationStartTrialForPlan` | `OrganizationStartTrialPayload!` | blocked_needs_design | mutation needs product and safety design |
| `organizationUpdate` | `OrganizationPayload!` | blocked_needs_design | mutation needs product and safety design |
| `passkeyLoginFinish` | `AuthResolverResponse!` | blocked_needs_design | mutation needs product and safety design |
| `passkeyLoginStart` | `PasskeyLoginStartResponse!` | blocked_needs_design | mutation needs product and safety design |
| `projectAddLabel` | `ProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectArchive` | `ProjectArchivePayload!` | implemented | root field used by local GraphQL operation |
| `projectCreate` | `ProjectPayload!` | implemented | root field used by local GraphQL operation |
| `projectCreateSlackChannel` | `ProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectDelete` | `ProjectArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectExternalSyncDisable` | `ProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelCreate` | `ProjectLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectLabelRestore` | `ProjectLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelRetire` | `ProjectLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectLabelUpdate` | `ProjectLabelPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectMilestoneCreate` | `ProjectMilestonePayload!` | implemented | root field used by local GraphQL operation |
| `projectMilestoneDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectMilestoneMove` | `ProjectMilestoneMovePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectMilestoneUpdate` | `ProjectMilestonePayload!` | implemented | root field used by local GraphQL operation |
| `projectReassignStatus` | `SuccessPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectRelationCreate` | `ProjectRelationPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectRelationDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectRelationUpdate` | `ProjectRelationPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectRemoveLabel` | `ProjectPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectStatusArchive` | `ProjectStatusArchivePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectStatusCreate` | `ProjectStatusPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectStatusUnarchive` | `ProjectStatusArchivePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectStatusUpdate` | `ProjectStatusPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectUnarchive` | `ProjectArchivePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectUpdate` | `ProjectPayload!` | implemented | root field used by local GraphQL operation |
| `projectUpdateArchive` | `ProjectUpdateArchivePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectUpdateCreate` | `ProjectUpdatePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectUpdateDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `projectUpdateUnarchive` | `ProjectUpdateArchivePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `projectUpdateUpdate` | `ProjectUpdatePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `pushSubscriptionCreate` | `PushSubscriptionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `pushSubscriptionDelete` | `PushSubscriptionPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `reactionCreate` | `ReactionPayload!` | blocked_needs_design | mutation needs product and safety design |
| `reactionDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `refreshGoogleSheetsData` | `IntegrationPayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseArchive` | `ReleaseArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseComplete` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseCompleteByAccessKey` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseCreate` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseDelete` | `ReleaseArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `releaseNoteCreate` | `ReleaseNotePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseNoteDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `releaseNoteUpdate` | `ReleaseNotePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releasePipelineArchive` | `ReleasePipelineArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releasePipelineCreate` | `ReleasePipelinePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releasePipelineDelete` | `ReleasePipelineArchivePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `releasePipelineUnarchive` | `ReleasePipelineArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releasePipelineUpdate` | `ReleasePipelinePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseStageArchive` | `ReleaseStageArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseStageCreate` | `ReleaseStagePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseStageUnarchive` | `ReleaseStageArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseStageUpdate` | `ReleaseStagePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseSync` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseSyncByAccessKey` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseUnarchive` | `ReleaseArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseUpdate` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseUpdateByPipeline` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `releaseUpdateByPipelineByAccessKey` | `ReleasePayload!` | blocked_needs_design | mutation needs product and safety design |
| `resendOrganizationInvite` | `DeletePayload!` | blocked_needs_design | mutation needs product and safety design |
| `resendOrganizationInviteByEmail` | `DeletePayload!` | blocked_needs_design | mutation needs product and safety design |
| `roadmapArchive` | `RoadmapArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `roadmapCreate` | `RoadmapPayload!` | blocked_needs_design | mutation needs product and safety design |
| `roadmapDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `roadmapToProjectCreate` | `RoadmapToProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmapToProjectDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `roadmapToProjectUpdate` | `RoadmapToProjectPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `roadmapUnarchive` | `RoadmapArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `roadmapUpdate` | `RoadmapPayload!` | blocked_needs_design | mutation needs product and safety design |
| `samlTokenUserAccountAuth` | `AuthResolverResponse!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `teamCreate` | `TeamPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `teamCyclesDelete` | `TeamPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `teamDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `teamKeyDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `teamMembershipCreate` | `TeamMembershipPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `teamMembershipDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `teamMembershipUpdate` | `TeamMembershipPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `teamUnarchive` | `TeamArchivePayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `teamUpdate` | `TeamPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `templateCreate` | `TemplatePayload!` | blocked_needs_design | mutation needs product and safety design |
| `templateDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `templateUpdate` | `TemplatePayload!` | blocked_needs_design | mutation needs product and safety design |
| `timeScheduleCreate` | `TimeSchedulePayload!` | blocked_needs_design | mutation needs product and safety design |
| `timeScheduleDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `timeScheduleRefreshIntegrationSchedule` | `TimeSchedulePayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `timeScheduleUpdate` | `TimeSchedulePayload!` | blocked_needs_design | mutation needs product and safety design |
| `timeScheduleUpsertExternal` | `TimeSchedulePayload!` | blocked_needs_design | mutation needs product and safety design |
| `trackAnonymousEvent` | `EventTrackingPayload!` | blocked_needs_design | mutation needs product and safety design |
| `triageResponsibilityCreate` | `TriageResponsibilityPayload!` | blocked_needs_design | mutation needs product and safety design |
| `triageResponsibilityDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `triageResponsibilityUpdate` | `TriageResponsibilityPayload!` | blocked_needs_design | mutation needs product and safety design |
| `updateIntegrationSlackScopes` | `IntegrationPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `userChangeRole` | `UserAdminPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userDiscordConnect` | `UserPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userExternalUserDisconnect` | `UserPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userFlagUpdate` | `UserSettingsFlagPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userRevokeAllSessions` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userRevokeSession` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userSettingsFlagsReset` | `UserSettingsFlagsResetPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userSettingsUpdate` | `UserSettingsPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userSuspend` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userUnlinkFromIdentityProvider` | `UserAdminPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `userUnsuspend` | `UserAdminPayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `userUpdate` | `UserPayload!` | accepted_gap | repo-planned or likely useful CLI domain |
| `viewPreferencesCreate` | `ViewPreferencesPayload!` | blocked_needs_design | mutation needs product and safety design |
| `viewPreferencesDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `viewPreferencesUpdate` | `ViewPreferencesPayload!` | blocked_needs_design | mutation needs product and safety design |
| `webhookCreate` | `WebhookPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `webhookDelete` | `DeletePayload!` | blocked_needs_design | destructive or access-changing operation needs explicit safety model |
| `webhookRotateSecret` | `WebhookRotateSecretPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `webhookUpdate` | `WebhookPayload!` | intentionally_excluded | admin/auth/internal integration surface outside ordinary agent CLI |
| `workflowStateArchive` | `WorkflowStateArchivePayload!` | blocked_needs_design | mutation needs product and safety design |
| `workflowStateCreate` | `WorkflowStatePayload!` | blocked_needs_design | mutation needs product and safety design |
| `workflowStateUpdate` | `WorkflowStatePayload!` | blocked_needs_design | mutation needs product and safety design |

## Local Generated Go Operations

| Operation | Kind | Root fields | Status | Evidence |
| --- | --- | --- | --- | --- |
| `AllTeamIssues` | query | `issues` | implemented | `internal/client/generated.go` |
| `CompletedWorkflowStates` | query | `workflowStates` | implemented | `internal/client/generated.go` |
| `CycleArchive` | mutation | `cycleArchive` | implemented | `internal/client/generated.go` |
| `CycleCreate` | mutation | `cycleCreate` | implemented | `internal/client/generated.go` |
| `CycleReport` | query | `cycle` | implemented | `internal/client/generated.go` |
| `CycleUpdate` | mutation | `cycleUpdate` | implemented | `internal/client/generated.go` |
| `CyclesByTeam` | query | `cycles` | implemented | `internal/client/generated.go` |
| `Documents` | query | `documents` | implemented | `internal/client/generated.go` |
| `IssueArchive` | mutation | `issueArchive` | implemented | `internal/client/generated.go` |
| `IssueBlockedIssues` | query | `issue` | implemented | `internal/client/generated.go` |
| `IssueClose` | mutation | `issueUpdate` | implemented | `internal/client/generated.go` |
| `IssueCommentCreate` | mutation | `commentCreate` | implemented | `internal/client/generated.go` |
| `IssueComments` | query | `issue` | implemented | `internal/client/generated.go` |
| `IssueCreate` | mutation | `issueCreate` | implemented | `internal/client/generated.go` |
| `IssueDependencies` | query | `issue` | implemented | `internal/client/generated.go` |
| `IssueLabels` | query | `issueLabels` | implemented | `internal/client/generated.go` |
| `IssueSearch` | query | `issues` | implemented | `internal/client/generated.go` |
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
| `ProjectMembers` | query | `project` | implemented | `internal/client/generated.go` |
| `ProjectMilestoneCreate` | mutation | `projectMilestoneCreate` | implemented | `internal/client/generated.go` |
| `ProjectMilestoneUpdate` | mutation | `projectMilestoneUpdate` | implemented | `internal/client/generated.go` |
| `ProjectMilestones` | query | `project` | implemented | `internal/client/generated.go` |
| `ProjectUpdate` | mutation | `projectUpdate` | implemented | `internal/client/generated.go` |
| `ProjectUpdates` | query | `project` | implemented | `internal/client/generated.go` |
| `Projects` | query | `team` | implemented | `internal/client/generated.go` |
| `StartedWorkflowStates` | query | `workflowStates` | implemented | `internal/client/generated.go` |
| `TargetProject` | query | `project` | implemented | `internal/client/generated.go` |
| `TeamMembers` | query | `team` | implemented | `internal/client/generated.go` |
| `Teams` | query | `teams` | implemented | `internal/client/generated.go` |
| `Viewer` | query | `viewer` | implemented | `internal/client/generated.go` |
| `comment` | query | `comment` | implemented | `internal/client/generated.go` |
| `comments` | query | `comments` | implemented | `internal/client/generated.go` |
| `cycle` | query | `cycle` | implemented | `internal/client/generated.go` |
| `document` | query | `document` | implemented | `internal/client/generated.go` |
| `issue` | query | `issue` | implemented | `internal/client/generated.go` |
| `issueLabel` | query | `issueLabel` | implemented | `internal/client/generated.go` |
| `project` | query | `project` | implemented | `internal/client/generated.go` |
| `projectMilestone` | query | `projectMilestone` | implemented | `internal/client/generated.go` |
| `projectUpdate` | query | `projectUpdate` | implemented | `internal/client/generated.go` |
| `projectUpdates` | query | `projectUpdates` | implemented | `internal/client/generated.go` |
| `team` | query | `team` | implemented | `internal/client/generated.go` |
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

