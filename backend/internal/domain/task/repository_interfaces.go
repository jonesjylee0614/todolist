package task

import (
	"context"

	"gorm.io/gorm"
)

// TaskRepository defines the interface for task repository operations
// This breaks the circular dependency between task domain and repository
type TaskRepository interface {
	DB() *gorm.DB
	Create(ctx context.Context, tx interface{}, t *Task) error
	Update(ctx context.Context, tx interface{}, t *Task) error
	UpdateColumns(ctx context.Context, tx interface{}, uuid string, columns map[string]any) error
	DeleteByUUID(ctx context.Context, tx interface{}, uuid string) error
	GetByUUID(ctx context.Context, tx interface{}, uuid string) (*Task, error)
	GetByUUIDs(ctx context.Context, tx interface{}, uuids []string) ([]Task, error)
	List(ctx context.Context, filter ListFilter) ([]Task, int64, error)
	BulkUpdateStatus(ctx context.Context, tx interface{}, uuids []string, status Status, columns map[string]any) error
	BulkDelete(ctx context.Context, tx interface{}, uuids []string) error
	ReplaceSnapshots(ctx context.Context, tx interface{}, snapshots []Snapshot) error
	DeleteBySnapshots(ctx context.Context, tx interface{}, snapshots []Snapshot) error
}

// UndoRepository defines the interface for undo repository operations
type UndoRepository interface {
	Create(ctx context.Context, tx interface{}, op interface{}) error
	GetByToken(ctx context.Context, tx interface{}, token string) (interface{}, error)
	MarkConsumed(ctx context.Context, tx interface{}, token string, t interface{}) error
}
