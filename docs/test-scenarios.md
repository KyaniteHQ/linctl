# Test Scenarios

This file defines the repeatable scenario set for the coverage, logging, and live-smoke goal.

## Method

N is 5 for the local regression streak. Each local scenario runs under the
same unit-test conditions: fake GraphQL responses, no live Linear writes, and no
secret material in inputs or logs.

Success is pass/fail:

- The expected behavior is asserted by an automated test.
- Important failure paths produce actionable errors or diagnostic logs.
- Logs must not include Linear tokens, request bodies, response bodies, or fixture user data.
- The default suite remains fast enough for local iteration.
- Live smoke uses the same pass/fail rule, but runs only when a disposable
  Linear token is available.

## Scenarios

1. Target-pinned issue write
   - Success: creates, updates, comments, and closes only after target resolution and project/team checks.
   - Evidence: `go test ./internal/client`, `Test_ClientWriteScenarios_guard_writes_and_report_results`.

2. Target-pinned project write
   - Success: creates, updates, and archives only after target resolution and project/team checks.
   - Evidence: `go test ./internal/client`, `Test_ClientWriteScenarios_guard_writes_and_report_results`.

3. Read-only issue/project inspection
   - Success: list/get/member commands map generated GraphQL responses into compact models with pagination data.
   - Evidence: `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

4. Transport retry and diagnostics
   - Success: 429 responses retry, terminal failures remain errors, and diagnostics include attempt/status without secrets or bodies.
   - Evidence: `go test ./internal/client`, `Test_Transport_retries_429_with_retry_after_when_present`,
     `Test_Transport_logs_decode_failures_without_response_body`,
     `Test_Transport_logs_terminal_http_failures_without_response_body`.
   - Regression: request encoding failures used to return before writing any diagnostic event.
     `Test_Transport_returns_errors_for_request_and_body_failures/unmarshalable_variables`
     now asserts `graphql_encode_failed` is logged without request text.
   - Benchmark: `go test -run '^$' -bench Benchmark_Transport_make_request_diagnostics -benchmem ./internal/client`
     tracks the diagnostic writer cost with logging disabled and enabled.

5. Production-record classification boundary
   - Success: no repo-owned production-record dataset or classification logic is present; generated Linear schema references are external API surface, not linctl-owned production records.
   - Evidence: `rg -n "production records|classification|classify|record|allowed definition|production" README.md CONTEXT.md docs internal scripts test -S`.

6. Machine-readable CLI output controls
   - Success: CLI-level commands support compact JSON, JSON field projection, id-only output, quiet success output, fail-on-empty list behavior, deterministic sort/order, and minimal human format.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_print_compact_json_when_compact_flag_is_set`,
     `Test_CommandFlows_project_json_fields_when_fields_flag_is_set`,
     `Test_CommandFlows_print_only_id_when_id_only_flag_is_set`,
     `Test_CommandFlows_suppress_success_output_when_quiet_flag_is_set`,
     `Test_CommandFlows_fail_on_empty_list_when_fail_on_empty_flag_is_set`,
     `Test_CommandFlows_sort_issue_list_when_sort_flags_are_set`,
     `Test_CommandFlows_print_minimal_human_output_when_format_flag_is_set`.

7. Current Issue branch helpers
   - Success: branch-derived issue helpers print the Current Issue identifier, title, URL, and Linear branch name through the public CLI surface.
   - Evidence: `go test ./internal/cli`,
     `Test_CommandFlows_print_current_issue_identifier_from_issue_id`,
     `Test_CommandFlows_print_current_issue_title_from_issue_title`,
     `Test_CommandFlows_print_current_issue_url_from_issue_url`,
     `Test_CommandFlows_print_issue_branch_from_issue_branch`.

8. Stdin comment body
   - Success: `linctl issue comment ISSUE --body -` reads the full command stdin as the guarded comment body.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_read_issue_comment_body_from_stdin`.

9. Issue comment list
   - Success: `linctl issue comments ISSUE --limit N` lists issue discussion comments through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_comments`.

10. Issue state filter
   - Success: `linctl issue list --state started --limit N` lists issues whose workflow state type is `started`.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_state_filter`.

11. Issue text search
   - Success: `linctl issue search QUERY --limit N` lists matching issues in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_search`.

12. Issue project filter
   - Success: `linctl issue list --project PROJECT_ID --limit N` lists issues attached to that project in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_project_filter`.

13. Issue mine filter
   - Success: `linctl issue list --mine --limit N` lists issues assigned to the authenticated user in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_mine_filter`.

14. Issue assignee filter
   - Success: `linctl issue list --assignee USER_ID --limit N` lists issues assigned to that Linear user id in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_assignee_filter`.

15. Issue label filter
   - Success: `linctl issue list --label LABEL_ID --limit N` lists issues carrying that Linear label id in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_label_filter`.

16. Issue cycle filter
   - Success: `linctl issue list --cycle CYCLE_ID --limit N` lists issues attached to that Cycle in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_cycle_filter`.

17. Issue created-after filter
   - Success: `linctl issue list --created-after DATE --limit N` lists issues created on or after that date in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_created-after_filter`.

18. Issue created-before filter
   - Success: `linctl issue list --created-before DATE --limit N` lists issues created on or before that date in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_created-before_filter`.

19. Issue created-since filter
   - Success: `linctl issue list --created-since DATE --limit N` lists issues created on or after that date in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_created-since_filter`.

20. Issue all-teams list
   - Success: `linctl issue list --all-teams --limit N` lists issues across every visible Linear team without applying the resolved-team filter.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_all_teams`.

21. Issue has-blockers filter
   - Success: `linctl issue list --has-blockers --limit N` lists issues blocked by another issue in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_has_blockers_filter`.

22. Issue blocks filter
   - Success: `linctl issue list --blocks --limit N` lists issues blocking another issue in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_blocks_filter`.

23. Issue blocked-by filter
   - Success: `linctl issue list --blocked-by ISSUE --limit N` lists issues blocked by that issue in the resolved team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_list_blocked_by_filter`.

24. Issue dependency graph
   - Success: `linctl issue deps ISSUE --limit N` lists the issue's parent, children, blocked issues, and blockers.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_deps`.

25. Issue comment reply
   - Success: `linctl issue reply ISSUE COMMENT --body TEXT` creates a threaded reply through the guarded issue comment path.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_reply`.

26. Issue start
   - Success: `linctl issue start ISSUE` assigns the issue to the authenticated viewer and moves it to the team's first started workflow state through the guarded issue update path.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_start`;
     `go test ./internal/client`, `Test_StartIssue_assigns_viewer_and_moves_to_started_state_when_target_matches`.

27. Done current issue
   - Success: `linctl done` derives the current checkout issue and closes it through the guarded issue close path.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_close_current_issue_from_done`.

28. Issue PR plan
   - Success: `linctl issue pr [ISSUE]` reads an explicit or current checkout issue and prints a `gh pr create` title/body plan without calling GitHub.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_pr`,
     `Test_CommandFlows_print_issue_pr_from_current_branch`.

29. Issue description append
   - Success: `linctl issue update ISSUE --append TEXT` preserves the existing description and appends the text through the guarded issue update path.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_update_append`;
     `go test ./internal/client`, `Test_ClientWriteScenarios_guard_writes_and_report_results/issue_update_appends_to_description`.

30. Next dry-run issue picker
   - Success: `linctl next --dry-run` resolves the pinned target, reads unstarted issues with no blocking relations, ranks candidates by active unblock count, priority, then created date, and prints the selected candidate without creating a branch or worktree.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/next_dry_run`;
     `Test_CommandFlows_rank_next_issue_candidates`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

31. Project update history
   - Success: `linctl project updates PROJECT --limit N` lists project status updates with health, author, and body through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_updates`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

32. Project milestone list
   - Success: `linctl project-milestone list PROJECT --limit N` lists a project's milestones with status through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_milestone_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

33. Cycle list
   - Success: `linctl cycle list --limit N` lists Cycles for the resolved team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CycleCommandFlows_list_cycles`;
     `go test ./internal/client`, `Test_ListCyclesByTeam_returns_cycle_page`.

34. Cycle get
   - Success: `linctl cycle get CYCLE_ID` reads one Cycle by id or slug through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CycleCommandFlows_get_cycle`;
     `go test ./internal/client`, `Test_GetCycleByID_returns_cycle`.

35. ProjectMilestone get
   - Success: `linctl project-milestone get PROJECT_MILESTONE_ID` reads one ProjectMilestone by id through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_get_project_milestone`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

36. Sprint current
   - Success: `linctl sprint current` reports the active Cycle for the resolved team without exposing Sprint mutations.
   - Evidence: `go test ./internal/cli`, `Test_CycleCommandFlows_get_current_sprint`;
     `go test ./internal/client`, `Test_CurrentCycleByTeam_returns_active_cycle`.

37. Sprint report
   - Success: `linctl sprint report CYCLE_ID --limit N` reports one Cycle and its assigned issues without exposing Sprint mutations.
   - Evidence: `go test ./internal/cli`, `Test_CycleCommandFlows_report_sprint`;
     `go test ./internal/client`, `Test_GetSprintReport_returns_cycle_and_issues`.

38. ProjectMilestone create
   - Success: `linctl project-milestone create PROJECT_ID --name NAME` creates a ProjectMilestone only after resolving and comparing the target project.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_milestone_create`;
     `go test ./internal/client`, `Test_CreateProjectMilestone_returns_created_milestone_when_target_matches`.

39. ProjectMilestone update
   - Success: `linctl project-milestone update PROJECT_MILESTONE_ID --name NAME` updates a ProjectMilestone only after resolving the ProjectMilestone's project and comparing the pinned target.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_milestone_update`;
     `go test ./internal/client`, `Test_UpdateProjectMilestone_returns_updated_milestone_when_target_matches`,
     `Test_UpdateProjectMilestone_refuses_when_pinned_project_differs`,
     `Test_UpdateProjectMilestone_refuses_when_project_team_differs`.

40. Cycle create
   - Success: `linctl cycle create --starts-at START --ends-at END` creates a Cycle only in the pinned team.
   - Evidence: `go test ./internal/cli`, `Test_CycleCommandFlows_write_cycles/create`;
     `go test ./internal/client`, `Test_CreateCycle_returns_created_cycle_when_target_matches`.

41. Cycle update
   - Success: `linctl cycle update CYCLE_ID --name NAME` updates a Cycle only after resolving and comparing its team.
   - Evidence: `go test ./internal/cli`, `Test_CycleCommandFlows_write_cycles/update`;
     `go test ./internal/client`, `Test_UpdateCycle_returns_updated_cycle_when_target_matches`,
     `Test_UpdateCycle_refuses_when_team_differs`.

42. Cycle archive
   - Success: `linctl cycle archive CYCLE_ID` archives a Cycle only after resolving and comparing its team.
   - Evidence: `go test ./internal/cli`, `Test_CycleCommandFlows_write_cycles/archive`;
     `go test ./internal/client`, `Test_ArchiveCycle_returns_archived_cycle_when_target_matches`.

43. Document list
   - Success: `linctl document list --limit N` lists visible Documents with parent type metadata through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/document_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

44. Document get
   - Success: `linctl document get DOCUMENT_ID` reads one Document by id or slug and reports its resolved parent type.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/document_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

45. Label list
   - Success: `linctl label list --limit N` lists visible Linear IssueLabels through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/label_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

46. Label get
   - Success: `linctl label get LABEL_ID` reads one Linear IssueLabel by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/label_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

47. Team list
   - Success: `linctl team list --limit N` lists visible teams through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

48. Team get
   - Success: `linctl team get TEAM_ID` reads one Team by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

49. Team members
   - Success: `linctl team members TEAM_ID --limit N` lists users on one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_members`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

50. User list
   - Success: `linctl user list --limit N` lists visible users through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

51. User get
   - Success: `linctl user get USER_ID` reads one Linear User by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

52. User me
   - Success: `linctl user me` reads the authenticated Linear User through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_me`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

53. WorkflowState list
   - Success: `linctl workflow-state list --limit N` lists visible WorkflowStates through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/workflow_state_list`;
     `Test_CommandFlows_print_workflow_state_list_as_json`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

54. WorkflowState get
   - Success: `linctl workflow-state get WORKFLOW_STATE_ID` reads one WorkflowState by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/workflow_state_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

55. Comment list
   - Success: `linctl comment list --limit N` lists visible comments through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/comment_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

56. Comment get
   - Success: `linctl comment get COMMENT_ID` reads one Linear Comment by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/comment_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

57. ProjectUpdate list
   - Success: `linctl project-update list --limit N` lists visible ProjectUpdates through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_update_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

58. ProjectUpdate get
   - Success: `linctl project-update get PROJECT_UPDATE_ID` reads one ProjectUpdate by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_update_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

59. InitiativeUpdate list
   - Success: `linctl initiative-update list --limit N` lists visible InitiativeUpdates through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_update_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

60. InitiativeUpdate get
   - Success: `linctl initiative-update get INITIATIVE_UPDATE_ID` reads one InitiativeUpdate by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_update_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

61. InitiativeRelation list
   - Success: `linctl initiative-relation list --limit N` lists visible InitiativeRelations through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_relation_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

62. InitiativeRelation get
   - Success: `linctl initiative-relation get INITIATIVE_RELATION_ID` reads one InitiativeRelation by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_relation_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

63. InitiativeToProject list
   - Success: `linctl initiative-to-project list --limit N` lists visible InitiativeToProjects through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_to_project_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

64. InitiativeToProject get
   - Success: `linctl initiative-to-project get INITIATIVE_TO_PROJECT_ID` reads one InitiativeToProject by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_to_project_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

65. Doctor health check
   - Success: `linctl doctor` reports config load, token presence, and target confirmation without printing token values.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/doctor`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands/--json/doctor`.

66. File-backed issue text
   - Success: `linctl issue create --description-file FILE`, `linctl issue update --append-file FILE`, `linctl issue comment --body-file FILE`, and `linctl issue reply --body-file FILE` read local file contents before the existing guarded write path.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_read_issue_text_from_files`.

67. Initiative list
   - Success: `linctl initiative list --limit N` lists visible Initiatives through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

68. Initiative get
   - Success: `linctl initiative get INITIATIVE_ID` reads one Initiative by id or slug.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

73. Initiative history
   - Success: `linctl initiative history INITIATIVE_ID --limit N` lists history records associated with one Initiative through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_history`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

74. Initiative links
   - Success: `linctl initiative links INITIATIVE_ID --limit N` lists external links associated with one Initiative through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_links`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

75. Initiative sub-initiatives
   - Success: `linctl initiative sub-initiatives INITIATIVE_ID --limit N` lists child Initiatives associated with one Initiative through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_sub-initiatives`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

76. Initiative updates
   - Success: `linctl initiative updates INITIATIVE_ID --limit N` lists InitiativeUpdates associated with one Initiative through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/initiative_updates`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

77. CustomView list
   - Success: `linctl custom-view list --limit N` lists visible CustomViews through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

78. CustomView get
   - Success: `linctl custom-view get CUSTOM_VIEW_ID` reads one CustomView by id or slug.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

79. CustomView subscribers
   - Success: `linctl custom-view subscribers CUSTOM_VIEW_ID` reports whether a CustomView has active notification subscribers through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_subscribers`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

80. CustomView initiatives
   - Success: `linctl custom-view initiatives CUSTOM_VIEW_ID --limit N` lists Initiatives matching a CustomView initiative filter through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_initiatives`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

81. CustomView issues
   - Success: `linctl custom-view issues CUSTOM_VIEW_ID --limit N` lists Issues matching a CustomView issue filter through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_issues`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

82. CustomView organization preferences
   - Success: `linctl custom-view organization-preferences CUSTOM_VIEW_ID` reads organization default ViewPreferences for one CustomView through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_organization_preferences`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

83. CustomView organization preference values
   - Success: `linctl custom-view organization-preferences values CUSTOM_VIEW_ID` reads organization default ViewPreferencesValues for one CustomView through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_organization_preference_values`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

84. CustomView projects
   - Success: `linctl custom-view projects CUSTOM_VIEW_ID --limit N` lists Projects matching a CustomView project filter through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_projects`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

85. CustomView user preferences
   - Success: `linctl custom-view user-preferences CUSTOM_VIEW_ID` reads the current user's ViewPreferences for one CustomView through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_user_preferences`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

86. CustomView user preference values
   - Success: `linctl custom-view user-preferences values CUSTOM_VIEW_ID` reads the current user's ViewPreferencesValues for one CustomView through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_user_preference_values`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

87. CustomView preference values
   - Success: `linctl custom-view preference-values CUSTOM_VIEW_ID` reads effective ViewPreferencesValues for one CustomView through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/custom_view_preference_values`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

84. Favorite list
   - Success: `linctl favorite list --limit N` lists the authenticated user's Favorites through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/favorite_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

85. Favorite get
   - Success: `linctl favorite get FAVORITE_ID` reads one Favorite by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/favorite_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

86. Favorite children
   - Success: `linctl favorite children FAVORITE_ID --limit N` lists child Favorites under a folder Favorite through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/favorite_children`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

87. Emoji list
   - Success: `linctl emoji list --limit N` lists workspace custom Emojis through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/emoji_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

88. Emoji get
   - Success: `linctl emoji get EMOJI_ID` reads one custom Emoji by id or name.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/emoji_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

89. Attachment list
   - Success: `linctl attachment list --limit N` lists visible issue Attachments through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/attachment_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

90. Attachment get
   - Success: `linctl attachment get ATTACHMENT_ID` reads one Attachment by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/attachment_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

91. Attachment URL lookup
   - Success: `linctl attachment url URL --limit N` lists issue Attachments linked to a URL through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/attachment_url`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

92. Organization exists
   - Success: `linctl organization exists URL_KEY` reports whether a Linear organization URL key exists through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/organization_exists`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

93. Organization labels
   - Success: `linctl organization labels --limit N` lists workspace-level issue labels through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/organization_labels`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

94. Organization project labels
   - Success: `linctl organization project-labels --limit N` lists workspace-level project labels through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/organization_project_labels`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

95. Organization teams
   - Success: `linctl organization teams --limit N` lists workspace teams visible to the authenticated user through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/organization_teams`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

96. Organization users
   - Success: `linctl organization users --limit N` lists active workspace users visible to the authenticated user through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/organization_users`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

97. Organization templates
   - Success: `linctl organization templates --limit N` lists workspace-level Linear Templates through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/organization_templates`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

98. Rate-limit status
   - Success: `linctl rate-limit status` reports the authenticated Linear client's current quota buckets through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/rate_limit_status`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

99. Customer list
   - Success: `linctl customer list --limit N` lists visible Linear Customers through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/customer_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

100. Customer get
   - Success: `linctl customer get CUSTOMER_ID` reads one Linear Customer by id or slug.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/customer_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

101. Customer need list
   - Success: `linctl customer-need list --limit N` lists visible Linear CustomerNeeds through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/customer_need_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

102. Customer need get
   - Success: `linctl customer-need get CUSTOMER_NEED_ID` reads one Linear CustomerNeed by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/customer_need_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

103. Customer status list
   - Success: `linctl customer-status list --limit N` lists workspace CustomerStatuses through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/customer_status_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

104. Customer status get
   - Success: `linctl customer-status get CUSTOMER_STATUS_ID` reads one Linear CustomerStatus by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/customer_status_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

105. Customer tier list
   - Success: `linctl customer-tier list --limit N` lists workspace CustomerTiers through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/customer_tier_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

106. Customer tier get
   - Success: `linctl customer-tier get CUSTOMER_TIER_ID` reads one Linear CustomerTier by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/customer_tier_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

107. Roadmap list
   - Success: `linctl roadmap list --limit N` lists visible Linear Roadmaps through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/roadmap_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

108. Roadmap get
   - Success: `linctl roadmap get ROADMAP_ID` reads one Linear Roadmap by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/roadmap_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

109. Time schedule list
   - Success: `linctl time-schedule list --limit N` lists visible Linear TimeSchedules through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/time_schedule_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

110. Time schedule get
   - Success: `linctl time-schedule get TIME_SCHEDULE_ID` reads one Linear TimeSchedule by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/time_schedule_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

111. Template list
   - Success: `linctl template list --limit N` lists visible Linear Templates through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/template_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

112. Template get
   - Success: `linctl template get TEMPLATE_ID` reads one Linear Template by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/template_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

113. Notification list
   - Success: `linctl notification list --limit N` lists authenticated-user Notifications through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/notification_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

114. Notification get
   - Success: `linctl notification get NOTIFICATION_ID` reads one Notification by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/notification_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

115. Notification subscription list
   - Success: `linctl notification subscription list --limit N` lists authenticated-user NotificationSubscriptions through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/notification_subscription_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

116. Notification subscription get
   - Success: `linctl notification subscription get NOTIFICATION_SUBSCRIPTION_ID` reads one NotificationSubscription by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/notification_subscription_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

117. Release pipeline list
   - Success: `linctl release-pipeline list --limit N` lists visible Linear ReleasePipelines through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_pipeline_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

118. Release pipeline get
   - Success: `linctl release-pipeline get RELEASE_PIPELINE_ID` reads one Linear ReleasePipeline by id or slug.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_pipeline_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

119. Release pipeline releases
   - Success: `linctl release-pipeline releases RELEASE_PIPELINE_ID --limit N` lists Releases associated with one ReleasePipeline through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_pipeline_releases`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

120. Release pipeline stages
   - Success: `linctl release-pipeline stages RELEASE_PIPELINE_ID --limit N` lists ReleaseStages associated with one ReleasePipeline through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_pipeline_stages`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

121. Release stage list
   - Success: `linctl release-stage list --limit N` lists visible Linear ReleaseStages through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_stage_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

122. Release stage get
   - Success: `linctl release-stage get RELEASE_STAGE_ID` reads one Linear ReleaseStage by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_stage_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

123. Release stage releases
   - Success: `linctl release-stage releases RELEASE_STAGE_ID --limit N` lists Releases associated with one ReleaseStage through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_stage_releases`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

124. Release list
   - Success: `linctl release list --limit N` lists visible Linear Releases through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

125. Release search
   - Success: `linctl release search TERM --limit N` searches Linear Releases by text through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_search`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

126. Release get
   - Success: `linctl release get RELEASE_ID` reads one Linear Release by id or slug.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

127. Release history
   - Success: `linctl release history RELEASE_ID --limit N` lists history records associated with one Release through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_history`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

128. Release links
   - Success: `linctl release links RELEASE_ID --limit N` lists external links associated with one Release through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_links`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

129. External link get
   - Success: `linctl external-link get EXTERNAL_LINK_ID` reads one Linear EntityExternalLink by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/external_link_get`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

130. Release note list
   - Success: `linctl release-note list --limit N` lists visible Linear ReleaseNotes through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_note_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

131. Release note get
   - Success: `linctl release-note get RELEASE_NOTE_ID` reads one Linear ReleaseNote by id or slug.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/release_note_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

132. Application info
   - Success: `linctl application info CLIENT_ID` reads public Linear OAuth Application metadata by client id through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/application_info`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

133. AgentSkill list
   - Success: `linctl agent-skill list --limit N` lists visible Linear AgentSkills through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/agent_skill_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

134. AgentSkill get
   - Success: `linctl agent-skill get AGENT_SKILL_ID` reads one Linear AgentSkill by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/agent_skill_get`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

135. AgentActivity list
   - Success: `linctl agent-activity list --limit N` lists visible Linear AgentActivities through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/agent_activity_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

136. AgentActivity get
   - Success: `linctl agent-activity get AGENT_ACTIVITY_ID` reads one Linear AgentActivity by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/agent_activity_get`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

137. AuditEntry types
   - Success: `linctl audit-entry types` lists Linear audit entry type names and descriptions without audit log metadata.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/audit_entry_types`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

138. TriageResponsibility list
   - Success: `linctl triage-responsibility list --limit N` lists visible Linear TriageResponsibility records.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/triage_responsibility_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

139. TriageResponsibility get
   - Success: `linctl triage-responsibility get TRIAGE_RESPONSIBILITY_ID` reads one Linear TriageResponsibility by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/triage_responsibility_get`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

140. TriageResponsibility manual selection
   - Success: `linctl triage-responsibility manual-selection TRIAGE_RESPONSIBILITY_ID` reads the manual user selection for one Linear TriageResponsibility.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/triage_responsibility_manual_selection`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

141. SLA configuration list
   - Success: `linctl sla-configuration list TEAM_ID_OR_KEY` lists active Linear SLA configurations that can apply to a team.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/SLA_configuration_list`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

142. Semantic search
   - Success: `linctl semantic-search QUERY --limit N` returns compact references from Linear semantic search.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/semantic_search`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

143. User drafts
   - Success: `linctl user drafts --limit N` lists authenticated-user draft parent metadata without draft body/data.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_drafts`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

144. User assigned issues
   - Success: `linctl user assigned-issues USER_ID --limit N` lists issues assigned to one User without issue body content.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_assigned_issues`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

145. User created issues
   - Success: `linctl user created-issues USER_ID --limit N` lists issues created by one User without issue body content.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_created_issues`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

146. User delegated issues
   - Success: `linctl user delegated-issues USER_ID --limit N` lists issues delegated to one User without issue body content.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_delegated_issues`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

147. User team memberships
   - Success: `linctl user team-memberships USER_ID --limit N` lists TeamMembership summaries for one User.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_team_memberships`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

148. User teams
   - Success: `linctl user teams USER_ID --limit N` lists Teams for one User.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_teams`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

149. User my assigned issues
   - Success: `linctl user my-assigned-issues --limit N` lists issues assigned to the authenticated User without issue body content.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_my_assigned_issues`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

150. User my created issues
   - Success: `linctl user my-created-issues --limit N` lists issues created by the authenticated User without issue body content.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_my_created_issues`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

151. User my delegated issues
   - Success: `linctl user my-delegated-issues --limit N` lists issues delegated to the authenticated User without issue body content.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_my_delegated_issues`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

152. User my team memberships
   - Success: `linctl user my-team-memberships --limit N` lists TeamMembership summaries for the authenticated User.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_my_team_memberships`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

153. User my teams
   - Success: `linctl user my-teams --limit N` lists Teams for the authenticated User.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/user_my_teams`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

154. ProjectStatus list
   - Success: `linctl project-status list --limit N` lists visible ProjectStatuses through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_status_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

155. ProjectStatus get
   - Success: `linctl project-status get PROJECT_STATUS_ID` reads one ProjectStatus by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_status_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

156. ProjectLabel list
   - Success: `linctl project-label list --limit N` lists visible ProjectLabels through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_label_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

157. ProjectLabel get
   - Success: `linctl project-label get PROJECT_LABEL_ID` reads one ProjectLabel by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_label_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

158. ProjectLabel children
   - Success: `linctl project-label children PROJECT_LABEL_ID --limit N` lists child ProjectLabels without exposing unrelated project data.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_label_children`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

159. ProjectLabel projects
   - Success: `linctl project-label projects PROJECT_LABEL_ID --limit N` lists projects associated with one ProjectLabel.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_label_projects`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

160. ProjectRelation list
   - Success: `linctl project-relation list --limit N` lists visible ProjectRelations through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_relation_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

161. ProjectRelation get
   - Success: `linctl project-relation get PROJECT_RELATION_ID` reads one ProjectRelation by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/project_relation_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

162. TeamMembership list
   - Success: `linctl team-membership list --limit N` lists visible TeamMemberships through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_membership_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

163. TeamMembership get
   - Success: `linctl team-membership get TEAM_MEMBERSHIP_ID` reads one TeamMembership by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_membership_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

164. RoadmapToProject list
   - Success: `linctl roadmap-to-project list --limit N` lists visible RoadmapToProjects through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/roadmap_to_project_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

165. RoadmapToProject get
   - Success: `linctl roadmap-to-project get ROADMAP_TO_PROJECT_ID` reads one RoadmapToProject by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/roadmap_to_project_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

166. IssueRelation list
   - Success: `linctl issue-relation list --limit N` lists visible IssueRelations through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_relation_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

167. IssueRelation get
   - Success: `linctl issue-relation get ISSUE_RELATION_ID` reads one IssueRelation by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_relation_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

168. IssueToRelease list
   - Success: `linctl issue-to-release list --limit N` lists visible IssueToReleases through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_to_release_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

169. IssueToRelease get
   - Success: `linctl issue-to-release get ISSUE_TO_RELEASE_ID` reads one IssueToRelease by id.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/issue_to_release_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

170. ExternalUser list
   - Success: `linctl external-user list --limit N` lists visible Linear ExternalUsers without selecting email through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/external_user_list`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

171. ExternalUser get
   - Success: `linctl external-user get EXTERNAL_USER_ID` reads one Linear ExternalUser by id without selecting email.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/external_user_get`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

172. Team cycles
   - Success: `linctl team cycles TEAM_ID --limit N` lists Cycles associated with one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_cycles`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

173. Team issues
   - Success: `linctl team issues TEAM_ID --limit N` lists Issues associated with one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_issues`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

174. Team labels
   - Success: `linctl team labels TEAM_ID --limit N` lists IssueLabels associated with one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_labels`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

175. Team memberships
   - Success: `linctl team memberships TEAM_ID --limit N` lists TeamMemberships associated with one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_memberships`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

176. Team projects
   - Success: `linctl team projects TEAM_ID --limit N` lists Projects associated with one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_projects`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

177. Team release pipelines
   - Success: `linctl team release-pipelines TEAM_ID --limit N` lists ReleasePipelines associated with one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_release_pipelines`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

178. Team states
   - Success: `linctl team states TEAM_ID --limit N` lists WorkflowStates associated with one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_states`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

179. Team templates
   - Success: `linctl team templates TEAM_ID --limit N` lists Templates associated with one Team through the public CLI and JSON output controls.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/team_templates`;
     `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

## Current Outcome

All one hundred seventy-one local scenarios pass under the method above. The complete local suite also passes with `go test ./...`.

Coverage is enforced with `task coverage`, which runs uncached tests and excludes generated GraphQL code, the thin process entrypoint, and repo maintenance scripts from the product behavior metric. The enforced product statement coverage target is 100.0%.

## Live Smoke

Run the complete live smoke suite with:

```bash
task live-smoke
```

The command requires a disposable Linear token in `LINCTL_TEST_TOKEN`, `LINCTL_TOKEN`, or `LINEAR_API_KEY`.
It builds a temporary `linctl` binary, smoke-tests read-only CLI commands, then runs the integration-tagged
client round trips. Write smoke tests create `linctl-it-<runid>` resources and archive them during cleanup.

The missing-token readiness check is `env -u LINCTL_TEST_TOKEN -u LINCTL_TOKEN -u LINEAR_API_KEY bash scripts/live-smoke.sh`,
which must exit 2 with a missing-token message and without printing secret values.
