package group

import (
	"context"
	"errors"

	"github.com/ianugroho1994/todo/shared"
	"github.com/jackc/pgx/v5"
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
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.groupRepository.GetAll(ctx, tx)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get groups")
		tx.Rollback(ctx)
		return nil, err
	}

	return res, tx.Commit(ctx)
}

func (s *GroupServiceImpl) CreateGroup(ctx context.Context, name string) (*GroupItem, error) {
	groupItem, err := NewGroupItem(
		name)

	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to create group item")
		return nil, err
	}

	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	err = s.groupRepository.Store(ctx, tx, groupItem)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store group item")
		tx.Rollback(ctx)
		return nil, err
	}

	return groupItem, tx.Commit(ctx)
}

func (s *GroupServiceImpl) UpdateGroup(ctx context.Context, id string, name string) (*GroupItem, error) {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.fetchByID(ctx, tx, id)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	res.Name = name

	err = s.groupRepository.Store(ctx, tx, res)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to store group item")
		tx.Rollback(ctx)
		return nil, err
	}

	return res, tx.Commit(ctx)
}

func (s *GroupServiceImpl) DeleteGroup(ctx context.Context, id string) error {
	tx, err := shared.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = s.groupRepository.Delete(ctx, tx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to delete group item")
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (s *GroupServiceImpl) fetchByID(ctx context.Context, tx pgx.Tx, id string) (*GroupItem, error) {
	res, err := s.groupRepository.GetByID(ctx, tx, id)
	if err != nil {
		shared.Log.Error().Err(err).Msg("todo: failed to get group by id")
		return nil, err
	}

	if res == nil {
		return nil, errors.New("todo: group not found")
	}

	return res, nil
}
