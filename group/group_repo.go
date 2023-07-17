package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/ianugroho1994/todo/shared"
	"github.com/jackc/pgx/v5"
)

type GroupRepository interface {
	Store(ctx context.Context, tx pgx.Tx, task *GroupItem) error
	GetAll(ctx context.Context, tx pgx.Tx) ([]*GroupItem, error)
	GetByID(ctx context.Context, tx pgx.Tx, id string) (*GroupItem, error)
	Delete(ctx context.Context, tx pgx.Tx, id string) error
}

type GroupRepositoryForProject interface {
	GetParentGroupID(ctx context.Context, tx pgx.Tx) (string, error)
}

type GroupRepositoryImpl struct{}

func NewGroupRepository() GroupRepository {
	return &GroupRepositoryImpl{}
}

func NewProjectRepositoryForTask() GroupRepositoryForProject {
	return &GroupRepositoryImpl{}
}

func (r *GroupRepositoryImpl) Store(ctx context.Context, tx pgx.Tx, project *GroupItem) error {
	query := `
	INSERT INTO groups (id, name) VALUES ($1, $2)
	ON CONFLICT(id)
	DO UPDATE SET name = $2`

	res, err := tx.Exec(ctx, query, project.ID, project.Name)
	if err != nil {
		return err
	}

	affected := res.RowsAffected()
	if affected != 1 {
		return errors.New("group - repo: failed to store group")
	}

	return nil
}

func (r *GroupRepositoryImpl) GetParentGroupID(ctx context.Context, tx pgx.Tx) (string, error) {
	parentTitle := "parent"
	query := `SELECT * FROM groups WHERE name = $1`
	res, err := r.fetch(ctx, tx, query, parentTitle)
	if err != nil {
		return "", err
	}

	if len(res) <= 0 {
		return "", nil
	}

	return res[0].ID, nil
}

func (r *GroupRepositoryImpl) GetByID(ctx context.Context, tx pgx.Tx, id string) (*GroupItem, error) {
	query := `SELECT * FROM groups WHERE id = $1`
	res, err := r.fetch(ctx, tx, query, id)
	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res[0], nil
}

func (r *GroupRepositoryImpl) GetAll(ctx context.Context, tx pgx.Tx) ([]*GroupItem, error) {
	query := `SELECT * FROM groups ORDER BY created_at DESC`
	res, err := r.fetch(ctx, tx, query)
	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res, nil
}

func (r *GroupRepositoryImpl) fetch(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (result []*GroupItem, err error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		t := &GroupItem{}
		err = rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (r *GroupRepositoryImpl) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM groups WHERE id = $1`

	res, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	affected := res.RowsAffected()
	shared.Log.Info().Msg(fmt.Sprintf("Delete affected: %d", affected))
	return nil
}
