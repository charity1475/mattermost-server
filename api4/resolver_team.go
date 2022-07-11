// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader/v6"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/web"
)

func getGraphQLTeam(ctx context.Context, id string) (*model.Team, error) {
	loader, err := getTeamsLoader(ctx)
	if err != nil {
		return nil, err
	}

	thunk := loader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	team := result.(*model.Team)
	return team, nil
}

func graphQLTeamsLoader(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	stringKeys := keys.Keys()
	result := make([]*dataloader.Result, len(stringKeys))

	c, err := getCtx(ctx)
	if err != nil {
		for i := range result {
			result[i] = &dataloader.Result{Error: err}
		}
		return result
	}

	teams, err := getGraphQLTeams(c, stringKeys)
	if err != nil {
		for i := range result {
			result[i] = &dataloader.Result{Error: err}
		}
		return result
	}

	for i, ch := range teams {
		result[i] = &dataloader.Result{Data: ch}
	}
	return result
}

func getGraphQLTeams(c *web.Context, teamIDs []string) ([]*model.Team, error) {
	teams, appErr := c.App.GetTeams(teamIDs)
	if appErr != nil {
		return nil, appErr
	}

	if len(teams) != len(teamIDs) {
		return nil, fmt.Errorf("All teams were not found. Requested %d; Found %d", len(teamIDs), len(teams))
	}

	var teamsToCheck []string
	for i := range teams {
		if !teams[i].AllowOpenInvite || teams[i].Type != model.TeamOpen {
			teamsToCheck = append(teamsToCheck, teams[i].Id)
		}
	}

	if !c.App.SessionHasPermissionToTeams(*c.AppContext.Session(), teamsToCheck, model.PermissionViewMembers) {
		c.SetPermissionError(model.PermissionViewTeam)
		return nil, c.Err
	}

	for i := range teams {
		teams[i] = c.App.SanitizeTeam(*c.AppContext.Session(), teams[i])
	}

	// The teams need to be in the exact same order as the input slice.
	tmp := make(map[string]*model.Team)
	for _, ch := range teams {
		tmp[ch.Id] = ch
	}

	// We reuse the same slice and just rewrite the teams.
	for i, id := range teamIDs {
		teams[i] = tmp[id]
	}

	return teams, nil
}
