package project

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

type ProjectItem struct {
	ID        ulid.ULID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	GroupID   string    `json:"group_id" db:"group_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewProjectItem(title, groupID string) (*ProjectItem, error) {
	if title == "" {
		return nil, errors.New("todo: title is required")
	}

	return &ProjectItem{
		ID:      ulid.Make(),
		Title:   title,
		GroupID: groupID,
	}, nil
}

func (t *ProjectItem) AssignToGroup(groupID string) error {
	t.GroupID = groupID
	return nil
}
