package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codegp/game-object-types/types"
)

func GetMoveType(w http.ResponseWriter, r *http.Request) *requestError {
	moveTypeID, rerr := readIDFromRequest(r, "moveTypeID")
	if rerr != nil {
		return rerr
	}

	moveType, err := cp.GetMoveType(moveTypeID)
	if err != nil {
		return requestErrorf(err, "Error getting moveType from datastore")
	}

	return marshalAndWriteResponse(w, moveType)
}

func PostMoveType(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling moveType")
	}

	var moveType *types.MoveType
	err = json.Unmarshal(body, &moveType)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling moveType")
	}

	moveType, err = cp.AddMoveType(moveType)
	if err != nil {
		return requestErrorf(err, "Error updating moveType in datastore")
	}

	return marshalAndWriteResponse(w, moveType)
}
