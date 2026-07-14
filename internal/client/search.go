package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// SearchDocumentSummary is a compact document search result.
type SearchDocumentSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	SlugID     string `json:"slug_id"`
	URL        string `json:"url"`
	ParentType string `json:"parent_type,omitempty"`
	ParentID   string `json:"parent_id,omitempty"`
	ParentName string `json:"parent_name,omitempty"`
}

// SearchDocumentList is a page of document search results.
type SearchDocumentList struct {
	Documents   []SearchDocumentSummary `json:"documents"`
	TotalCount  float64                 `json:"total_count"`
	HasNextPage bool                    `json:"has_next_page"`
	EndCursor   *string                 `json:"end_cursor,omitempty"`
}

// SearchIssueSummary is a compact issue search result.
type SearchIssueSummary struct {
	ID          string `json:"id"`
	Identifier  string `json:"identifier"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	TeamID      string `json:"team_id"`
	TeamKey     string `json:"team_key"`
	TeamName    string `json:"team_name"`
	StateID     string `json:"state_id"`
	StateName   string `json:"state_name"`
	StateType   string `json:"state_type"`
	ProjectID   string `json:"project_id,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
}

// SearchIssueList is a page of issue search results.
type SearchIssueList struct {
	Issues      []SearchIssueSummary `json:"issues"`
	TotalCount  float64              `json:"total_count"`
	HasNextPage bool                 `json:"has_next_page"`
	EndCursor   *string              `json:"end_cursor,omitempty"`
}

// SearchProjectSummary is a compact project search result.
type SearchProjectSummary struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	SlugID string        `json:"slug_id"`
	URL    string        `json:"url"`
	Status ProjectStatus `json:"status"`
	Lead   string        `json:"lead,omitempty"`
	Teams  []ProjectTeam `json:"teams"`
}

// SearchProjectList is a page of project search results.
type SearchProjectList struct {
	Projects    []SearchProjectSummary `json:"projects"`
	TotalCount  float64                `json:"total_count"`
	HasNextPage bool                   `json:"has_next_page"`
	EndCursor   *string                `json:"end_cursor,omitempty"`
}

// SearchDocuments returns compact document search results.
func SearchDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	term string,
	limit int,
) (SearchDocumentList, error) {
	result, err := gql.XSearchDocuments(ctx, graphqlClient, term, intPtr(limit), nil)
	if err != nil {
		return SearchDocumentList{}, fmt.Errorf("search documents: %w", err)
	}

	documents := mapNodes(result.SearchDocuments.Nodes, func(
		node gql.XSearchDocumentsSearchDocumentsDocumentSearchPayloadNodesDocumentSearchResult,
	) SearchDocumentSummary {
		return searchDocumentSummary(node.SearchDocumentSummaryFields)
	})

	return SearchDocumentList{
		Documents:   documents,
		TotalCount:  result.SearchDocuments.TotalCount,
		HasNextPage: result.SearchDocuments.PageInfo.HasNextPage,
		EndCursor:   result.SearchDocuments.PageInfo.EndCursor,
	}, nil
}

// SearchIssues returns compact issue search results.
func SearchIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	term string,
	limit int,
) (SearchIssueList, error) {
	result, err := gql.XSearchIssues(ctx, graphqlClient, term, intPtr(limit), nil)
	if err != nil {
		return SearchIssueList{}, fmt.Errorf("search issues: %w", err)
	}

	issues := mapNodes(result.SearchIssues.Nodes, func(
		node gql.XSearchIssuesSearchIssuesIssueSearchPayloadNodesIssueSearchResult,
	) SearchIssueSummary {
		return typedSearchIssueSummary(node.SearchIssueSummaryFields)
	})

	return SearchIssueList{
		Issues:      issues,
		TotalCount:  result.SearchIssues.TotalCount,
		HasNextPage: result.SearchIssues.PageInfo.HasNextPage,
		EndCursor:   result.SearchIssues.PageInfo.EndCursor,
	}, nil
}

// SearchProjects returns compact project search results.
func SearchProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	term string,
	limit int,
) (SearchProjectList, error) {
	result, err := gql.XSearchProjects(ctx, graphqlClient, term, intPtr(limit), nil)
	if err != nil {
		return SearchProjectList{}, fmt.Errorf("search projects: %w", err)
	}

	projects := mapNodes(result.SearchProjects.Nodes, func(
		node gql.XSearchProjectsSearchProjectsProjectSearchPayloadNodesProjectSearchResult,
	) SearchProjectSummary {
		return searchProjectSummary(node.SearchProjectSummaryFields)
	})

	return SearchProjectList{
		Projects:    projects,
		TotalCount:  result.SearchProjects.TotalCount,
		HasNextPage: result.SearchProjects.PageInfo.HasNextPage,
		EndCursor:   result.SearchProjects.PageInfo.EndCursor,
	}, nil
}

func searchDocumentSummary(fields gql.SearchDocumentSummaryFields) SearchDocumentSummary {
	summary := SearchDocumentSummary{
		ID:     fields.Id,
		Title:  fields.Title,
		SlugID: fields.SlugId,
		URL:    fields.Url,
	}
	if fields.Project != nil {
		summary.ParentType = "project"
		summary.ParentID = fields.Project.Id
		summary.ParentName = fields.Project.Name
	}
	if fields.Initiative != nil {
		summary.ParentType = "initiative"
		summary.ParentID = fields.Initiative.Id
		summary.ParentName = fields.Initiative.Name
	}
	if fields.Team != nil {
		summary.ParentType = "team"
		summary.ParentID = fields.Team.Id
		summary.ParentName = fields.Team.Name
	}
	if fields.Issue != nil {
		summary.ParentType = "issue"
		summary.ParentID = fields.Issue.Id
		summary.ParentName = fields.Issue.Identifier
	}
	if fields.Release != nil {
		summary.ParentType = "release"
		summary.ParentID = fields.Release.Id
		summary.ParentName = fields.Release.Name
	}
	if fields.Cycle != nil {
		summary.ParentType = "cycle"
		summary.ParentID = fields.Cycle.Id
		summary.ParentName = fmt.Sprintf("Cycle %.0f", fields.Cycle.Number)
		if fields.Cycle.Name != nil && *fields.Cycle.Name != "" {
			summary.ParentName = *fields.Cycle.Name
		}
	}

	return summary
}

func typedSearchIssueSummary(fields gql.SearchIssueSummaryFields) SearchIssueSummary {
	summary := SearchIssueSummary{
		ID:         fields.Id,
		Identifier: fields.Identifier,
		Title:      fields.Title,
		URL:        fields.Url,
		TeamID:     fields.Team.Id,
		TeamKey:    fields.Team.Key,
		TeamName:   fields.Team.Name,
		StateID:    fields.State.Id,
		StateName:  fields.State.Name,
		StateType:  fields.State.Type,
	}
	if fields.Project != nil {
		summary.ProjectID = fields.Project.Id
		summary.ProjectName = fields.Project.Name
	}

	return summary
}

func searchProjectSummary(fields gql.SearchProjectSummaryFields) SearchProjectSummary {
	teams := mapNodes(fields.Teams.Nodes, func(
		team gql.SearchProjectSummaryFieldsTeamsTeamConnectionNodesTeam,
	) ProjectTeam {
		return ProjectTeam{
			ID:   team.Id,
			Key:  team.Key,
			Name: team.Name,
		}
	})

	summary := SearchProjectSummary{
		ID:     fields.Id,
		Name:   fields.Name,
		SlugID: fields.SlugId,
		URL:    fields.Url,
		Status: ProjectStatus{
			ID:   fields.Status.Id,
			Name: fields.Status.Name,
			Type: string(fields.Status.Type),
		},
		Teams: teams,
	}
	if fields.Lead != nil {
		summary.Lead = fields.Lead.DisplayName
	}

	return summary
}
