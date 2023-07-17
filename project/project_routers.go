package project

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ianugroho1994/todo/group"
	"github.com/ianugroho1994/todo/shared"
)

var (
	projectService ProjectService
)

func ProjectRouters() *chi.Mux {
	initTaskRouter()

	r := chi.NewMux()

	r.Get("/{group_id}", listProjectsByGroupHandler)
	r.Post("/", createProjectHandler)
	r.Put("/{id}", updateProjectHandler)
	r.Delete("/{id}", deleteProjectHandler)

	return r
}

func initTaskRouter() {
	projectService = NewProjectService(
		NewProjectRepository(),
		group.NewProjectRepositoryForTask())
}

func listProjectsByGroupHandler(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	id := chi.URLParam(request, "group_id")

	resp, err := projectService.ListProjectsByGroup(ctx, id)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func createProjectHandler(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := request.Context()

	newProject := &ProjectItem{}
	err := json.NewDecoder(request.Body).Decode(&newProject)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := projectService.CreateProject(ctx, newProject.Title, newProject.GroupID)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func updateProjectHandler(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := request.Context()

	newProject := &ProjectItem{}
	err := json.NewDecoder(request.Body).Decode(&newProject)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := projectService.UpdateProject(ctx, newProject.ID, newProject.Title, newProject.GroupID)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	err := projectService.DeleteProject(ctx, id)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
