package task

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type TaskItem struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description" db:"description"`
	Link        *string   `json:"link" db:"link"`
	ProjectID   string    `json:"project_id" db:"project_id"`
	IsTodo      bool      `json:"is_todo" db:"is_todo"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	DoneAt      null.Time `json:"done_at" db:"done_at"`
}

func NewTaskItem(title, description string, links string, projectID string) (*TaskItem, error) {
	if title == "" {
		return nil, errors.New("todo: title is required")
	}

	return &TaskItem{
		ID:          ulid.Make().String(),
		Title:       title,
		Description: description,
		Link:        links,
		ProjectID:   projectID,
		IsTodo:      false,
		CreatedAt:   time.Now(),
	}, nil
}

func (t *TaskItem) IsDone() bool {
	return t.DoneAt.Valid && t.DoneAt.Time.After(t.CreatedAt)
}

func (t *TaskItem) MakeDone() error {
	if t.IsDone() {
		return errors.New("todo: the item is done")
	}

	t.DoneAt = null.TimeFrom(time.Now())
	return nil
}

func (t *TaskItem) MakeAsTodo() error {
	if !t.IsDone() {
		return errors.New("todo: the item is todo")
	}

	t.IsTodo = true
	return nil
}

func (t *TaskItem) AssignToProject(projectID string) {
	t.ProjectID = projectID
}
