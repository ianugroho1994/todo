package task

import (
	"context"
	"errors"

	"github.com/ianugroho1994/todo/shared"
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
	res, err := s.taskRepository.GetByProjectID(ctx, projectID)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get task by project id")
		return nil, err
	}

	return res, nil
}

func (s *TaskServiceImpl) GetTaskByID(ctx context.Context, id string) (*TaskItem, error) {
	res, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get task by project id")
		return nil, err
	}

	return res, nil
}

func (s *TaskServiceImpl) CreateTask(ctx context.Context, title string, description string, links string, projectID string) (*TaskItem, error) {
	todoItem, err := NewTaskItem(title,
		description,
		links,
		projectID)

	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to create task item")
		return nil, err
	}

	err = s.taskRepository.Store(ctx, todoItem)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store task item")
		return nil, err
	}

	return todoItem, nil
}

func (s *TaskServiceImpl) UpdateTask(ctx context.Context, id string, title, description string, links string, projectID string) (*TaskItem, error) {
	res, err := s.fetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res.Title = title
	res.Description = description
	res.Link = links
	res.ProjectID = projectID

	err = s.taskRepository.Store(ctx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store task item")
		return nil, err
	}

	return res, nil
}

func (s *TaskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	err := s.taskRepository.Delete(ctx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to delete task item")
		return err
	}

	return nil
}

func (s *TaskServiceImpl) MakeTaskDone(ctx context.Context, id string) error {
	res, err := s.fetchByID(ctx, id)
	if err != nil {
		return err
	}

	if err = res.MakeDone(); err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to make task done")
		return err
	}

	err = s.taskRepository.Store(ctx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store task item")
		return err
	}

	return nil
}

func (s *TaskServiceImpl) MakeTaskTodo(ctx context.Context, id string) error {
	res, err := s.fetchByID(ctx, id)
	if err != nil {
		return err
	}

	if err = res.MakeDone(); err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to make task done")
		return err
	}

	err = s.taskRepository.Store(ctx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store task item")
		return err
	}

	return nil
}

func (s *TaskServiceImpl) fetchByID(ctx context.Context, id string) (*TaskItem, error) {
	res, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get task by project id")
		return nil, err
	}

	if res == nil {
		return nil, errors.New("todo: task not found")
	}

	return res, nil
}
