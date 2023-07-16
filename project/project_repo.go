package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ProjectRepository interface {
	Store(ctx context.Context, tx pgx.Tx, task *ProjectItem) error
	GetByGroupID(ctx context.Context, tx pgx.Tx, projectID string) ([]*ProjectItem, error)
	GetByID(ctx context.Context, tx pgx.Tx, id string) (*ProjectItem, error)
	Delete(ctx context.Context, tx pgx.Tx, id string) error
}

type ProjectRepositoryImpl struct{}

func NewProjectRepository() ProjectRepository {
	return &ProjectRepositoryImpl{}
}

func (r *ProjectRepositoryImpl) Store(ctx context.Context, tx pgx.Tx, project *ProjectItem) error {
	query := `INSERT INTO projects (id, title, group_id) VALUES ($1, $2, $3)
	ON CONFLICT(id)
	DO UPDATE SET title= ?$2 and group_id= $3`

	res, err := tx.Exec(ctx, query, project.ID, project.Title, project.GroupID)
	if err != nil {
		return err
	}

	affected := res.RowsAffected()

	if affected != 1 {
		return errors.New("todo: failed to store project")
	}

	return nil
}

func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, tx pgx.Tx, id string) (*ProjectItem, error) {
	query := `SELECT * FROM projects WHERE id = ?`
	res, err := r.fetch(ctx, tx, query, id)
	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res[0], nil
}

func (r *ProjectRepositoryImpl) GetByGroupID(ctx context.Context, tx pgx.Tx, groupID string) ([]*ProjectItem, error) {
	query := `SELECT * FROM projects WHERE group_id = $1`
	res, err := r.fetch(ctx, tx, query, groupID)
	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res, nil
}

func (r *ProjectRepositoryImpl) fetch(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (result []*ProjectItem, err error) {
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
		t := &ProjectItem{}
		err = rows.Scan(&t.ID, &t.Title, &t.GroupID, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (r *ProjectRepositoryImpl) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM projects WHERE id = $1`

	res, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	affected := res.RowsAffected()
	fmt.Println("Delete affected: %d", affected)
	return nil
}
