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

57. Doctor health check
   - Success: `linctl doctor` reports config load, token presence, and target confirmation without printing token values.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_execute_read_and_write_commands/doctor`;
     `Test_CommandFlows_print_json_for_read_and_comment_commands/--json/doctor`.

58. File-backed issue text
   - Success: `linctl issue create --description-file FILE`, `linctl issue update --append-file FILE`, `linctl issue comment --body-file FILE`, and `linctl issue reply --body-file FILE` read local file contents before the existing guarded write path.
   - Evidence: `go test ./internal/cli`, `Test_CommandFlows_read_issue_text_from_files`.

## Current Outcome

All fifty-eight local scenarios pass under the method above. The complete local suite also passes with `go test ./...`.

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
