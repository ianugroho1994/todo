package group

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type GroupItem struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt null.Time `json:"updated_at" db:"updated_at"`
}

func NewGroupItem(name string) (*GroupItem, error) {
	if name == "" {
		return nil, errors.New("group-item: name is required")
	}

	return &GroupItem{
		ID:   ulid.Make().String(),
		Name: name,
	}, nil
}
