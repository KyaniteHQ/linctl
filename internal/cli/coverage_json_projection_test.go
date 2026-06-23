package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CliOutputHelpers_cover_json_projection_and_sort_edges(t *testing.T) {
	projected, err := projectJSONFields(
		map[string]any{"issues": []any{map[string]any{"identifier": "LIT-1", "state": map[string]any{"name": "Todo"}}}},
		"identifier,state.name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"issues": []any{map[string]any{"identifier": "LIT-1", "state": map[string]any{"name": "Todo"}}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"projects": []any{map[string]any{"id": "project-id", "status": map[string]any{"name": "Backlog"}}}},
		"id,status.name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"projects": []any{map[string]any{"id": "project-id", "status": map[string]any{"name": "Backlog"}}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"members": []any{map[string]any{"id": "user-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"members": []any{map[string]any{"id": "user-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"customers": []any{map[string]any{"id": "customer-id", "status_name": "Active"}}},
		"id,status_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"customers": []any{map[string]any{"id": "customer-id", "status_name": "Active"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"roadmaps": []any{map[string]any{"id": "roadmap-id", "slug_id": "platform-roadmap"}}},
		"id,slug_id",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"roadmaps": []any{map[string]any{"id": "roadmap-id", "slug_id": "platform-roadmap"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"time_schedules": []any{map[string]any{"id": "time-schedule-id", "entry_count": float64(1)}}},
		"id,entry_count",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"time_schedules": []any{map[string]any{"id": "time-schedule-id", "entry_count": float64(1)}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"notifications": []any{map[string]any{"id": "notification-id", "category": "mentions"}}},
		"id,category",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"notifications": []any{map[string]any{"id": "notification-id", "category": "mentions"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{
			"notification_subscriptions": []any{
				map[string]any{"id": "notification-subscription-id", "target_type": "project"},
			},
		},
		"id,target_type",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"notification_subscriptions": []any{
			map[string]any{"id": "notification-subscription-id", "target_type": "project"},
		},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"triage_responsibilities": []any{
			map[string]any{"id": "triage-responsibility-id", "team_key": "LIT", "action": "notify"},
		}},
		"id,team_key,action",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"triage_responsibilities": []any{
			map[string]any{"id": "triage-responsibility-id", "team_key": "LIT", "action": "notify"},
		},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"sla_configurations": []any{
			map[string]any{"id": "sla-configuration-id", "name": "First response", "sla_type": "all"},
		}},
		"id,name,sla_type",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"sla_configurations": []any{
			map[string]any{"id": "sla-configuration-id", "name": "First response", "sla_type": "all"},
		},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"audit_entry_types": []any{map[string]any{"type": "user_login", "description": "User logged in"}}},
		"type,description",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"audit_entry_types": []any{map[string]any{"type": "user_login", "description": "User logged in"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"release_pipelines": []any{map[string]any{"id": "release-pipeline-id", "slug_id": "production"}}},
		"id,slug_id",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"release_pipelines": []any{map[string]any{"id": "release-pipeline-id", "slug_id": "production"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"release_stages": []any{map[string]any{"id": "release-stage-id", "pipeline_name": "Production"}}},
		"id,pipeline_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"release_stages": []any{map[string]any{"id": "release-stage-id", "pipeline_name": "Production"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"customer_needs": []any{map[string]any{"id": "customer-need-id", "customer_name": "Acme"}}},
		"id,customer_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"customer_needs": []any{map[string]any{"id": "customer-need-id", "customer_name": "Acme"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"customer_statuses": []any{map[string]any{"id": "customer-status-id", "display_name": "Active"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"customer_statuses": []any{map[string]any{"id": "customer-status-id", "display_name": "Active"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"customer_tiers": []any{map[string]any{"id": "customer-tier-id", "display_name": "Enterprise"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"customer_tiers": []any{map[string]any{"id": "customer-tier-id", "display_name": "Enterprise"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"comments": []any{map[string]any{"id": "comment-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"comments": []any{map[string]any{"id": "comment-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"documents": []any{map[string]any{"id": "document-id", "title": "Spec"}}},
		"id,title",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"documents": []any{map[string]any{"id": "document-id", "title": "Spec"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"labels": []any{map[string]any{"id": "label-id", "color": "#ff0000"}}},
		"id,color",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"labels": []any{map[string]any{"id": "label-id", "color": "#ff0000"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"teams": []any{map[string]any{"id": "team-id", "key": "LIT"}}},
		"id,key",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"teams": []any{map[string]any{"id": "team-id", "key": "LIT"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"users": []any{map[string]any{"id": "user-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"users": []any{map[string]any{"id": "user-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"drafts": []any{map[string]any{"id": "draft-id", "parent_key": "LIT-3"}}},
		"id,parent_key",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"drafts": []any{map[string]any{"id": "draft-id", "parent_key": "LIT-3"}},
	}, projected)

	projected, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "identifier")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	projected, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "identifier,, ")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	_, err = projectJSONFields(map[string]any{"bad": func() {}}, "bad")
	require.Error(t, err)
	require.Contains(t, err.Error(), "marshal output")

	_, err = projectJSONFields([]string{"not-an-object"}, "id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode output")

	_, err = projectJSONFields(map[string]any{"issues": []any{"bad-item"}}, "identifier")
	require.Error(t, err)
	require.Contains(t, err.Error(), "item is not an object")

	_, err = projectJSONFields(map[string]any{"issues": []any{map[string]any{"title": "Missing id"}}}, "identifier")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"identifier\" is not present")

	_, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "missing")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"missing\" is not present")

	_, err = projectJSONFields(map[string]any{"state": "Todo"}, "state.name")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"state\" is not an object")

	items := []client.IssueSummary{
		{Identifier: "LIT-2", Title: "Zebra"},
		{Identifier: "LIT-1", Title: "Alpha"},
	}
	sortedItems, err := sortByJSONField(items, "", "asc")
	require.NoError(t, err)
	require.Equal(t, items, sortedItems)

	sortedItems, err = sortByJSONField(items, "title", "asc")
	require.NoError(t, err)
	require.Equal(t, "Alpha", sortedItems[0].Title)

	_, err = sortByJSONField(items, "title", "sideways")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid sort order")

	_, err = sortByJSONField(items, "missing", "asc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "sort field \"missing\" is not present")

	_, err = sortByJSONField([]map[string]any{{"state": "Todo"}}, "state.name", "asc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not an object path")

	_, err = jsonFieldValue(map[string]any{"bad": func() {}}, "bad")
	require.Error(t, err)
	require.Contains(t, err.Error(), "marshal output")

	destination := map[string]any{}
	require.NoError(t, copyJSONPath(map[string]any{"id": "issue-id"}, destination, nil))
	require.Empty(t, destination)
}
