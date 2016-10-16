package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codegp/game-object-types/types"
)

func GetItemType(w http.ResponseWriter, r *http.Request) *requestError {
	itemTypeID, rerr := readIDFromRequest(r, "itemTypeID")
	if rerr != nil {
		return rerr
	}

	itemType, err := cp.GetItemType(itemTypeID)
	if err != nil {
		return requestErrorf(err, "Error getting itemType from datastore")
	}

	return marshalAndWriteResponse(w, itemType)
}

func PostItemType(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling itemType")
	}

	var itemType *types.ItemType
	err = json.Unmarshal(body, &itemType)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling itemType")
	}

	itemType, err = cp.AddItemType(itemType)
	if err != nil {
		return requestErrorf(err, "Error updating itemType in datastore")
	}

	return marshalAndWriteResponse(w, itemType)
}
