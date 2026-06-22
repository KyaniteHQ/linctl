package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_normalizedPriorityValue_maps_all_aliases(t *testing.T) {
	tests := []struct {
		raw      string
		expected string
		changed  bool
	}{
		{raw: "urgent", expected: "1", changed: true},
		{raw: "high", expected: "2", changed: true},
		{raw: "medium", expected: "3", changed: true},
		{raw: "med", expected: "3", changed: true},
		{raw: "low", expected: "4", changed: true},
		{raw: "none", expected: "0", changed: true},
		{raw: "no priority", expected: "0", changed: true},
		// canonical pass-throughs
		{raw: "0", expected: "0", changed: false},
		{raw: "1", expected: "1", changed: false},
		{raw: "2", expected: "2", changed: false},
		{raw: "3", expected: "3", changed: false},
		{raw: "4", expected: "4", changed: false},
		// case-folding
		{raw: "URGENT", expected: "1", changed: true},
		{raw: "  High  ", expected: "2", changed: true},
	}
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			value, changed, err := normalizedPriorityValue(tt.raw)
			require.NoError(t, err)
			require.Equal(t, tt.expected, value)
			require.Equal(t, tt.changed, changed)
		})
	}
}

func Test_normalizedPriorityValue_returns_error_for_unknown(t *testing.T) {
	_, _, err := normalizedPriorityValue("critical")

	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown priority")
}

func Test_normalizedStateType_maps_all_aliases(t *testing.T) {
	tests := []struct {
		raw      string
		expected string
		changed  bool
	}{
		// canonical pass-throughs
		{raw: "triage", expected: "triage", changed: false},
		{raw: "backlog", expected: "backlog", changed: false},
		{raw: "unstarted", expected: "unstarted", changed: false},
		{raw: "started", expected: "started", changed: false},
		{raw: "completed", expected: "completed", changed: false},
		{raw: "cancelled", expected: "cancelled", changed: false},
		// aliases
		{raw: "todo", expected: "unstarted", changed: true},
		{raw: "to do", expected: "unstarted", changed: true},
		{raw: "in progress", expected: "started", changed: true},
		{raw: "in-progress", expected: "started", changed: true},
		{raw: "done", expected: "completed", changed: true},
		{raw: "complete", expected: "completed", changed: true},
		{raw: "closed", expected: "completed", changed: true},
		{raw: "canceled", expected: "cancelled", changed: true},
		{raw: "wont do", expected: "cancelled", changed: true},
		{raw: "wont-do", expected: "cancelled", changed: true},
		{raw: "won't do", expected: "cancelled", changed: true},
		// case-folding
		{raw: "TODO", expected: "unstarted", changed: true},
		{raw: "  Started  ", expected: "started", changed: true},
	}
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			value, changed, err := normalizedStateType(tt.raw)
			require.NoError(t, err)
			require.Equal(t, tt.expected, value)
			require.Equal(t, tt.changed, changed)
		})
	}
}

func Test_normalizedStateType_returns_error_for_unknown(t *testing.T) {
	_, _, err := normalizedStateType("sprinting")

	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown state type")
}

func Test_mergedStateFlag_returns_state_when_non_empty(t *testing.T) {
	result := mergedStateFlag("started", "todo")

	require.Equal(t, "started", result)
}

func Test_mergedStateFlag_returns_status_when_state_empty(t *testing.T) {
	result := mergedStateFlag("", "todo")

	require.Equal(t, "todo", result)
}

func Test_mergedStateFlag_returns_empty_when_both_empty(t *testing.T) {
	result := mergedStateFlag("", "")

	require.Empty(t, result)
}

func Test_writeNote_writes_to_stderr(t *testing.T) {
	var buf bytes.Buffer
	command := &cobra.Command{}
	command.SetErr(&buf)

	err := writeNote(command, "state %q normalized to %q", "todo", "unstarted")

	require.NoError(t, err)
	require.Equal(t, "note: state \"todo\" normalized to \"unstarted\"\n", buf.String())
}

func Test_writeNote_returns_error_on_write_failure(t *testing.T) {
	command := &cobra.Command{}
	command.SetErr(commandFailingWriter{})

	err := writeNote(command, "hello")

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func Test_normalizeAndNote_returns_empty_for_empty_raw(t *testing.T) {
	var buf bytes.Buffer
	command := &cobra.Command{}
	command.SetErr(&buf)

	value, err := normalizeAndNote(command, "state", "", normalizedStateType)

	require.NoError(t, err)
	require.Empty(t, value)
	require.Empty(t, buf.String())
}

func Test_normalizeAndNote_emits_note_when_value_changes(t *testing.T) {
	var buf bytes.Buffer
	command := &cobra.Command{}
	command.SetErr(&buf)

	value, err := normalizeAndNote(command, "state", "todo", normalizedStateType)

	require.NoError(t, err)
	require.Equal(t, "unstarted", value)
	require.Equal(t, "note: state \"todo\" normalized to \"unstarted\"\n", buf.String())
}

func Test_normalizeAndNote_is_silent_for_canonical_value(t *testing.T) {
	var buf bytes.Buffer
	command := &cobra.Command{}
	command.SetErr(&buf)

	value, err := normalizeAndNote(command, "state", "started", normalizedStateType)

	require.NoError(t, err)
	require.Equal(t, "started", value)
	require.Empty(t, buf.String())
}

func Test_normalizeAndNote_returns_normalize_error(t *testing.T) {
	command := &cobra.Command{}
	command.SetErr(&bytes.Buffer{})

	_, err := normalizeAndNote(command, "state", "sprinting", normalizedStateType)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown state type")
}

func Test_normalizeAndNote_returns_note_write_error(t *testing.T) {
	command := &cobra.Command{}
	command.SetErr(commandFailingWriter{})

	_, err := normalizeAndNote(command, "priority", "urgent", normalizedPriorityValue)

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}
