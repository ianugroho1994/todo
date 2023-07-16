package task

import (
	"context"
	"errors"
	"fmt"

	"github.com/ianugroho1994/todo/shared"
	"github.com/jackc/pgx/v5"
)

type TaskRepository interface {
	Store(ctx context.Context, tx pgx.Tx, task *TaskItem) error
	GetByID(ctx context.Context, tx pgx.Tx, id string) (*TaskItem, error)
	GetByProjectID(ctx context.Context, tx pgx.Tx, projectID string) ([]*TaskItem, error)
	Delete(ctx context.Context, tx pgx.Tx, id string) error
}

type TaskRepositoryImpl struct{}

func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{}
}

func (r *TaskRepositoryImpl) Store(ctx context.Context, tx pgx.Tx, task *TaskItem) error {
	query := `
	INSERT INTO tasks (id, title, description, links, project_id, is_todo, created_at, done_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT(id)
	DO UPDATE SET title = $2, description = $3, link = $4, project_id = $5, is_todo = $6, done_at = $8`

	res, err := tx.Exec(ctx, query, task.ID, task.Title, task.Description, task.Link, task.ProjectID, task.IsTodo, task.CreatedAt, task.DoneAt)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_repo: failed to store task ")
		return err
	}

	affected := res.RowsAffected()
	if affected != 1 {
		return errors.New("task_repo: there is no rows affected")
	}

	return nil
}

func (r *TaskRepositoryImpl) GetByID(ctx context.Context, tx pgx.Tx, id string) (*TaskItem, error) {
	query := `SELECT * FROM tasks WHERE id = $1`

	res, err := r.fetch(ctx, tx, query, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_repo: failed to fetch task by id: " + id)
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res[0], nil
}

func (r *TaskRepositoryImpl) GetByProjectID(ctx context.Context, tx pgx.Tx, projectID string) ([]*TaskItem, error) {
	query := `SELECT * FROM tasks WHERE project_id = $1`
	res, err := r.fetch(ctx, tx, query, projectID)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_repo: failed to fetch task by project id: " + projectID)
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res, nil
}

func (r *TaskRepositoryImpl) fetch(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (result []*TaskItem, err error) {
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
		t := &TaskItem{}
		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Link, &t.ProjectID, &t.IsTodo, &t.CreatedAt, &t.DoneAt)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (r *TaskRepositoryImpl) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	res, err := tx.Exec(ctx, query, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_repo: failed to delete task by id: " + id)
		return err
	}
	affected := res.RowsAffected()
	shared.Log.Info().Msg(fmt.Sprintf("task_repo: delete affected: %d", affected))
	return nil
}
