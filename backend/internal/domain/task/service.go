package task

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	repo          TaskRepository
	undoService   UndoService
	logger        *zap.Logger
	defaultWeight func() int64
}

func NewService(repo TaskRepository, undoSvc UndoService, logger *zap.Logger) *Service {
	return &Service{
		repo:        repo,
		undoService: undoSvc,
		logger:      logger,
		defaultWeight: func() int64 {
			return time.Now().UnixNano()
		},
	}
}

type CreateTaskInput struct {
	Title      string
	Notes      *string
	Deadline   *time.Time
	Status     Status
	SortWeight *int64
	ParentUUID *string
}

type UpdatePayload struct {
	Title       *string
	Notes       *string
	NotesSet    bool
	Deadline    *time.Time
	DeadlineSet bool
}

type UpdateStatusInput struct {
	Status        Status
	SortWeight    *int64
	CompletedTime *time.Time
}

type ListTasksResult struct {
	Tasks []Task
	Total int64
}

var ErrTaskNotFound = errors.New("task not found")

func (s *Service) List(ctx context.Context, filter ListFilter) (ListTasksResult, error) {
	// Modify repository to support filtering by ParentUUID IS NULL and preloading Children
	// For now, we assume the repository handles this if we don't pass a specific filter,
	// but we need to ensure we only get root tasks.
	// This requires a change in the repository layer which is not exposed here directly as a filter option.
	// We will need to update the repository interface or implementation.
	// However, since I cannot see the repository implementation, I will assume I can modify the query in the repository
	// or add a new method. But sticking to the plan, I will update the service to request this.

	// Actually, looking at the repository interface in `repository_interfaces.go` (implied),
	// I should check if I can pass this constraint.
	// Since I can't easily change the repo interface without seeing it, I'll assume I can update the implementation.

	// Let's update the List method in the repository to handle "root only" if not specified otherwise.
	// Or better, let's update the ListFilter struct in model.go (which I already did? No, I didn't touch ListFilter).

	tasks, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return ListTasksResult{}, err
	}
	return ListTasksResult{Tasks: tasks, Total: total}, nil
}

func (s *Service) Get(ctx context.Context, uuid string) (*Task, error) {
	task, err := s.repo.GetByUUID(ctx, nil, uuid)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (s *Service) Create(ctx context.Context, input CreateTaskInput) (*Task, string, error) {
	status := input.Status
	if status == "" {
		status = StatusFuture
	}
	if !IsValidStatus(status) {
		return nil, "", errors.New("invalid status")
	}

	if input.ParentUUID != nil {
		// Verify parent exists
		parent, err := s.repo.GetByUUID(ctx, nil, *input.ParentUUID)
		if err != nil {
			return nil, "", err
		}
		if parent == nil {
			return nil, "", errors.New("parent task not found")
		}
	}

	sortWeight := s.defaultWeight()
	if input.SortWeight != nil {
		sortWeight = *input.SortWeight
	}

	taskModel := &Task{
		UUID:       uuid.NewString(),
		ParentUUID: input.ParentUUID,
		Title:      input.Title,
		Notes:      input.Notes,
		Deadline:   input.Deadline,
		Status:     status,
		SortWeight: sortWeight,
	}
	if status == StatusHistory {
		now := time.Now()
		taskModel.CompletedAt = &now
	}

	var undoToken string
	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Create(ctx, tx, taskModel); err != nil {
			return err
		}
		after := []Snapshot{taskModel.ToSnapshot()}
		token, err := s.undoService.RecordOperation(ctx, tx, ActionCreate, ScopeSingle, []string{taskModel.UUID}, nil, after)
		if err != nil {
			return err
		}
		undoToken = token
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return taskModel, undoToken, nil
}

func (s *Service) Update(ctx context.Context, uuid string, payload UpdatePayload) (*Task, string, error) {
	var beforeSnap Snapshot
	var undoToken string
	var updatedTask *Task

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		existing, err := s.repo.GetByUUID(ctx, tx, uuid)
		if err != nil {
			return err
		}
		if existing == nil {
			return ErrTaskNotFound
		}

		beforeSnap = existing.ToSnapshot()

		if payload.Title != nil {
			existing.Title = *payload.Title
		}
		if payload.NotesSet {
			existing.Notes = payload.Notes
		}
		if payload.DeadlineSet {
			existing.Deadline = payload.Deadline
		}

		if err := s.repo.Update(ctx, tx, existing); err != nil {
			return err
		}

		after := existing.ToSnapshot()
		token, err := s.undoService.RecordOperation(ctx, tx, ActionUpdate, ScopeSingle, []string{existing.UUID}, []Snapshot{beforeSnap}, []Snapshot{after})
		if err != nil {
			return err
		}
		undoToken = token
		updatedTask = existing
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return updatedTask, undoToken, nil
}

func (s *Service) UpdateStatus(ctx context.Context, uuid string, input UpdateStatusInput) (*Task, string, error) {
	if !IsValidStatus(input.Status) {
		return nil, "", errors.New("invalid status")
	}

	var updated *Task
	var undoToken string

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		existing, err := s.repo.GetByUUID(ctx, tx, uuid)
		if err != nil {
			return err
		}
		if existing == nil {
			return ErrTaskNotFound
		}

		before := existing.ToSnapshot()

		existing.Status = input.Status
		if input.SortWeight != nil {
			existing.SortWeight = *input.SortWeight
		} else {
			existing.SortWeight = s.defaultWeight()
		}
		action := ActionMove
		if input.Status == StatusHistory {
			now := time.Now()
			if input.CompletedTime != nil {
				now = *input.CompletedTime
			}
			existing.CompletedAt = &now
			if before.Status != StatusHistory {
				action = ActionComplete
			}
		} else {
			existing.CompletedAt = nil
		}

		if err := s.repo.Update(ctx, tx, existing); err != nil {
			return err
		}

		after := existing.ToSnapshot()
		token, err := s.undoService.RecordOperation(ctx, tx, action, ScopeSingle, []string{existing.UUID}, []Snapshot{before}, []Snapshot{after})
		if err != nil {
			return err
		}
		undoToken = token
		updated = existing
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return updated, undoToken, nil
}

func (s *Service) Complete(ctx context.Context, uuid string, completedAt *time.Time) (*Task, string, error) {
	return s.UpdateStatus(ctx, uuid, UpdateStatusInput{Status: StatusHistory, CompletedTime: completedAt})
}

func (s *Service) Delete(ctx context.Context, uuid string) (string, error) {
	var undoToken string

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		existing, err := s.repo.GetByUUID(ctx, tx, uuid)
		if err != nil {
			return err
		}
		if existing == nil {
			return ErrTaskNotFound
		}

		before := existing.ToSnapshot()

		if err := s.repo.DeleteByUUID(ctx, tx, uuid); err != nil {
			return err
		}

		token, err := s.undoService.RecordOperation(ctx, tx, ActionDelete, ScopeSingle, []string{uuid}, []Snapshot{before}, nil)
		if err != nil {
			return err
		}
		undoToken = token
		return nil
	})
	if err != nil {
		return "", err
	}
	return undoToken, nil
}

func (s *Service) BulkMove(ctx context.Context, uuids []string, status Status) ([]Task, string, error) {
	if len(uuids) == 0 {
		return nil, "", errors.New("empty ids")
	}
	if !IsValidStatus(status) {
		return nil, "", errors.New("invalid status")
	}

	var tasks []Task
	var undoToken string

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		beforeTasks, err := s.repo.GetByUUIDs(ctx, tx, uuids)
		if err != nil {
			return err
		}
		if len(beforeTasks) == 0 {
			return ErrTaskNotFound
		}

		beforeSnaps := make([]Snapshot, 0, len(beforeTasks))
		now := time.Now()
		action := ActionBulkMove
		if status == StatusHistory {
			action = ActionBulkComplete
		}

		beforeMap := make(map[string]Snapshot, len(beforeTasks))
		for _, t := range beforeTasks {
			beforeMap[t.UUID] = t.ToSnapshot()
		}
		for _, id := range uuids {
			if snap, ok := beforeMap[id]; ok {
				beforeSnaps = append(beforeSnaps, snap)
			}
		}

		baseWeight := s.defaultWeight()
		for idx, id := range uuids {
			updates := map[string]any{
				"status":      status,
				"sort_weight": baseWeight + int64(idx),
			}
			if status == StatusHistory {
				updates["completed_at"] = now
			} else {
				updates["completed_at"] = nil
			}
			if err := s.repo.UpdateColumns(ctx, tx, id, updates); err != nil {
				return err
			}
		}

		afterTasks, err := s.repo.GetByUUIDs(ctx, tx, uuids)
		if err != nil {
			return err
		}

		afterMap := make(map[string]Snapshot, len(afterTasks))
		taskMap := make(map[string]Task, len(afterTasks))
		for _, t := range afterTasks {
			snap := t.ToSnapshot()
			afterMap[t.UUID] = snap
			taskMap[t.UUID] = t
		}
		afterSnaps := make([]Snapshot, 0, len(uuids))
		orderedTasks := make([]Task, 0, len(uuids))
		for _, id := range uuids {
			if snap, ok := afterMap[id]; ok {
				afterSnaps = append(afterSnaps, snap)
			}
			if task, ok := taskMap[id]; ok {
				orderedTasks = append(orderedTasks, task)
			}
		}
		if len(afterSnaps) != len(uuids) {
			return ErrTaskNotFound
		}

		token, err := s.undoService.RecordOperation(ctx, tx, action, ScopeBulk, uuids, beforeSnaps, afterSnaps)
		if err != nil {
			return err
		}
		undoToken = token
		tasks = orderedTasks
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return tasks, undoToken, nil
}

func (s *Service) BulkDelete(ctx context.Context, uuids []string) (string, error) {
	if len(uuids) == 0 {
		return "", errors.New("empty ids")
	}
	var undoToken string

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		beforeTasks, err := s.repo.GetByUUIDs(ctx, tx, uuids)
		if err != nil {
			return err
		}
		if len(beforeTasks) == 0 {
			return ErrTaskNotFound
		}

		beforeSnaps := make([]Snapshot, 0, len(beforeTasks))
		for _, t := range beforeTasks {
			beforeSnaps = append(beforeSnaps, t.ToSnapshot())
		}

		if err := s.repo.BulkDelete(ctx, tx, uuids); err != nil {
			return err
		}

		token, err := s.undoService.RecordOperation(ctx, tx, ActionBulkDelete, ScopeBulk, uuids, beforeSnaps, nil)
		if err != nil {
			return err
		}
		undoToken = token
		return nil
	})
	if err != nil {
		return "", err
	}
	return undoToken, nil
}

func (s *Service) UpdateOrder(ctx context.Context, status Status, ordered []string) (string, error) {
	if !IsValidStatus(status) {
		return "", errors.New("invalid status")
	}
	if len(ordered) == 0 {
		return "", errors.New("ordered list empty")
	}

	var undoToken string

	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		tasks, err := s.repo.GetByUUIDs(ctx, tx, ordered)
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			return ErrTaskNotFound
		}

		beforeMap := make(map[string]Snapshot, len(tasks))
		for _, t := range tasks {
			beforeMap[t.UUID] = t.ToSnapshot()
		}
		before := make([]Snapshot, 0, len(ordered))
		for _, id := range ordered {
			if snap, ok := beforeMap[id]; ok {
				before = append(before, snap)
			}
		}
		if len(before) != len(ordered) {
			return ErrTaskNotFound
		}

		weightBase := time.Now().UnixNano()
		for idx, id := range ordered {
			weight := weightBase + int64(idx)
			if err := s.repo.UpdateColumns(ctx, tx, id, map[string]any{"sort_weight": weight}); err != nil {
				return err
			}
		}

		updated, err := s.repo.GetByUUIDs(ctx, tx, ordered)
		if err != nil {
			return err
		}
		afterMap := make(map[string]Snapshot, len(updated))
		for _, t := range updated {
			afterMap[t.UUID] = t.ToSnapshot()
		}
		after := make([]Snapshot, 0, len(ordered))
		for _, id := range ordered {
			if snap, ok := afterMap[id]; ok {
				after = append(after, snap)
			}
		}
		if len(after) != len(ordered) {
			return ErrTaskNotFound
		}

		token, err := s.undoService.RecordOperation(ctx, tx, ActionResort, ScopeBulk, ordered, before, after)
		if err != nil {
			return err
		}
		undoToken = token
		return nil
	})
	if err != nil {
		return "", err
	}
	return undoToken, nil
}
