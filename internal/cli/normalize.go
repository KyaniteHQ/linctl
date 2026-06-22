package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var priorityAliases = map[string]string{
	"urgent":      "1",
	"high":        "2",
	"medium":      "3",
	"med":         "3",
	"low":         "4",
	"none":        "0",
	"no priority": "0",
	"0":           "0",
	"1":           "1",
	"2":           "2",
	"3":           "3",
	"4":           "4",
}

var stateTypeAliases = map[string]string{
	"triage":      "triage",
	"backlog":     "backlog",
	"unstarted":   "unstarted",
	"started":     "started",
	"completed":   "completed",
	"cancelled":   "cancelled",
	"todo":        "unstarted",
	"to do":       "unstarted",
	"in progress": "started",
	"in-progress": "started",
	"done":        "completed",
	"complete":    "completed",
	"closed":      "completed",
	"canceled":    "cancelled",
	"wont do":     "cancelled",
	"wont-do":     "cancelled",
	"won't do":    "cancelled",
}

var healthAliases = map[string]string{
	"ontrack":   "onTrack",
	"on track":  "onTrack",
	"on-track":  "onTrack",
	"atrisk":    "atRisk",
	"at risk":   "atRisk",
	"at-risk":   "atRisk",
	"offtrack":  "offTrack",
	"off track": "offTrack",
	"off-track": "offTrack",
}

// normalizedHealthValue maps a raw project-update health string to its canonical
// Linear ProjectUpdateHealthType. It returns whether the value changed and any error.
func normalizedHealthValue(raw string) (value string, changed bool, err error) {
	key := strings.ToLower(strings.TrimSpace(raw))
	canonical, ok := healthAliases[key]
	if !ok {
		return "", false, fmt.Errorf("unknown health %q: use on-track, at-risk, or off-track", raw)
	}

	return canonical, canonical != raw, nil
}

// normalizedPriorityValue maps a raw priority string to a canonical numeric
// string value. It returns whether the value changed and any parse error.
func normalizedPriorityValue(raw string) (value string, changed bool, err error) {
	key := strings.ToLower(strings.TrimSpace(raw))
	canonical, ok := priorityAliases[key]
	if !ok {
		return "", false, fmt.Errorf("unknown priority %q: use urgent/high/medium/low/none or 0-4", raw)
	}

	return canonical, canonical != raw, nil
}

// normalizedStateType maps a raw state type string to a canonical Linear
// workflow state type. It returns whether the value changed and any parse error.
func normalizedStateType(raw string) (value string, changed bool, err error) {
	key := strings.ToLower(strings.TrimSpace(raw))
	canonical, ok := stateTypeAliases[key]
	if !ok {
		return "", false, fmt.Errorf(
			"unknown state type %q: use triage/backlog/unstarted/started/completed/cancelled",
			raw,
		)
	}

	return canonical, canonical != raw, nil
}

// mergedStateFlag returns state when non-empty, and status otherwise.
// It models the same pattern as issueListCreatedAfter: the canonical flag wins.
func mergedStateFlag(state string, status string) string {
	if state != "" {
		return state
	}

	return status
}

// normalizeAndNote normalizes raw with normalize, emitting a single stderr note
// (labelled by field) when normalization changed the value. An empty raw is
// returned unchanged without invoking normalize.
func normalizeAndNote(
	command *cobra.Command,
	field string,
	raw string,
	normalize func(string) (string, bool, error),
) (string, error) {
	if raw == "" {
		return "", nil
	}
	canonical, changed, err := normalize(raw)
	if err != nil {
		return "", err
	}
	if changed {
		if noteErr := writeNote(command, "%s %q normalized to %q", field, raw, canonical); noteErr != nil {
			return "", noteErr
		}
	}

	return canonical, nil
}

// writeNote writes a "note: " prefix line to the command's stderr stream.
func writeNote(command *cobra.Command, format string, args ...any) error {
	_, err := fmt.Fprintf(command.ErrOrStderr(), "note: "+format+"\n", args...)

	return err
}
