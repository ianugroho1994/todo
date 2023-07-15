package task

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ianugroho1994/todo/shared"
	"github.com/oklog/ulid/v2"
)

var (
	taskService TaskService
)

func TaskRouters() *chi.Mux {
	initTaskRouter()

	r := chi.NewMux()

	r.Get("/{project_id}", listTasksByProjectHandler)
	r.Get("/{id}", getTaskByIDHandler)
	r.Post("/", createTaskHandler)
	r.Put("/{id}", createTaskHandler)
	r.Delete("/{id}", deleteTaskHandler)
	r.Put("/{id}/done", makeTaskDoneHandler)
	r.Put("/{id}/todo", makeTaskTodoHandler)

	return r
}

func initTaskRouter() {
	taskService = NewTaskService(NewTaskRepository())
}

func listTasksByProjectHandler(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	itemId := chi.URLParam(request, "project_id")
	id, err := ulid.Parse(itemId)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := taskService.ListTasksByProject(ctx, id.String())
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func getTaskByIDHandler(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	itemId := chi.URLParam(request, "id")
	id, err := ulid.Parse(itemId)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := taskService.GetTaskByID(ctx, id.String())
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func createTaskHandler(w http.ResponseWriter, request *http.Request) {
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

	resp, err := taskService.UpdateTask(ctx, id.String(), title, description, links, projectID.String())
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemId := chi.URLParam(r, "id")
	id, err := ulid.Parse(itemId)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = taskService.DeleteTask(ctx, id.String())
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func makeTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemId := chi.URLParam(r, "id")
	id, err := ulid.Parse(itemId)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = taskService.MakeTaskDone(ctx, id.String())
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func makeTaskTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemId := chi.URLParam(r, "id")
	id, err := ulid.Parse(itemId)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = taskService.MakeTaskTodo(ctx, id.String())
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
