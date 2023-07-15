package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/ianugroho1994/todo/shared"

	"github.com/jmoiron/sqlx"
)

type GroupRepository interface {
	Store(ctx context.Context, task *GroupItem) error
	GetAll(ctx context.Context) ([]*GroupItem, error)
	GetByID(ctx context.Context, id string) (*GroupItem, error)
	Delete(ctx context.Context, id string) error
}

type GroupRepositoryImpl struct {
	DBConnection *sqlx.DB
}

func NewGroupRepository() GroupRepository {
	return &GroupRepositoryImpl{
		DBConnection: shared.DBConnection,
	}
}

func (r *GroupRepositoryImpl) Store(ctx context.Context, project *GroupItem) error {
	query := `INSERT INTO groups (id, name) VALUES (?, ?)
	ON CONFLICT(id)
	DO UPDATE SET name = $2`

	res, err := r.DBConnection.ExecContext(ctx, query, project.ID, project.Name)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("todo: failed to store group")
	}

	return nil
}

func (r *GroupRepositoryImpl) GetByID(ctx context.Context, id string) (*GroupItem, error) {
	query := `SELECT * FROM groups WHERE id = ?`
	res, err := r.fetch(ctx, query, id)
	if err != nil {
		err = errors.New("todo: failed to fetch group")
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res[0], nil
}

func (r *GroupRepositoryImpl) GetAll(ctx context.Context) ([]*GroupItem, error) {
	query := `SELECT * FROM groups ORDER BY created_at DESC`
	res, err := r.fetch(ctx, query)
	if err != nil {
		err = errors.New("todo: failed to fetch task")
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res, nil
}

func (r *GroupRepositoryImpl) fetch(ctx context.Context, query string, args ...interface{}) (result []*GroupItem, err error) {
	rows, err := r.DBConnection.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		t := &GroupItem{}
		err = rows.StructScan(t)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (r *GroupRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM groups WHERE id = ?`

	res, err := r.DBConnection.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("Delete affected: %d", affect)
	return nil
}
