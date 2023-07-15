package project

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ianugroho1994/todo/shared"
	"github.com/oklog/ulid/v2"
)

var (
	projectService ProjectService
)

func ProjectRouters() *chi.Mux {
	initTaskRouter()

	r := chi.NewMux()

	r.Route("/projects", func(r chi.Router) {
		r.Get("/{group_id}", listProjectsByGroupHandler)
		r.Post("/", createProjectHandler)
		r.Put("/{id}", createProjectHandler)
		r.Delete("/{id}", deleteProjectHandler)
	})

	return r
}

func initTaskRouter() {
	projectService = NewProjectService(NewProjectRepository())
}

func listProjectsByGroupHandler(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	itemId := chi.URLParam(request, "group_id")
	id, err := ulid.Parse(itemId)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := projectService.ListProjectsByGroup(ctx, id.String())
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

	itemId := chi.URLParam(request, "id")
	id, err := ulid.Parse(itemId)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	title := request.FormValue("title")
	description := request.FormValue("description")
	links := request.FormValue("links")
	projectIDForm := request.FormValue("project_id")
	projectID, err := ulid.Parse(projectIDForm)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := projectService.UpdateProject(ctx, id.String(), title, description, links, projectID.String())
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

	itemId := chi.URLParam(r, "id")
	id, err := ulid.Parse(itemId)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = projectService.DeleteProject(ctx, id.String())
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
