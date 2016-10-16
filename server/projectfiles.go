package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func getTemplate(lang string) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("templates/%s", lang))
}

func GetProjectFile(w http.ResponseWriter, r *http.Request) *requestError {
	vars := mux.Vars(r)
	fileName := vars["name"]
	projectID, rerr := readIDFromRequest(r, "projectID")
	if rerr != nil {
		return rerr
	}

	content, err := cp.ReadProjectFile(projectID, fileName)
	if err != nil {
		return requestErrorf(err, "Error reading file from storage")
	}

	_, err = w.Write(content)
	if err != nil {
		return requestErrorf(err, "Error writing response")
	}

	return nil
}

func PostProjectFile(w http.ResponseWriter, r *http.Request) *requestError {
	vars := mux.Vars(r)
	fileName := vars["name"]
	projectID, rerr := readIDFromRequest(r, "projectID")
	if rerr != nil {
		return rerr
	}

	project, err := cp.GetProject(projectID)
	if err != nil {
		return requestErrorf(err, "Error getting project from datastore")
	}

	project.FileNames = append(project.FileNames, fileName)
	err = cp.UpdateProject(project)
	if err != nil {
		return requestErrorf(err, "Error putting project to datastore")
	}

	template, err := getTemplate(project.Language)
	if err != nil {
		return requestErrorf(err, "Error reading template")
	}

	err = cp.WriteProjectFile(projectID, fileName, template)
	if err != nil {
		return requestErrorf(err, "Error writing file to storage")
	}

	_, err = w.Write(template)
	if err != nil {
		return requestErrorf(err, "Error writing file to response")
	}

	return nil
}

func PutProjectFile(w http.ResponseWriter, r *http.Request) *requestError {
	vars := mux.Vars(r)
	fileName := vars["name"]
	projectID, rerr := readIDFromRequest(r, "projectID")
	if rerr != nil {
		return rerr
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error reading file from body")
	}

	err = cp.WriteProjectFile(projectID, fileName, body)
	if err != nil {
		return requestErrorf(err, "Error writing file to storage")
	}

	return nil
}
