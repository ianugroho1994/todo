package group

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ianugroho1994/todo/shared"
)

var (
	groupService GroupService
)

func GroupRouters() *chi.Mux {
	initGroupRouter()

	r := chi.NewMux()

	r.Get("/", listGroupsByGroupHandler)
	r.Post("/", createGroupHandler)
	r.Put("/{id}", createGroupHandler)
	r.Delete("/{id}", deleteProjectHandler)

	return r
}

func initGroupRouter() {
	groupService = NewGroupService(NewGroupRepository())
}

func listGroupsByGroupHandler(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	resp, err := groupService.ListAllGroup(ctx)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func createGroupHandler(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := request.Context()

	name := request.FormValue("name")
	resp, err := groupService.CreateGroup(ctx, name)
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

	err := groupService.DeleteGroup(ctx, id)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
