package main

import "net/http"

func GetUser(w http.ResponseWriter, r *http.Request) *requestError {
	logger.Info("GetUser")
	user, err := getUserFromContext(r)
	if err != nil {
		return requestErrorf(err, "Error getting user from context")
	}

	return marshalAndWriteResponse(w, user)
}
