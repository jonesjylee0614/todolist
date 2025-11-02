package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// TaskOperation represents an undo operation record
type TaskOperation struct {
	ID          uint   `gorm:"primaryKey"`
	Token       string `gorm:"uniqueIndex;size:26;not null"`
	Action      string `gorm:"size:20;not null"`
	Scope       string `gorm:"size:10;not null"`
	TaskIDs     string `gorm:"type:text;not null"`
	BeforeState string `gorm:"type:longtext"`
	AfterState  string `gorm:"type:longtext"`
	ConsumedAt  *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	ExpireAt    time.Time `gorm:"not null;index"`
}

// IsConsumed checks if the operation has been consumed
func (op *TaskOperation) IsConsumed() bool {
	return op.ConsumedAt != nil
}

// IsExpired checks if the operation has expired
func (op *TaskOperation) IsExpired(now time.Time) bool {
	return now.After(op.ExpireAt)
}

type UndoRepository struct {
	db *gorm.DB
}

func NewUndoRepository(db *gorm.DB) *UndoRepository {
	return &UndoRepository{db: db}
}

func (r *UndoRepository) Create(ctx context.Context, tx *gorm.DB, op *TaskOperation) error {
	return r.dbWith(tx).WithContext(ctx).Create(op).Error
}

func (r *UndoRepository) GetByToken(ctx context.Context, tx *gorm.DB, token string) (*TaskOperation, error) {
	var op TaskOperation
	err := r.dbWith(tx).WithContext(ctx).Where("token = ?", token).First(&op).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &op, nil
}

func (r *UndoRepository) MarkConsumed(ctx context.Context, tx *gorm.DB, token string, t time.Time) error {
	return r.dbWith(tx).WithContext(ctx).
		Model(&TaskOperation{}).
		Where("token = ?", token).
		Update("consumed_at", t).
		Error
}

func (r *UndoRepository) dbWith(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
