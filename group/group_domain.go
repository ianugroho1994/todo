package group

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

type GroupItem struct {
	ID        ulid.ULID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewGroupItem(name string) (*GroupItem, error) {
	if name == "" {
		return nil, errors.New("todo: name is required")
	}

	return &GroupItem{
		ID:   ulid.Make(),
		Name: name,
	}, nil
}