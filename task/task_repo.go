package task

import (
	"context"
	"errors"
	"fmt"

	"github.com/ianugroho1994/todo/shared"

	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	Store(ctx context.Context, task *TaskItem) error
	GetByID(ctx context.Context, id string) (*TaskItem, error)
	GetByProjectID(ctx context.Context, projectID string) ([]*TaskItem, error)
	Delete(ctx context.Context, id string) error
}

type TaskRepositoryImpl struct {
	DBConnection *sqlx.DB
}

func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{
		DBConnection: shared.DBConnection,
	}
}

func (r *TaskRepositoryImpl) Store(ctx context.Context, task *TaskItem) error {
	query := `INSERT INTO tasks (id, title, description, links, project_id, is_todo, created_at, done_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id)
	DO UPDATE SET title= ?$2 and done_at= $7 and description= $3 and link= $4 and project_id= $5 and is_todo= $6`

	res, err := r.DBConnection.ExecContext(ctx, query, task.ID, task.Title, task.Description, task.Link, task.ProjectID, task.IsTodo, task.CreatedAt, task.DoneAt)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("todo: failed to store task")
	}

	return nil
}

func (r *TaskRepositoryImpl) GetByID(ctx context.Context, id string) (*TaskItem, error) {
	query := `SELECT * FROM tasks WHERE id = ?`

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

func (r *TaskRepositoryImpl) GetByProjectID(ctx context.Context, projectID string) ([]*TaskItem, error) {
	query := `SELECT * FROM tasks WHERE project_id = ?`
	res, err := r.fetch(ctx, query, projectID)
	if err != nil {
		err = errors.New("todo: failed to fetch task")
		return nil, err
	}

	if len(res) <= 0 {
		return nil, nil
	}

	return res, nil
}

func (r *TaskRepositoryImpl) fetch(ctx context.Context, query string, args ...interface{}) (result []*TaskItem, err error) {
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
		t := &TaskItem{}
		err = rows.StructScan(t)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (r *TaskRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`

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
