package client

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func ifUpdatedAtIssuePayloads(updatedAt string, teamID string) map[string]string {
	body := issueJSON(b1IssueFixture("LIT-1"))
	if updatedAt != "" {
		body = strings.Replace(body, `"assignee":null,`, `"assignee":null,"updatedAt":"`+updatedAt+`",`, 1)
	}
	if teamID != "" {
		body = strings.Replace(body, `"team":{"id":"team-id","key":"LIT"`, `"team":{"id":"`+teamID+`","key":"OTHER"`, 1)
	}

	return map[string]string{
		"issue":       `{"issue":` + body + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + body + `}}`,
	}
}

func Test_UpdateIssue_refuses_when_the_issue_changed_since_if_updated_at(t *testing.T) {
	recorder := &mutationRecordingClient{
		inner: issueWriteFakeClient(ifUpdatedAtIssuePayloads("2026-09-16T14:30:00.000Z", "")),
	}

	_, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:          "LIT-1",
		Title:       "new title",
		IfUpdatedAt: "2026-09-16T14:26:12Z",
	})

	require.ErrorIs(t, err, ErrWriteConflict)
	require.Contains(t, err.Error(), "2026-09-16T14:30:00.000Z")
	require.Contains(t, err.Error(), "2026-09-16T14:26:12Z")
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_UpdateIssue_proceeds_when_if_updated_at_names_the_same_instant(t *testing.T) {
	for name, expected := range map[string]string{
		"fewer fractional digits": "2026-09-16T14:26:12Z",
		"explicit zero fraction":  "2026-09-16T14:26:12.000Z",
		"other utc offset":        "2026-09-16T16:26:12+02:00",
	} {
		t.Run(name, func(t *testing.T) {
			recorder := &mutationRecordingClient{
				inner: issueWriteFakeClient(ifUpdatedAtIssuePayloads("2026-09-16T14:26:12.000Z", "")),
			}

			_, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
				ID:          "LIT-1",
				Title:       "new title",
				IfUpdatedAt: expected,
			})

			require.NoError(t, err)
			require.Equal(t, 1, countOperation(recorder.sent, "IssueUpdate"))
		})
	}
}

func Test_UpdateIssue_rejects_a_malformed_if_updated_at_before_any_request(t *testing.T) {
	for _, malformed := range []string{"yesterday", "2026-09-16", "2026-09-16 14:26:12"} {
		t.Run(malformed, func(t *testing.T) {
			recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

			_, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
				ID:          "LIT-1",
				Title:       "new title",
				IfUpdatedAt: malformed,
			})

			require.ErrorIs(t, err, ErrWriteInvalid)
			require.Empty(t, recorder.sent)
		})
	}
}

func Test_UpdateIssue_reports_target_mismatch_before_a_stale_if_updated_at(t *testing.T) {
	recorder := &mutationRecordingClient{
		inner: issueWriteFakeClient(ifUpdatedAtIssuePayloads("2026-09-16T14:30:00.000Z", "other-team-id")),
	}

	_, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:          "LIT-1",
		Title:       "new title",
		IfUpdatedAt: "2026-09-16T14:26:12Z",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.NotErrorIs(t, err, ErrWriteConflict)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_UpdateIssue_refuses_when_the_issue_has_no_readable_updated_at(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(ifUpdatedAtIssuePayloads("", ""))}

	_, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:          "LIT-1",
		Title:       "new title",
		IfUpdatedAt: "2026-09-16T14:26:12Z",
	})

	require.ErrorIs(t, err, ErrWriteConflict)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func countOperation(sent []string, name string) int {
	count := 0
	for _, operation := range sent {
		if operation == name {
			count++
		}
	}

	return count
}
