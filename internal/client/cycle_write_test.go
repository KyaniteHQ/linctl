package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_CreateCycle_returns_created_cycle_when_target_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"CycleCreate": `{"cycleCreate":{"success":true,"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}}`,
	})

	cycle, err := CreateCycle(context.Background(), graphqlClient, matchingTarget(), CycleCreateRequest{
		Name:        "Planning cycle",
		Description: "cycle body",
		StartsAt:    "2026-07-01T00:00:00Z",
		EndsAt:      "2026-07-15T00:00:00Z",
	})

	require.NoError(t, err)
	require.Equal(t, "cycle-id", cycle.ID)
	require.Equal(t, "Planning cycle", cycle.Name)
	require.Equal(t, "team-id", cycle.TeamID)
	require.Equal(t, "LIT", cycle.TeamKey)
}

func Test_CreateCycle_returns_mutation_failed_when_payload_omits_cycle(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"CycleCreate": `{"cycleCreate":{"success":true,"cycle":null}}`,
	})

	_, err := CreateCycle(context.Background(), graphqlClient, matchingTarget(), CycleCreateRequest{
		Name:     "Planning cycle",
		StartsAt: "2026-07-01T00:00:00Z",
		EndsAt:   "2026-07-15T00:00:00Z",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_UpdateCycle_returns_updated_cycle_when_target_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"cycle":       `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
		"CycleUpdate": `{"cycleUpdate":{"success":true,"cycle":` + cycleJSON("Updated cycle", "team-id", "LIT") + `}}`,
	})

	cycle, err := UpdateCycle(context.Background(), graphqlClient, matchingTarget(), CycleUpdateRequest{
		ID:   "cycle-id",
		Name: "Updated cycle",
	})

	require.NoError(t, err)
	require.Equal(t, "cycle-id", cycle.ID)
	require.Equal(t, "Updated cycle", cycle.Name)
}

func Test_UpdateCycle_returns_mutation_failed_when_payload_omits_cycle(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"cycle":       `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
		"CycleUpdate": `{"cycleUpdate":{"success":true,"cycle":null}}`,
	})

	_, err := UpdateCycle(context.Background(), graphqlClient, matchingTarget(), CycleUpdateRequest{
		ID:   "cycle-id",
		Name: "Updated cycle",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_ArchiveCycle_returns_archived_cycle_when_target_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"cycle":        `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
		"CycleArchive": `{"cycleArchive":{"success":true,"entity":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}}`,
	})

	cycle, err := ArchiveCycle(context.Background(), graphqlClient, matchingTarget(), "cycle-id")

	require.NoError(t, err)
	require.Equal(t, "cycle-id", cycle.ID)
	require.Equal(t, "Planning cycle", cycle.Name)
}

func Test_ArchiveCycle_returns_mutation_failed_when_payload_omits_entity(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"cycle":        `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
		"CycleArchive": `{"cycleArchive":{"success":true,"entity":null}}`,
	})

	_, err := ArchiveCycle(context.Background(), graphqlClient, matchingTarget(), "cycle-id")

	require.Error(t, err)
	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_UpdateCycle_refuses_when_team_differs(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"cycle": `{"cycle":` + cycleJSON("Wrong team cycle", "other-team", "OTHER") + `}`,
	})

	_, err := UpdateCycle(
		context.Background(),
		graphqlClient,
		config.Target{
			OrgID:   "org-id",
			TeamKey: "LIT",
			TeamID:  "team-id",
		},
		CycleUpdateRequest{
			ID:   "cycle-id",
			Name: "Updated cycle",
		},
	)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func cycleJSON(name string, teamID string, teamKey string) string {
	return `{
		"id":"cycle-id",
		"number":12,
		"name":"` + name + `",
		"description":"cycle body",
		"startsAt":"2026-07-01T00:00:00Z",
		"endsAt":"2099-01-01T00:00:00Z",
		"completedAt":null,
		"progress":0.25,
		"team":{"id":"` + teamID + `","key":"` + teamKey + `","name":"linctl"}
	}`
}
