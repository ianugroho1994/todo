package project

import (
	"context"
	"errors"

	"github.com/ianugroho1994/todo/shared"
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
	res, err := s.projectRepository.GetByGroupID(ctx, groupID)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get project by group id")
		return nil, err
	}

	return res, nil
}

func (s *ProjectServiceImpl) CreateProject(ctx context.Context, title string, groupID string) (*ProjectItem, error) {
	projectItem, err := NewProjectItem(
		title,
		groupID)

	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to create project item")
		return nil, err
	}

	err = s.projectRepository.Store(ctx, projectItem)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store project item")
		return nil, err
	}

	return projectItem, nil
}

func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, id string, title, groupID string) (*ProjectItem, error) {
	res, err := s.fetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res.Title = title
	res.GroupID = groupID

	err = s.projectRepository.Store(ctx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store project item")
		return nil, err
	}

	return res, nil
}

func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, id string) error {
	err := s.projectRepository.Delete(ctx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to delete project item")
		return err
	}

	return nil
}

func (s *ProjectServiceImpl) fetchByID(ctx context.Context, id string) (*ProjectItem, error) {
	res, err := s.projectRepository.GetByID(ctx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get task by project id")
		return nil, err
	}

	if res == nil {
		return nil, errors.New("todo: task not found")
	}

	return res, nil
}
