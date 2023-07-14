package task

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func TaskRouters() *chi.Mux {
	r := chi.NewMux()

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/{project_id}", listTasksByIDHandler)
		r.Get("/{id}", getTaskByIDHandler)
		r.Post("/", createTaskHandler)
		r.Put("/{id}", updateTaskHandler)
		r.Delete("/{id}", deleteTaskHandler)
		r.Put("/{id}/done", makeTaskDoneHandler)
		r.Put("/{id}/todo", makeTaskTodoHandler)
	})

	return r
}

func listTasksByIDHandler(w http.ResponseWriter, r *http.Request) {}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {}

func makeTaskDoneHandler(w http.ResponseWriter, r *http.Request) {}

func makeTaskTodoHandler(w http.ResponseWriter, r *http.Request) {}
