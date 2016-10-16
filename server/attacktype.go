package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codegp/game-object-types/types"
)

func GetAttackType(w http.ResponseWriter, r *http.Request) *requestError {
	attackTypeID, rerr := readIDFromRequest(r, "attackTypeID")
	if rerr != nil {
		return rerr
	}

	attackType, err := cp.GetAttackType(attackTypeID)
	if err != nil {
		return requestErrorf(err, "Error getting attackType from datastore")
	}

	return marshalAndWriteResponse(w, attackType)
}

func PostAttackType(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling attackType")
	}

	var attackType *types.AttackType
	err = json.Unmarshal(body, &attackType)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling attackType")
	}

	attackType, err = cp.AddAttackType(attackType)
	if err != nil {
		return requestErrorf(err, "Error updating attackType in datastore")
	}

	return marshalAndWriteResponse(w, attackType)
}
