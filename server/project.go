package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codegp/cloud-persister/models"
)

func GetProject(w http.ResponseWriter, r *http.Request) *requestError {
	projectID, rerr := readIDFromRequest(r, "projectID")
	if rerr != nil {
		return rerr
	}

	project, err := cp.GetProject(projectID)
	if err != nil {
		return requestErrorf(err, "Error getting project from datastore")
	}

	return marshalAndWriteResponse(w, project)
}

func PostProject(w http.ResponseWriter, r *http.Request) *requestError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling project")
	}

	var project *models.Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		return requestErrorf(err, "Error unmarshalling project")
	}

	user, err := getUserFromContext(r)
	if err != nil {
		return requestErrorf(err, "Error getting user from session")
	}
	project.UserID = user.ID
	logger.Debugf("Project %v", project)
	project, err = cp.AddProject(project)
	if err != nil {
		return requestErrorf(err, "Error adding project to datastore")
	}

	user.ProjectIDs = append(user.ProjectIDs, project.ID)
	err = cp.UpdateUser(user)
	if err != nil {
		return requestErrorf(err, "Error updating user to datastore")
	}

	return marshalAndWriteResponse(w, project)
}

func GetProjects(w http.ResponseWriter, r *http.Request) *requestError {
	user, err := getUserFromContext(r)
	if err != nil {
		return requestErrorf(err, "Error getting user from datastore")
	}

	projects := []*models.Project{}
	for _, projectID := range user.ProjectIDs {
		project, err := cp.GetProject(projectID)
		if err != nil {
			return requestErrorf(err, "Error getting project from datastore")
		}
		projects = append(projects, project)
	}

	return marshalAndWriteResponse(w, projects)
}
