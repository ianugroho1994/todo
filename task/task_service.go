package task

import (
	"context"
	"errors"

	"github.com/ianugroho1994/todo/shared"
	"github.com/jackc/pgx/v5"
)

type TaskService interface {
	ListTasksByProject(ctx context.Context, projectID string) ([]*TaskItem, error)
	GetTaskByID(ctx context.Context, id string) (*TaskItem, error)
	CreateTask(ctx context.Context, title string, description string, links string, projectID string) (*TaskItem, error)
	UpdateTask(ctx context.Context, id string, title, description string, links string, projectID string) (*TaskItem, error)
	DeleteTask(ctx context.Context, id string) error
	MakeTaskDone(ctx context.Context, id string) error
	MakeTaskTodo(ctx context.Context, id string) error
}

type TaskServiceImpl struct {
	taskRepository TaskRepository
}

func NewTaskService(taskRepo TaskRepository) TaskService {
	return &TaskServiceImpl{
		taskRepository: taskRepo,
	}
}

func (s *TaskServiceImpl) ListTasksByProject(ctx context.Context, projectID string) ([]*TaskItem, error) {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.taskRepository.GetByProjectID(ctx, tx, projectID)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to get task by project id")
		tx.Rollback(ctx)
		return nil, err
	}

	return res, tx.Commit(ctx)
}

func (s *TaskServiceImpl) GetTaskByID(ctx context.Context, id string) (*TaskItem, error) {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.taskRepository.GetByID(ctx, tx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to get task by project id")
		tx.Rollback(ctx)
		return nil, err
	}

	return res, tx.Commit(ctx)
}

func (s *TaskServiceImpl) CreateTask(ctx context.Context, title string, description string, links string, projectID string) (*TaskItem, error) {
	todoItem, err := NewTaskItem(title,
		description,
		links,
		projectID)

	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to create task item")
		return nil, err
	}

	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.taskRepository.Store(ctx, tx, todoItem)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to store task item")
		tx.Rollback(ctx)
		return nil, err
	}

	return todoItem, tx.Commit(ctx)
}

func (s *TaskServiceImpl) UpdateTask(ctx context.Context, id string, title, description string, links string, projectID string) (*TaskItem, error) {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.fetchByID(ctx, tx, id)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	res.Title = title
	res.Description = description
	res.Link = links
	res.ProjectID = projectID

	err = s.taskRepository.Store(ctx, tx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to store task item")
		tx.Rollback(ctx)
		return nil, err
	}

	return res, tx.Commit(ctx)
}

func (s *TaskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = s.taskRepository.Delete(ctx, tx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to delete task item")
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *TaskServiceImpl) MakeTaskDone(ctx context.Context, id string) error {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	res, err := s.fetchByID(ctx, tx, id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err = res.MakeDone(); err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to make task done")
		tx.Rollback(ctx)
		return err
	}

	err = s.taskRepository.Store(ctx, tx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to store task item")
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *TaskServiceImpl) MakeTaskTodo(ctx context.Context, id string) error {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	res, err := s.fetchByID(ctx, tx, id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err = res.MakeAsTodo(); err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to make task done")
		tx.Rollback(ctx)
		return err
	}

	err = s.taskRepository.Store(ctx, tx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to store task item")
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *TaskServiceImpl) fetchByID(ctx context.Context, tx pgx.Tx, id string) (*TaskItem, error) {
	res, err := s.taskRepository.GetByID(ctx, tx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("task_service: failed to get task by project id")
		return nil, err
	}

	if res == nil {
		return nil, errors.New("task_service: task not found")
	}

	return res, nil
}
