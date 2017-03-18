package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/codegp/cloud-persister/models"
	"github.com/codegp/env"
	"github.com/gorilla/mux"
)

type requestHandler func(http.ResponseWriter, *http.Request) *requestError
type requestError struct {
	Error   error
	Message string
	Code    int
}

func requestErrorf(err error, format string, v ...interface{}) *requestError {
	return &requestError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    http.StatusInternalServerError,
	}
}

func sessionMiddleware(h requestHandler) func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
		if !env.IsLocal() && profileFromSession(r) == nil && githubUserFromSession(r) == nil {
			http.Error(w, "Please go to the splash page and login", http.StatusUnauthorized)
			return
		}

		executeHandler(h, w, r)
	})
}

func errorMiddleware(h requestHandler) func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
		executeHandler(h, w, r)
	})
}

func executeHandler(h requestHandler, w http.ResponseWriter, r *http.Request) {
	if e := h(w, r); e != nil {
		logger.Errorf("Handler error: status code: %d, message: %s, underlying err: %v", e.Code, e.Message, e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

func readIDFromRequest(r *http.Request, varName string) (int64, *requestError) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseInt(vars[varName], 10, 64)
	if err != nil {
		return ID, requestErrorf(err, "Invalid %v, must be integer", varName)
	}

	return ID, nil
}

func marshalAndWriteResponse(w http.ResponseWriter, toMarshal interface{}) *requestError {
	content, err := json.Marshal(toMarshal)
	if err != nil {
		return requestErrorf(err, "Error marshalling json")
	}
	_, err = w.Write(content)
	if err != nil {
		return requestErrorf(err, "Error writing response")
	}
	return nil
}

func getUserFromContext(r *http.Request) (*models.User, error) {
	if env.IsLocal() {
		// if local use a generic user
		return getUserByID(12345)
	}

	// get the profile info from the users session
	profile := profileFromSession(r)
	if profile != nil {
		ID, err := strconv.ParseInt(profile.Id[2:], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing google+ profile ID: %v", err)
		}

		return getUserByID(ID)
	}

	ghUser := githubUserFromSession(r)
	if ghUser != nil {
		return getUserByID(int64(*ghUser.ID))
	}

	return nil, fmt.Errorf("No profiles found in session")
}

func getUserByID(ID int64) (*models.User, error) {
	user, err := cp.GetUser(ID)
	logger.Infof("user %v\nerr %v", user, err)
	if err != nil {
		// assume user entity does not exist for this profile
		// TODO: find a way to only do this if entity dne

		return createUser(ID)
	}

	return user, nil
}

func createUser(ID int64) (*models.User, error) {
	user := &models.User{
		ID:         ID,
		ProjectIDs: []int64{},
	}
	err := cp.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %v", err)
	}

	return user, nil
}
