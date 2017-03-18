package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/codegp/cloud-persister/models"
)

func GetGame(w http.ResponseWriter, r *http.Request) *requestError {
	gameID, rerr := readIDFromRequest(r, "gameID")
	if rerr != nil {
		return rerr
	}

	game, err := cp.GetGame(gameID)
	if err != nil {
		return requestErrorf(err, "Error getting game from datastore")
	}

	return marshalAndWriteResponse(w, game)
}

func PostGame(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling itemType")
	}

	// TODO: if gametype.numteams > 1 find competitors
	game := &models.Game{
		Created:  time.Now(),
		Complete: false,
	}

	err = json.Unmarshal(body, game)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling game")
	}

	game, err = cp.AddGame(game)
	if err != nil {
		return requestErrorf(err, "Error adding game to datastore")
	}

	projects, err := cp.GetMultiProject(game.ProjectIDs)
	if err != nil {
		return requestErrorf(err, "Error getting projects from datastore")
	}

	_, err = kc.StartGame(game, projects)
	if err != nil {
		return requestErrorf(err, "Error starting game pod")
	}

	for _, proj := range projects {
		proj.GameIDs = append(proj.GameIDs, game.ID)

		// TODO: put multi
		if cp.UpdateProject(proj) != nil {
			return requestErrorf(err, "Error updating project")
		}
	}

	return marshalAndWriteResponse(w, game)
}

func UpdateGame(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling game")
	}

	var game *models.Game
	err = json.Unmarshal(body, game)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling game")
	}

	err = cp.UpdateGame(game)
	if err != nil {
		return requestErrorf(err, "Error updating game in datastore")
	}

	return marshalAndWriteResponse(w, game)
}

func GetGames(w http.ResponseWriter, r *http.Request) *requestError {
	projectID, rerr := readIDFromRequest(r, "projectID")
	if rerr != nil {
		return rerr
	}

	project, err := cp.GetProject(projectID)
	if err != nil {
		return requestErrorf(err, "Error getting project from datastore")
	}

	games := make([]*models.Game, len(project.GameIDs))
	for i, gameID := range project.GameIDs {
		if games[i], err = cp.GetGame(gameID); err != nil {
			return requestErrorf(err, "Error getting game from datastore")
		}
	}

	return marshalAndWriteResponse(w, games)
}

func GetHistory(w http.ResponseWriter, r *http.Request) *requestError {
	gameID, rerr := readIDFromRequest(r, "gameID")
	if rerr != nil {
		return rerr
	}

	json, err := cp.ReadHistory(gameID)
	if err != nil {
		return requestErrorf(err, "Error reading history from storage")
	}

	if _, err := w.Write(json); err != nil {
		return requestErrorf(err, "Error writing history")
	}
	return nil
}
