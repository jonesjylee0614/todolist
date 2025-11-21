package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "todolist/backend/internal/domain/task"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) DB() *gorm.DB {
	return r.db
}

func (r *TaskRepository) dbWith(tx interface{}) *gorm.DB {
	if tx != nil {
		if db, ok := tx.(*gorm.DB); ok {
			return db
		}
	}
	return r.db
}

func (r *TaskRepository) Create(ctx context.Context, tx interface{}, t *domain.Task) error {
	return r.dbWith(tx).WithContext(ctx).Create(t).Error
}

func (r *TaskRepository) Update(ctx context.Context, tx interface{}, t *domain.Task) error {
	return r.dbWith(tx).WithContext(ctx).Save(t).Error
}

func (r *TaskRepository) UpdateColumns(ctx context.Context, tx interface{}, uuid string, columns map[string]any) error {
	return r.dbWith(tx).WithContext(ctx).Model(&domain.Task{}).Where("uuid = ?", uuid).Updates(columns).Error
}

func (r *TaskRepository) DeleteByUUID(ctx context.Context, tx interface{}, uuid string) error {
	return r.dbWith(tx).WithContext(ctx).Where("uuid = ?", uuid).Delete(&domain.Task{}).Error
}

func (r *TaskRepository) GetByUUID(ctx context.Context, tx interface{}, uuid string) (*domain.Task, error) {
	var t domain.Task
	err := r.dbWith(tx).WithContext(ctx).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_weight ASC")
		}).
		Where("uuid = ?", uuid).
		First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) GetByUUIDs(ctx context.Context, tx interface{}, uuids []string) ([]domain.Task, error) {
	if len(uuids) == 0 {
		return []domain.Task{}, nil
	}
	var tasks []domain.Task
	err := r.dbWith(tx).WithContext(ctx).
		Where("uuid IN ?", uuids).
		Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.Task, int64, error) {
	query := r.db.WithContext(ctx).Model(&domain.Task{})

	// Only show root tasks in the main list
	query = query.Where("parent_uuid IS NULL")

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR notes LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 200 {
		filter.PageSize = 50
	}

	offset := (filter.Page - 1) * filter.PageSize

	order := "sort_weight ASC"
	if filter.Status != nil {
		if *filter.Status == domain.StatusHistory {
			order = "completed_at DESC"
		} else {
			order = "CASE WHEN deadline IS NULL THEN 1 ELSE 0 END ASC, deadline ASC, sort_weight ASC"
		}
	}

	var tasks []domain.Task
	err := query.
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_weight ASC")
		}).
		Order(order).
		Offset(offset).
		Limit(filter.PageSize).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (r *TaskRepository) BulkUpdateStatus(ctx context.Context, tx interface{}, uuids []string, status domain.Status, columns map[string]any) error {
	q := r.dbWith(tx).WithContext(ctx).Model(&domain.Task{}).Where("uuid IN ?", uuids)
	updates := map[string]any{
		"status": status,
	}
	for k, v := range columns {
		updates[k] = v
	}
	return q.Updates(updates).Error
}

func (r *TaskRepository) BulkDelete(ctx context.Context, tx interface{}, uuids []string) error {
	return r.dbWith(tx).WithContext(ctx).Where("uuid IN ?", uuids).Delete(&domain.Task{}).Error
}

func (r *TaskRepository) ReplaceSnapshots(ctx context.Context, tx interface{}, snapshots []domain.Snapshot) error {
	if len(snapshots) == 0 {
		return nil
	}
	for _, snap := range snapshots {
		taskModel := domain.FromSnapshot(snap)
		err := r.dbWith(tx).WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "uuid"}},
				UpdateAll: true,
			}).
			Create(taskModel).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TaskRepository) DeleteBySnapshots(ctx context.Context, tx interface{}, snapshots []domain.Snapshot) error {
	if len(snapshots) == 0 {
		return nil
	}
	uuids := make([]string, 0, len(snapshots))
	for _, snap := range snapshots {
		uuids = append(uuids, snap.UUID)
	}
	return r.BulkDelete(ctx, tx, uuids)
}
