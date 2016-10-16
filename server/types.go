package main

import (
	"net/http"

	"github.com/codegp/game-object-types/types"
)

type TypeLists struct {
	BotTypes     []*types.BotType     `json:"botTypes"`
	AttackTypes  []*types.AttackType  `json:"attackTypes"`
	MoveTypes    []*types.MoveType    `json:"moveTypes"`
	ItemTypes    []*types.ItemType    `json:"itemTypes"`
	TerrainTypes []*types.TerrainType `json:"terrainTypes"`
}

func GetTypes(w http.ResponseWriter, r *http.Request) *requestError {
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

	typeLists := &TypeLists{
		BotTypes:     botTypes,
		AttackTypes:  attackTypes,
		MoveTypes:    moveTypes,
		ItemTypes:    itemTypes,
		TerrainTypes: terrainTypes,
	}

	return marshalAndWriteResponse(w, typeLists)
}
