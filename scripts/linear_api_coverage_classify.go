//go:build ignore

package main

import (
	"strings"

	"github.com/KyaniteHQ/linctl/internal/cli"
)

var blockedDomainCommands = map[string]bool{
	"document create":                                   true,
	"document update":                                   true,
	"comment resolve":                                   true,
	"comment unresolve":                                 true,
	"issue-relation create":                             true,
	"issue-relation update":                             true,
	"issue-relation delete":                             true,
	"project-update create":                             true,
	"project-update update":                             true,
	"project-update archive":                            true,
	"project-status create":                             true,
	"project-status update":                             true,
	"project-status archive":                            true,
	"project-status unarchive":                          true,
	"project-label create":                              true,
	"project-label update":                              true,
	"project-label delete":                              true,
	"project-label retire":                              true,
	"project-label restore":                             true,
	"project-relation create":                           true,
	"project-relation update":                           true,
	"project-relation delete":                           true,
	"label create":                                      true,
	"label update":                                      true,
	"team create":                                       true,
	"team update":                                       true,
	"team delete":                                       true,
	"team-membership create":                            true,
	"team-membership update":                            true,
	"team-membership delete":                            true,
	"workflow-state create":                             true,
	"workflow-state update":                             true,
	"workflow-state archive":                            true,
	"time-schedule create":                              true,
	"time-schedule update":                              true,
	"time-schedule delete":                              true,
	"time-schedule upsert-external":                     true,
	"template create":                                   true,
	"template update":                                   true,
	"template delete":                                   true,
	"initiative-relation create":                        true,
	"initiative-relation update":                        true,
	"initiative-relation delete":                        true,
	"initiative-to-project create":                      true,
	"initiative-to-project update":                      true,
	"initiative-to-project delete":                      true,
	"roadmap-to-project create":                         true,
	"roadmap-to-project update":                         true,
	"roadmap-to-project delete":                         true,
	"initiative-update create":                          true,
	"initiative-update update":                          true,
	"initiative-update archive":                         true,
	"initiative-update unarchive":                       true,
	"initiative create":                                 true,
	"initiative update":                                 true,
	"initiative archive":                                true,
	"roadmap create":                                    true,
	"roadmap update":                                    true,
	"roadmap archive":                                   true,
	"roadmap delete":                                    true,
	"custom-view create":                                true,
	"custom-view update":                                true,
	"customer create":                                   true,
	"customer update":                                   true,
	"customer archive":                                  true,
	"customer-need create":                              true,
	"customer-need update":                              true,
	"customer-need archive":                             true,
	"customer-need delete":                              true,
	"customer-status create":                            true,
	"customer-status update":                            true,
	"customer-status delete":                            true,
	"customer-tier create":                              true,
	"customer-tier update":                              true,
	"customer-tier delete":                              true,
	"favorite create":                                   true,
	"favorite update":                                   true,
	"emoji create":                                      true,
	"attachment create":                                 true,
	"attachment update":                                 true,
	"notification archive":                              true,
	"notification archive all":                          true,
	"notification update":                               true,
	"notification mark read all":                        true,
	"notification mark unread all":                      true,
	"notification snooze all":                           true,
	"notification unsnooze all":                         true,
	"notification category channel subscription update": true,
	"notification subscription create":                  true,
	"notification subscription update":                  true,
	"notification subscription delete":                  true,
	"release-pipeline create":                           true,
	"release-pipeline update":                           true,
	"release-pipeline archive":                          true,
	"release-pipeline unarchive":                        true,
	"release-pipeline delete":                           true,
	"release-stage create":                              true,
	"release-stage update":                              true,
	"release-stage archive":                             true,
	"release-stage unarchive":                           true,
	"release create":                                    true,
	"release update":                                    true,
	"release archive":                                   true,
	"release unarchive":                                 true,
	"release delete":                                    true,
	"release complete":                                  true,
	"release sync":                                      true,
	"release-note create":                               true,
	"release-note update":                               true,
	"release-note archive":                              true,
	"release-note delete":                               true,
	"issue-to-release create":                           true,
	"issue-to-release update":                           true,
	"issue-to-release delete":                           true,
}

func domainCommandBlocked(command string) bool {
	return blockedDomainCommands[command]
}

func classifySDKMethod(
	name string,
	commandInventory map[string]cli.CommandInfo,
	implementedRoots map[string]bool,
	rootKinds map[string]string,
) (string, string) {
	if sdkImplemented(name, implementedRoots) {
		return "generated_operation", "local GraphQL operation uses this root"
	}
	if command, ok := commandInventory[commandLookupName(name)]; ok && command.OperationBacking != "" {
		if command.Safety == cli.CommandSafetyWrite {
			return "guarded_write_command", "public CLI command exposes this operation with write safety metadata"
		}
		return "public_command", "public CLI command exposes this operation"
	}
	if kind, ok := sdkRootKind(name, rootKinds); ok {
		return classifyLoose(name, kind)
	}
	if status, rationale, ok := explicitRiskClassification(strings.ToLower(name)); ok {
		return status, rationale
	}

	return "blocked_needs_design", "SDK method is not matched to a GraphQL root field; explicit classification required"
}

func classifyRoot(field rootField, implementedRoots map[string]bool) (string, string) {
	if implementedRoots[rootKey(field.Kind, field.Name)] {
		return "generated_operation", "root field used by local GraphQL operation"
	}
	return classifyLoose(field.Name, field.Kind)
}

func classifyLoose(name string, kind string) (string, string) {
	lower := strings.ToLower(name)
	if status, rationale, ok := explicitRiskClassification(lower); ok {
		return status, rationale
	}
	switch {
	case strings.Contains(lower, "latestreleasebyaccesskey"),
		strings.Contains(lower, "releasepipelinebyaccesskey"):
		return "intentionally_excluded", accessKeyReleaseRationale()
	case strings.Contains(lower, "documentcontent"),
		strings.Contains(lower, "archivepayload"),
		strings.Contains(lower, "externalthread"):
		return "blocked_needs_design", contentPayloadReadRationale()
	case strings.Contains(lower, "delete"),
		strings.Contains(lower, "remove"),
		strings.Contains(lower, "revoke"),
		strings.Contains(lower, "suspend"):
		return "blocked_needs_design", "destructive or access-changing operation needs explicit safety model"
	case strings.Contains(lower, "admin"),
		strings.Contains(lower, "auth"),
		strings.Contains(lower, "oauth"),
		strings.Contains(lower, "session"),
		strings.Contains(lower, "webhook"),
		strings.Contains(lower, "integration"):
		return "intentionally_excluded", "admin/auth/internal integration surface outside ordinary agent CLI"
	case hasWritePrefix(lower):
		return "blocked_needs_design", "write operation needs guarded target semantics before exposure"
	case strings.Contains(lower, "resolve"):
		return "blocked_needs_design", "state-changing operation needs guarded target semantics before exposure"
	case strings.Contains(lower, "issue"),
		strings.Contains(lower, "project"),
		strings.Contains(lower, "cycle"),
		strings.Contains(lower, "document"),
		strings.Contains(lower, "label"),
		strings.Contains(lower, "team"),
		strings.Contains(lower, "user"),
		strings.Contains(lower, "comment"):
		return "accepted_gap", "repo-planned or likely useful CLI domain"
	default:
		if kind == "mutation" {
			return "blocked_needs_design", "mutation needs product and safety design"
		}
		return "safe_candidate", "read operation may fit future CLI coverage"
	}
}

type riskClassification struct {
	status    string
	rationale string
}

var explicitRiskClassifications = map[string]riskClassification{
	"auditentries": {
		status: "blocked_needs_design",
		rationale: "audit logs can expose actor, IP, country, and request metadata; " +
			"needs explicit admin/security output model",
	},
	"emailintakeaddress": {
		status:    "intentionally_excluded",
		rationale: "email intake administration sits outside the ordinary agent CLI read surface",
	},
	"emailintakeaddress_sesdomainidentity": {
		status:    "intentionally_excluded",
		rationale: "email domain identity administration sits outside the ordinary agent CLI read surface",
	},
	"attachmentlinkgithubissue": {
		status: "blocked_needs_design",
		rationale: "attachment-to-GitHub linking mutates third-party integration state; " +
			"needs explicit integration guard semantics",
	},
	"attachmentlinkjiraissue": {
		status: "blocked_needs_design",
		rationale: "attachment-to-Jira linking mutates third-party integration state; " +
			"needs explicit integration guard semantics",
	},
	"availableusers": {
		status: "intentionally_excluded",
		rationale: "available-user picker enumeration is a specialized product resolver; " +
			"`user list` is the supported user read surface",
	},
	"cycleshiftall": {
		status: "blocked_needs_design",
		rationale: "bulk Cycle date shifting is a state-changing organization operation that " +
			"needs target-pinned guard semantics",
	},
	"cyclestartupcomingcycletoday": {
		status: "blocked_needs_design",
		rationale: "starting an upcoming Cycle changes team planning state and needs " +
			"target-pinned guard semantics",
	},
	"issueaddlabel": {
		status:    "blocked_needs_design",
		rationale: "issue label mutation needs issue target pinning and target-mismatch tests",
	},
	"issueexternalsyncdisable": {
		status: "blocked_needs_design",
		rationale: "issue external-sync disable changes integration state and needs explicit " +
			"integration guard semantics",
	},
	"issueimportcheckcsv": {
		status: "blocked_needs_design",
		rationale: "CSV import validation can expose imported row payloads and needs an " +
			"explicit redaction/output model",
	},
	"issueimportchecksync": {
		status: "blocked_needs_design",
		rationale: "sync import validation can expose external tracker payloads and needs an " +
			"explicit redaction/output model",
	},
	"issueimportcreateasana": {
		status: "blocked_needs_design",
		rationale: "Asana issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreatecsvjira": {
		status: "blocked_needs_design",
		rationale: "CSV/Jira issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreateclubhouse": {
		status: "blocked_needs_design",
		rationale: "Clubhouse issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreategithub": {
		status: "blocked_needs_design",
		rationale: "GitHub issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreatejira": {
		status: "blocked_needs_design",
		rationale: "Jira issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportjqlcheck": {
		status: "blocked_needs_design",
		rationale: "JQL import validation can expose external tracker payloads and needs an " +
			"explicit redaction/output model",
	},
	"issueimportprocess": {
		status: "blocked_needs_design",
		rationale: "issue import processing advances external import workflow state and needs " +
			"explicit integration guard semantics",
	},
	"issuelabelrestore": {
		status:    "blocked_needs_design",
		rationale: "issue label lifecycle restore needs explicit organization/admin safety semantics",
	},
	"issuelabelretire": {
		status:    "blocked_needs_design",
		rationale: "issue label lifecycle retire needs explicit organization/admin safety semantics",
	},
	"issuereminder": {
		status: "blocked_needs_design",
		rationale: "issue reminder mutation changes notification state and needs target-pinned " +
			"guard semantics",
	},
	"issuerepositorysuggestions": {
		status: "intentionally_excluded",
		rationale: "repository suggestion reads expose VCS integration metadata outside the " +
			"default Linear work CLI surface",
	},
	"issuedescriptionupdatefromfront": {
		status: "blocked_needs_design",
		rationale: "Front-origin description updates mutate issue content through integration state; " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreatelinearv2": {
		status: "blocked_needs_design",
		rationale: "Linear v2 issue import creation starts import workflow state and needs explicit " +
			"import guard semantics",
	},
	"issueshare": {
		status:    "blocked_needs_design",
		rationale: "issue sharing changes access state and needs target-pinned guard semantics",
	},
	"issuesubscribe": {
		status: "blocked_needs_design",
		rationale: "issue subscription changes notification state and needs target-pinned " +
			"guard semantics",
	},
	"issueunshare": {
		status:    "blocked_needs_design",
		rationale: "issue unsharing changes access state and needs target-pinned guard semantics",
	},
	"issueunsubscribe": {
		status: "blocked_needs_design",
		rationale: "issue unsubscribe changes notification state and needs target-pinned " +
			"guard semantics",
	},
	"latestreleasebyaccesskey": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"latestreleasebyaccesskey_history": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"latestreleasebyaccesskey_links": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"initiativelabeladd": {
		status:    "blocked_needs_design",
		rationale: "initiative label mutation needs initiative target pinning and target-mismatch tests",
	},
	"initiativeaddlabel": {
		status:    "blocked_needs_design",
		rationale: "initiative label mutation needs initiative target pinning and target-mismatch tests",
	},
	"microsoftteamschannels": {
		status: "intentionally_excluded",
		rationale: "Microsoft Teams channel enumeration exposes chat integration metadata outside the " +
			"default Linear work CLI surface",
	},
	"organizationinvite": {
		status:    "intentionally_excluded",
		rationale: organizationInviteRationale(),
	},
	"organizationinvites": {
		status:    "intentionally_excluded",
		rationale: organizationInviteRationale(),
	},
	"organizationinvitedetails": {
		status:    "intentionally_excluded",
		rationale: organizationInviteRationale(),
	},
	"organizationdomainclaimrequest": {
		status: "intentionally_excluded",
		rationale: "organization domain claim requests expose org-admin domain-verification " +
			"metadata outside the ordinary agent CLI surface",
	},
	"organization_subscription": {
		status:    "intentionally_excluded",
		rationale: "organization subscription and billing state is outside the ordinary agent CLI surface",
	},
	"pushsubscriptiontest": {
		status: "intentionally_excluded",
		rationale: "push subscription diagnostics are notification-device integration plumbing " +
			"outside the CLI surface",
	},
	"projectlabelrestore": {
		status:    "blocked_needs_design",
		rationale: "project label lifecycle restore needs explicit organization/admin safety semantics",
	},
	"projectlabelretire": {
		status:    "blocked_needs_design",
		rationale: "project label lifecycle retire needs explicit organization/admin safety semantics",
	},
	"projectaddlabel": {
		status:    "blocked_needs_design",
		rationale: "project label mutation needs project target pinning and target-mismatch tests",
	},
	"projectexternalsyncdisable": {
		status: "blocked_needs_design",
		rationale: "project external-sync disable changes integration state and needs explicit " +
			"integration guard semantics",
	},
	"projectcreateslackchannel": {
		status: "blocked_needs_design",
		rationale: "project Slack channel creation mutates chat integration state and needs explicit " +
			"integration guard semantics",
	},
	"projectreassignstatus": {
		status: "blocked_needs_design",
		rationale: "project status reassignment mutates project workflow state and needs target-pinned " +
			"guard semantics",
	},
	"recentreleasesbyaccesskey": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"releasepipelinebyaccesskey": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"releasepipelinebyaccesskey_releases": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"releasepipelinebyaccesskey_stages": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"ssourlfromemail": {
		status:    "intentionally_excluded",
		rationale: "SSO discovery from email belongs to auth flow tooling, not the Linear work CLI",
	},
	"userchangerole": {
		status:    "intentionally_excluded",
		rationale: "user role changes are organization administration outside the ordinary agent CLI surface",
	},
	"userdiscordconnect": {
		status:    "intentionally_excluded",
		rationale: "Discord account connection belongs to user auth/integration setup, not work CLI reads",
	},
	"userexternaluserdisconnect": {
		status: "intentionally_excluded",
		rationale: "external-user disconnection is identity integration administration outside the " +
			"ordinary agent CLI surface",
	},
	"usersettingsflagsreset": {
		status: "intentionally_excluded",
		rationale: "user settings flag reset is internal preference administration outside the " +
			"ordinary agent CLI surface",
	},
	"userunlinkfromidentityprovider": {
		status:    "intentionally_excluded",
		rationale: "identity-provider unlinking is auth administration outside the ordinary agent CLI surface",
	},
	"verifygithubenterpriseserverinstallation": {
		status: "intentionally_excluded",
		rationale: "GitHub Enterprise installation verification is integration administration " +
			"outside the CLI surface",
	},
}

func explicitRiskClassification(lowerName string) (string, string, bool) {
	classification, ok := explicitRiskClassifications[lowerName]
	return classification.status, classification.rationale, ok
}

func accessKeyReleaseRationale() string {
	return "access-key release reads are unauthenticated sharing surfaces " +
		"outside the auth-scoped agent CLI"
}

func contentPayloadReadRationale() string {
	return "content, thread, and archive payload reads can expose body/blob data; " +
		"needs explicit opt-in projection before CLI exposure"
}

func organizationInviteRationale() string {
	return "organization invite reads can expose invitee and admin metadata " +
		"outside an agent-safe CLI surface"
}

func countImplemented(fields []rootField, implementedRoots map[string]bool) int {
	return countWhere(fields, func(field rootField) bool {
		return implementedRoots[rootKey(field.Kind, field.Name)]
	})
}

func countImplementedSDK(methods []sdkMethod, implementedRoots map[string]bool) int {
	return countWhere(methods, func(method sdkMethod) bool {
		return sdkImplemented(method.Name, implementedRoots)
	})
}

func sdkImplemented(name string, implementedRoots map[string]bool) bool {
	if implementedRoots[rootKey("query", name)] || implementedRoots[rootKey("mutation", name)] {
		return true
	}
	for _, candidate := range sdkMutationRootCandidates(name) {
		if implementedRoots[rootKey("mutation", candidate)] {
			return true
		}
	}
	return false
}

func sdkRootKind(name string, rootKinds map[string]string) (string, bool) {
	if kind, ok := rootKinds[strings.ToLower(name)]; ok {
		return kind, true
	}
	for _, candidate := range sdkMutationRootCandidates(name) {
		if kind, ok := rootKinds[strings.ToLower(candidate)]; ok {
			return kind, true
		}
	}

	return "", false
}

func rootKey(kind string, name string) string {
	return strings.ToLower(kind) + ":" + name
}

func sdkMutationRootCandidates(name string) []string {
	candidates := []string{}
	for _, prefix := range []string{"create", "update", "archive", "delete", "unarchive", "cancel"} {
		if strings.HasPrefix(name, prefix) && len(name) > len(prefix) {
			entity := lowerFirst(strings.TrimPrefix(name, prefix))
			candidates = append(candidates, entity+upperFirst(prefix))
		}
	}
	return candidates
}

func hasWritePrefix(lowerName string) bool {
	for _, prefix := range []string{
		"create",
		"update",
		"archive",
		"delete",
		"unarchive",
		"cancel",
		"mark",
		"move",
		"rotate",
	} {
		if strings.HasPrefix(lowerName, prefix) || strings.HasSuffix(lowerName, prefix) {
			return true
		}
	}
	return false
}
