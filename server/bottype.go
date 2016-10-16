package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codegp/game-object-types/types"
)

func GetBotType(w http.ResponseWriter, r *http.Request) *requestError {
	botTypeID, rerr := readIDFromRequest(r, "botTypeID")
	if rerr != nil {
		return rerr
	}

	botType, err := cp.GetBotType(botTypeID)
	if err != nil {
		return requestErrorf(err, "Error getting botType from datastore")
	}

	return marshalAndWriteResponse(w, botType)
}

func PostBotType(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling botType")
	}

	var botType *types.BotType
	err = json.Unmarshal(body, &botType)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling botType")
	}

	botType, err = cp.AddBotType(botType)
	if err != nil {
		return requestErrorf(err, "Error updating botType in datastore")
	}

	return marshalAndWriteResponse(w, botType)
}
