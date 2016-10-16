package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codegp/cloud-persister/models"
)

func GetGameType(w http.ResponseWriter, r *http.Request) *requestError {
	gametypeID, rerr := readIDFromRequest(r, "gameTypeID")
	if rerr != nil {
		return rerr
	}

	gameType, err := cp.GetGameType(gametypeID)
	if err != nil {
		return requestErrorf(err, "Error getting gameType from datastore")
	}

	return marshalAndWriteResponse(w, gameType)
}

func PostGameType(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling gameType")
	}

	var gameType *models.GameType
	err = json.Unmarshal(body, &gameType)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling gameType")
	}

	user, err := getUserFromContext(r)
	if err != nil {
		return requestErrorf(err, "Error getting user to set creatorID")
	}

	gameType.CreatorID = user.ID
	gameType, err = cp.AddGameType(gameType)
	if err != nil {
		return requestErrorf(err, "Error adding gameType to datastore")
	}

	return marshalAndWriteResponse(w, gameType)
}

func PostGameTypeCode(w http.ResponseWriter, r *http.Request) *requestError {
	gameTypeID, rerr := readIDFromRequest(r, "gameTypeID")
	if rerr != nil {
		return rerr
	}

	f, _, err := r.FormFile("code")
	if err != nil {
		return requestErrorf(err, "Error getting code from form body")
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return requestErrorf(err, "Error reading file into byte buffer")
	}

	err = cp.WriteGameTypeCode(gameTypeID, buf.Bytes())
	if err != nil {
		return requestErrorf(err, "Error writing file to storage")
	}

	return nil
}

func GetGameTypes(w http.ResponseWriter, r *http.Request) *requestError {
	gameTypes, err := cp.ListGameTypes()
	if err != nil {
		return requestErrorf(err, "Error querying gameTypes from datastore")
	}

	return marshalAndWriteResponse(w, gameTypes)
}
