package task

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ianugroho1994/todo/shared"
)

var (
	taskService TaskService
)

func TaskRouters() *chi.Mux {
	initTaskRouter()

	r := chi.NewMux()

	r.Get("/project/{project_id}", listTasksByProjectHandler)
	r.Get("/{id}", getTaskByIDHandler)
	r.Post("/", updateTaskHandler)
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

	projectID := chi.URLParam(request, "project_id")
	shared.Log.Info().Msg("task_router: id: " + projectID)

	resp, err := taskService.ListTasksByProject(ctx, projectID)
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

	id := chi.URLParam(request, "id")
	shared.Log.Info().Msg("task_router: id: " + id)

	resp, err := taskService.GetTaskByID(ctx, id)
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

	title := request.FormValue("title")
	description := request.FormValue("description")
	links := request.FormValue("links")
	projectIDForm := request.FormValue("project_id")

	resp, err := taskService.CreateTask(ctx, title, description, links, projectIDForm)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func updateTaskHandler(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := request.Context()

	id := chi.URLParam(request, "id")
	shared.Log.Info().Msg("task_router: id: " + id)

	title := request.FormValue("title")
	description := request.FormValue("description")
	links := request.FormValue("links")
	projectIDForm := request.FormValue("project_id")

	resp, err := taskService.UpdateTask(ctx, id, title, description, links, projectIDForm)
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

	id := chi.URLParam(r, "id")
	shared.Log.Info().Msg("task_router: id: " + id)

	err := taskService.DeleteTask(ctx, id)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func makeTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	shared.Log.Info().Msg("task_router: id: " + id)

	err := taskService.MakeTaskDone(ctx, id)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func makeTaskTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	shared.Log.Info().Msg("task_router: id: " + id)

	err := taskService.MakeTaskTodo(ctx, id)
	if err != nil {
		shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
