package group

import (
	"context"
	"errors"

	"github.com/ianugroho1994/todo/shared"
)

type GroupService interface {
	ListAllGroup(ctx context.Context) ([]*GroupItem, error)
	CreateGroup(ctx context.Context, name string) (*GroupItem, error)
	UpdateGroup(ctx context.Context, id string, name string) (*GroupItem, error)
	DeleteGroup(ctx context.Context, id string) error
}

type GroupServiceImpl struct {
	groupRepository GroupRepository
}

func NewGroupService(groupRepo GroupRepository) GroupService {
	return &GroupServiceImpl{
		groupRepository: groupRepo,
	}
}

func (s *GroupServiceImpl) ListAllGroup(ctx context.Context) ([]*GroupItem, error) {
	res, err := s.groupRepository.GetAll(ctx)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get groups")
		return nil, err
	}

	return res, nil
}

func (s *GroupServiceImpl) CreateGroup(ctx context.Context, name string) (*GroupItem, error) {
	groupItem, err := NewGroupItem(
		name)

	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to create group item")
		return nil, err
	}

	err = s.groupRepository.Store(ctx, groupItem)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store group item")
		return nil, err
	}

	return groupItem, nil
}

func (s *GroupServiceImpl) UpdateGroup(ctx context.Context, id string, name string) (*GroupItem, error) {
	res, err := s.fetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res.Name = name

	err = s.groupRepository.Store(ctx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store group item")
		return nil, err
	}

	return res, nil
}

func (s *GroupServiceImpl) DeleteGroup(ctx context.Context, id string) error {
	err := s.groupRepository.Delete(ctx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to delete group item")
		return err
	}

	return nil
}

func (s *GroupServiceImpl) fetchByID(ctx context.Context, id string) (*GroupItem, error) {
	res, err := s.groupRepository.GetByID(ctx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get group by id")
		return nil, err
	}

	if res == nil {
		return nil, errors.New("todo: group not found")
	}

	return res, nil
}
