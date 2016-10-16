package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codegp/game-object-types/types"
)

func GetTerrainType(w http.ResponseWriter, r *http.Request) *requestError {
	terrainTypeID, rerr := readIDFromRequest(r, "terrainTypeID")
	if rerr != nil {
		return rerr
	}

	terrainType, err := cp.GetTerrainType(terrainTypeID)
	if err != nil {
		return requestErrorf(err, "Error getting terrainType from datastore")
	}

	return marshalAndWriteResponse(w, terrainType)
}

func PostTerrainType(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling terrainType")
	}

	var terrainType *types.TerrainType
	err = json.Unmarshal(body, &terrainType)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling terrainType")
	}

	terrainType, err = cp.AddTerrainType(terrainType)
	if err != nil {
		return requestErrorf(err, "Error updating terrainType in datastore")
	}

	return marshalAndWriteResponse(w, terrainType)
}
