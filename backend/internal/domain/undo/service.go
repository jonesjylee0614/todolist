package undo

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"todolist/backend/internal/domain/task"
	"todolist/backend/internal/repository"
)

type Service struct {
	repo     *repository.UndoRepository
	taskRepo *repository.TaskRepository
	ttl      time.Duration
	logger   *zap.Logger
}

func NewService(repo *repository.UndoRepository, taskRepo *repository.TaskRepository, ttl time.Duration, logger *zap.Logger) *Service {
	return &Service{repo: repo, taskRepo: taskRepo, ttl: ttl, logger: logger}
}

func (s *Service) RecordOperation(ctx context.Context, tx interface{}, action task.Action, scope task.Scope, taskIDs []string, before, after []task.Snapshot) (string, error) {
	token := generateToken()
	beforeJSON, err := json.Marshal(before)
	if err != nil {
		return "", err
	}
	afterJSON, err := json.Marshal(after)
	if err != nil {
		return "", err
	}
	idsJSON, err := json.Marshal(taskIDs)
	if err != nil {
		return "", err
	}

	op := &repository.TaskOperation{
		Token:       token,
		Action:      string(action),
		Scope:       string(scope),
		TaskIDs:     string(idsJSON),
		BeforeState: string(beforeJSON),
		AfterState:  string(afterJSON),
		ExpireAt:    time.Now().Add(s.ttl),
	}

	if err := s.repo.Create(ctx, tx.(*gorm.DB), op); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Undo(ctx context.Context, token string) ([]string, string, error) {
	op, err := s.repo.GetByToken(ctx, nil, token)
	if err != nil {
		return nil, "", err
	}
	if op == nil {
		return nil, "", ErrTokenNotFound
	}
	if op.IsConsumed() {
		return nil, "", ErrTokenConsumed
	}
	if op.IsExpired(time.Now()) {
		return nil, "", ErrTokenExpired
	}

	var before []task.Snapshot
	var after []task.Snapshot
	var ids []string

	if err := json.Unmarshal([]byte(op.BeforeState), &before); err != nil {
		return nil, "", err
	}
	if err := json.Unmarshal([]byte(op.AfterState), &after); err != nil {
		return nil, "", err
	}
	if err := json.Unmarshal([]byte(op.TaskIDs), &ids); err != nil {
		return nil, "", err
	}

	var reverseToken string
	err = s.taskRepo.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.applyUndo(ctx, tx, task.Action(op.Action), before, after); err != nil {
			return err
		}
		if err := s.repo.MarkConsumed(ctx, tx, token, time.Now()); err != nil {
			return err
		}
		newAction := reverseAction(task.Action(op.Action))
		newToken, err := s.RecordOperation(ctx, tx, newAction, task.Scope(op.Scope), ids, after, before)
		if err != nil {
			return err
		}
		reverseToken = newToken
		return nil
	})
	if err != nil {
		return nil, "", err
	}

	return ids, reverseToken, nil
}

func (s *Service) applyUndo(ctx context.Context, tx *gorm.DB, action task.Action, before, after []task.Snapshot) error {
	switch action {
	case task.ActionCreate:
		return s.taskRepo.DeleteBySnapshots(ctx, tx, after)
	case task.ActionDelete, task.ActionBulkDelete:
		return s.taskRepo.ReplaceSnapshots(ctx, tx, before)
	case task.ActionMove, task.ActionComplete, task.ActionUpdate, task.ActionBulkMove, task.ActionBulkComplete, task.ActionResort:
		return s.taskRepo.ReplaceSnapshots(ctx, tx, before)
	default:
		return errors.New("unsupported action for undo")
	}
}

func reverseAction(action task.Action) task.Action {
	switch action {
	case task.ActionCreate:
		return task.ActionDelete
	case task.ActionDelete:
		return task.ActionCreate
	case task.ActionMove:
		return task.ActionMove
	case task.ActionComplete:
		return task.ActionComplete
	case task.ActionUpdate:
		return task.ActionUpdate
	case task.ActionBulkMove:
		return task.ActionBulkMove
	case task.ActionBulkComplete:
		return task.ActionBulkComplete
	case task.ActionBulkDelete:
		return task.ActionBulkDelete
	case task.ActionResort:
		return task.ActionResort
	default:
		return action
	}
}

func generateToken() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")[:26]
}

var (
	ErrTokenNotFound = errors.New("undo token not found")
	ErrTokenConsumed = errors.New("undo token already consumed")
	ErrTokenExpired  = errors.New("undo token expired")
)
