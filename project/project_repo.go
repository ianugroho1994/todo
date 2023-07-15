package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/ianugroho1994/todo/shared"

	"github.com/jmoiron/sqlx"
)

type ProjectRepository interface {
	Store(ctx context.Context, task *ProjectItem) error
	GetByGroupID(ctx context.Context, projectID string) ([]*ProjectItem, error)
	GetByID(ctx context.Context, id string) (*ProjectItem, error)
	Delete(ctx context.Context, id string) error
}

type ProjectRepositoryImpl struct {
	DBConnection *sqlx.DB
}

func NewProjectRepository() ProjectRepository {
	return &ProjectRepositoryImpl{
		DBConnection: shared.DBConnection,
	}
}

func (r *ProjectRepositoryImpl) Store(ctx context.Context, project *ProjectItem) error {
	query := `INSERT INTO projects (id, title, group_id) VALUES (?, ?, ?)
	ON CONFLICT(id)
	DO UPDATE SET title= ?$2 and done_at= $7 and description= $3 and link= $4 and project_id= $5 and is_todo= $6`

	res, err := r.DBConnection.ExecContext(ctx, query, project.ID, project.Title, project.GroupID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("todo: failed to store project")
	}

	return nil
}

func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, id string) (*ProjectItem, error) {
	query := `SELECT * FROM projects WHERE id = ?`
	res, err := r.fetch(ctx, query, id)
	if err != nil {
		err = errors.New("todo: failed to fetch task")
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res[0], nil
}

func (r *ProjectRepositoryImpl) GetByGroupID(ctx context.Context, groupID string) ([]*ProjectItem, error) {
	query := `SELECT * FROM projects WHERE group_id = ?`
	res, err := r.fetch(ctx, query, groupID)
	if err != nil {
		err = errors.New("todo: failed to fetch task")
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res, nil
}

func (r *ProjectRepositoryImpl) fetch(ctx context.Context, query string, args ...interface{}) (result []*ProjectItem, err error) {
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
		t := &ProjectItem{}
		err = rows.StructScan(t)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (r *ProjectRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = ?`

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