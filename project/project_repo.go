package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/ianugroho1994/todo/shared"
	"github.com/jackc/pgx/v5"
)

type ProjectRepository interface {
	Store(ctx context.Context, tx pgx.Tx, task *ProjectItem) error
	GetByGroupID(ctx context.Context, tx pgx.Tx, projectID string) ([]*ProjectItem, error)
	GetByID(ctx context.Context, tx pgx.Tx, id string) (*ProjectItem, error)
	Delete(ctx context.Context, tx pgx.Tx, id string) error
}

type ProjectRepositoryForTask interface {
	GetParentProjectID(ctx context.Context, tx pgx.Tx) (string, error)
}

type ProjectRepositoryImpl struct{}

func NewProjectRepository() ProjectRepository {
	return &ProjectRepositoryImpl{}
}

func NewProjectRepositoryForTask() ProjectRepositoryForTask {
	return &ProjectRepositoryImpl{}
}

func (r *ProjectRepositoryImpl) Store(ctx context.Context, tx pgx.Tx, project *ProjectItem) error {
	query := `
	INSERT INTO projects (id, title, group_id) VALUES ($1, $2, $3)
	ON CONFLICT(id)
	DO UPDATE SET title= $2, group_id= $3`

	res, err := tx.Exec(ctx, query, project.ID, project.Title, project.GroupID)
	if err != nil {
		shared.Log.Err(err).Msg("project-repo: failed to store project")
		return err
	}

	affected := res.RowsAffected()

	if affected != 1 {
		return errors.New("projects-repo: failed to store project, no rows affected")
	}

	return nil
}

func (r *ProjectRepositoryImpl) GetParentProjectID(ctx context.Context, tx pgx.Tx) (string, error) {
	parentTitle := "parent"
	query := `SELECT * FROM projects WHERE title = $1`
	res, err := r.fetch(ctx, tx, query, parentTitle)
	if err != nil {
		return "", err
	}

	if len(res) <= 0 {
		return "", nil
	}

	return res[0].ID, nil
}

func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, tx pgx.Tx, id string) (*ProjectItem, error) {
	query := `SELECT * FROM projects WHERE id = $1`
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
	shared.Log.Info().Msg(fmt.Sprintf("Delete affected: %d", affected))
	return nil
}
