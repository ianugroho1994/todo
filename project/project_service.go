package project

import (
	"context"
	"errors"

	"github.com/ianugroho1994/todo/shared"
	"github.com/jackc/pgx/v5"
)

type ProjectService interface {
	ListProjectsByGroup(ctx context.Context, groupID string) ([]*ProjectItem, error)
	CreateProject(ctx context.Context, title string, groupID string) (*ProjectItem, error)
	UpdateProject(ctx context.Context, id string, title, groupID string) (*ProjectItem, error)
	DeleteProject(ctx context.Context, id string) error
}

type ProjectServiceImpl struct {
	projectRepository ProjectRepository
}

func NewProjectService(projectRepo ProjectRepository) ProjectService {
	return &ProjectServiceImpl{
		projectRepository: projectRepo,
	}
}

func (s *ProjectServiceImpl) ListProjectsByGroup(ctx context.Context, groupID string) ([]*ProjectItem, error) {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.projectRepository.GetByGroupID(ctx, tx, groupID)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get project by group id")
		tx.Rollback(ctx)
		return nil, err
	}

	return res, tx.Commit(ctx)
}

func (s *ProjectServiceImpl) CreateProject(ctx context.Context, title string, groupID string) (*ProjectItem, error) {
	projectItem, err := NewProjectItem(
		title,
		groupID)

	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to create project item")
		return nil, err
	}

	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.projectRepository.Store(ctx, tx, projectItem)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store project item")
		tx.Rollback(ctx)
		return nil, err
	}

	return projectItem, tx.Commit(ctx)
}

func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, id string, title, groupID string) (*ProjectItem, error) {
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
	res.GroupID = groupID

	err = s.projectRepository.Store(ctx, tx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store project item")
		tx.Rollback(ctx)
		return nil, err
	}

	return res, tx.Commit(ctx)
}

func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, id string) error {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = s.projectRepository.Delete(ctx, tx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to delete project item")
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *ProjectServiceImpl) fetchByID(ctx context.Context, tx pgx.Tx, id string) (*ProjectItem, error) {
	res, err := s.projectRepository.GetByID(ctx, tx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get task by project id")
		return nil, err
	}

	if res == nil {
		return nil, errors.New("todo: task not found")
	}

	return res, nil
}
