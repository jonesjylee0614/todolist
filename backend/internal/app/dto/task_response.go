package dto

import (
	"time"

	domain "todolist/backend/internal/domain/task"
)

type TaskResponse struct {
	UUID        string         `json:"uuid"`
	ParentUUID  *string        `json:"parentUuid,omitempty"`
	Children    []TaskResponse `json:"children,omitempty"`
	Title       string         `json:"title"`
	Notes       *string        `json:"notes,omitempty"`
	Deadline    *string        `json:"deadline,omitempty"`
	Status      string         `json:"status"`
	SortWeight  int64          `json:"sortWeight"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
	CompletedAt *string        `json:"completedAt,omitempty"`
}

type TaskListResponse struct {
	Items []TaskResponse `json:"items"`
	Total int64          `json:"total"`
}

func FromTask(model domain.Task) TaskResponse {
	resp := TaskResponse{
		UUID:       model.UUID,
		ParentUUID: model.ParentUUID,
		Title:      model.Title,
		Notes:      model.Notes,
		Status:     string(model.Status),
		SortWeight: model.SortWeight,
		CreatedAt:  model.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  model.UpdatedAt.Format(time.RFC3339),
	}
	if model.Deadline != nil {
		formatted := model.Deadline.Format("2006-01-02")
		resp.Deadline = &formatted
	}
	if model.CompletedAt != nil {
		formatted := model.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &formatted
	}
	if len(model.Children) > 0 {
		resp.Children = FromTasks(model.Children)
	}
	return resp
}

func FromTasks(list []domain.Task) []TaskResponse {
	result := make([]TaskResponse, 0, len(list))
	for _, t := range list {
		result = append(result, FromTask(t))
	}
	return result
}
