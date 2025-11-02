package task

import "context"

// Action and Scope types to avoid importing undo domain
type Action string
type Scope string

const (
	ActionCreate       Action = "create"
	ActionUpdate       Action = "update"
	ActionMove         Action = "move"
	ActionComplete     Action = "complete"
	ActionDelete       Action = "delete"
	ActionBulkMove     Action = "bulk_move"
	ActionBulkComplete Action = "bulk_complete"
	ActionBulkDelete   Action = "bulk_delete"
	ActionResort       Action = "resort"
)

const (
	ScopeSingle Scope = "single"
	ScopeBulk   Scope = "bulk"
)

// UndoService defines the interface for undo operations
// This breaks the circular dependency between task and undo domains
type UndoService interface {
	RecordOperation(ctx context.Context, tx interface{}, action Action, scope Scope, taskIDs []string, before, after []Snapshot) (string, error)
}
