package main

import (
	"net/http"

	"github.com/codegp/cloud-persister/models"

	"github.com/codegp/game-object-types/types"
)

func GetInitialData(w http.ResponseWriter, r *http.Request) *requestError {
	user, err := getUserFromContext(r)
	if err != nil {
		return requestErrorf(err, "Error getting user from session")
	}

	projects, err := cp.GetMultiProject(user.ProjectIDs)
	if err != nil {
		return requestErrorf(err, "Error getting projects for user")
	}

	gameTypes, err := cp.ListGameTypes()
	if err != nil {
		return requestErrorf(err, "Error querying gameTypes from datastore")
	}

	gameTypeCache := map[int64]*models.GameType{}
	var maps []*models.Map
	for _, proj := range projects {
		gameType, found := gameTypeCache[proj.GameTypeID]
		if !found {
			gameType, err = cp.GetGameType(proj.GameTypeID)
			if err != nil {
				return requestErrorf(err, "Error getting gametype for project")
			}

			gameTypeCache[proj.GameTypeID] = gameType
		}

		pmaps, err := cp.GetMultiMap(gameType.MapIDs)
		if err != nil {
			return requestErrorf(err, "Error getting maps for project")
		}
		maps = append(maps, pmaps...)
	}

	botTypes, err := cp.ListBotTypes()
	if err != nil {
		return requestErrorf(err, "Error getting botTypes from datastore")
	}

	attackTypes, err := cp.ListAttackTypes()
	if err != nil {
		return requestErrorf(err, "Error getting attackTypes from datastore")
	}

	moveTypes, err := cp.ListMoveTypes()
	if err != nil {
		return requestErrorf(err, "Error getting moveTypes from datastore")
	}

	itemTypes, err := cp.ListItemTypes()
	if err != nil {
		return requestErrorf(err, "Error getting itemTypes from datastore")
	}

	terrainTypes, err := cp.ListTerrainTypes()
	if err != nil {
		return requestErrorf(err, "Error getting terrainTypes from datastore")
	}

	return marshalAndWriteResponse(w, &InitialData{
		Maps:         maps,
		Projects:     projects,
		GameTypes:    gameTypes,
		BotTypes:     botTypes,
		AttackTypes:  attackTypes,
		ItemTypes:    itemTypes,
		TerrainTypes: terrainTypes,
		MoveTypes:    moveTypes,
	})
}

type InitialData struct {
	Maps         []*models.Map
	Projects     []*models.Project
	GameTypes    []*models.GameType
	BotTypes     []*types.BotType
	AttackTypes  []*types.AttackType
	ItemTypes    []*types.ItemType
	TerrainTypes []*types.TerrainType
	MoveTypes    []*types.MoveType
}
