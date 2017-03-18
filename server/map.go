package main

import (
	"io/ioutil"
	"net/http"

	"github.com/codegp/cloud-persister/models"
	"github.com/gorilla/mux"
)

func PostMap(w http.ResponseWriter, r *http.Request) *requestError {
	vars := mux.Vars(r)
	mapName := vars["mapName"]
	gameTypeID, rerr := readIDFromRequest(r, "gameTypeID")
	if rerr != nil {
		return rerr
	}

	gameType, err := cp.GetGameType(gameTypeID)
	if err != nil {
		return requestErrorf(err, "Error getting gametype")
	}

	m := &models.Map{
		Name:       mapName,
		GameTypeID: gameTypeID,
	}

	m, err = cp.AddMap(m)
	if err != nil {
		return requestErrorf(err, "Error adding map model to datastore")
	}

	gameType.MapIDs = append(gameType.MapIDs, m.ID)
	err = cp.UpdateGameType(gameType)
	if err != nil {
		return requestErrorf(err, "Error adding map id to gametype")
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling gameType")
	}

	err = cp.WriteMap(m.ID, body)
	if err != nil {
		return requestErrorf(err, "Error writing file to storage")
	}

	return nil
}
