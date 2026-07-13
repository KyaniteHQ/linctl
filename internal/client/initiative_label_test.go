package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_InitiativeLabelReads_return_compact_models(t *testing.T) {
	endCursor := "cursor-1"
	graphqlClient := fakeGraphQLClient{
		"initiativeLabels": `{"initiativeLabels":{"nodes":[{"id":"initiative-label-id","name":"Strategy","description":"Strategic theme","color":"#5e6ad2","isGroup":false,"lastAppliedAt":"2026-07-10T12:00:00Z","retiredAt":null,"archivedAt":null,"createdAt":"2026-07-01T12:00:00Z","updatedAt":"2026-07-10T12:00:00Z","parent":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"initiativeLabel":  `{"initiativeLabel":{"id":"initiative-label-id","name":"Strategy","description":"Strategic theme","color":"#5e6ad2","isGroup":false,"lastAppliedAt":"2026-07-10T12:00:00Z","retiredAt":null,"archivedAt":null,"createdAt":"2026-07-01T12:00:00Z","updatedAt":"2026-07-10T12:00:00Z","parent":{"id":"initiative-label-group-id","name":"Themes","color":"#8a8f98"}}}`,
	}

	labels, err := ListInitiativeLabels(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	label, err := GetInitiativeLabelByID(context.Background(), graphqlClient, "initiative-label-id")
	require.NoError(t, err)

	require.True(t, labels.HasNextPage)
	require.Equal(t, endCursor, *labels.EndCursor)
	require.Equal(t, "Strategy", labels.InitiativeLabels[0].Name)
	require.Empty(t, labels.InitiativeLabels[0].ParentID)
	require.Equal(t, "initiative-label-id", label.ID)
	require.Equal(t, "initiative-label-group-id", label.ParentID)
	require.Equal(t, "Themes", label.ParentName)
	require.Equal(t, "#8a8f98", label.ParentColor)
}

func Test_InitiativeLabelReads_wrap_errors(t *testing.T) {
	graphqlClient := errorGraphQLClient{err: errors.New("request failed")}

	_, err := ListInitiativeLabels(context.Background(), graphqlClient, 1)
	require.EqualError(t, err, "list initiative labels: request failed")
	_, err = GetInitiativeLabelByID(context.Background(), graphqlClient, "initiative-label-id")
	require.EqualError(t, err, "get initiative label initiative-label-id: request failed")
}
