package handler

import (
    "errors"
    "io"
    "time"

    "github.com/gin-gonic/gin"

    "todolist/backend/internal/app/dto"
    "todolist/backend/internal/domain/task"
    "todolist/backend/internal/pkg/response"
)

type TaskHandler struct {
    service *task.Service
}

func NewTaskHandler(service *task.Service) *TaskHandler {
    return &TaskHandler{service: service}
}

func (h *TaskHandler) List(c *gin.Context) {
    var query dto.ListQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    filter := task.ListFilter{
        Keyword:  query.Keyword,
        Page:     query.Page,
        PageSize: query.PageSize,
    }
    if query.Status != "" {
        status := task.Status(query.Status)
        if !task.IsValidStatus(status) {
            response.BadRequest(c, "invalid status")
            return
        }
        filter.Status = &status
    }

    result, err := h.service.List(c.Request.Context(), filter)
    if err != nil {
        response.InternalServerError(c, err.Error())
        return
    }

    response.Success(c, dto.TaskListResponse{Items: dto.FromTasks(result.Tasks), Total: result.Total})
}

func (h *TaskHandler) Get(c *gin.Context) {
    uuid := c.Param("uuid")
    taskModel, err := h.service.Get(c.Request.Context(), uuid)
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }

    response.Success(c, dto.FromTask(*taskModel))
}

func (h *TaskHandler) Create(c *gin.Context) {
    var req dto.CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    status := task.StatusFuture
    if req.Status != nil {
        status = task.Status(*req.Status)
        if !task.IsValidStatus(status) {
            response.BadRequest(c, "invalid status")
            return
        }
    }

    var deadline *time.Time
    if req.Deadline != nil && *req.Deadline != "" {
        parsed, err := time.Parse("2006-01-02", *req.Deadline)
        if err != nil {
            response.BadRequest(c, "invalid deadline format")
            return
        }
        deadline = &parsed
    }

    taskModel, undoToken, err := h.service.Create(c.Request.Context(), task.CreateTaskInput{
        Title:      req.Title,
        Notes:      req.Notes,
        Deadline:   deadline,
        Status:     status,
        SortWeight: req.SortWeight,
    })
    if err != nil {
        response.InternalServerError(c, err.Error())
        return
    }

    response.Created(c, dto.FromTask(*taskModel), undoToken)
}

func (h *TaskHandler) Update(c *gin.Context) {
    uuid := c.Param("uuid")
    var req dto.UpdateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    deadline := req.Deadline.Value
    payload := task.UpdatePayload{
        Title:       req.Title,
        Notes:       req.Notes.Value,
        NotesSet:    req.Notes.Set,
        Deadline:    deadline,
        DeadlineSet: req.Deadline.Set,
    }

    updated, undoToken, err := h.service.Update(c.Request.Context(), uuid, payload)
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }

    response.Success(c, dto.FromTask(*updated), undoToken)
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
    uuid := c.Param("uuid")
    var req dto.StatusUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    status := task.Status(req.Status)
    if !task.IsValidStatus(status) {
        response.BadRequest(c, "invalid status")
        return
    }
    var completedAt *time.Time
    if req.CompletedAt != nil && *req.CompletedAt != "" {
        parsed, err := time.Parse(time.RFC3339, *req.CompletedAt)
        if err != nil {
            response.BadRequest(c, "invalid completed time")
            return
        }
        completedAt = &parsed
    }

    updated, undoToken, err := h.service.UpdateStatus(c.Request.Context(), uuid, task.UpdateStatusInput{
        Status:        status,
        SortWeight:    req.SortWeight,
        CompletedTime: completedAt,
    })
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }

    response.Success(c, dto.FromTask(*updated), undoToken)
}

func (h *TaskHandler) Complete(c *gin.Context) {
    uuid := c.Param("uuid")
    var req struct {
        CompletedAt *string `json:"completedAt"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        if !errors.Is(err, io.EOF) {
            response.BadRequest(c, err.Error())
            return
        }
    }

    var completedAt *time.Time
    if req.CompletedAt != nil && *req.CompletedAt != "" {
        parsed, err := time.Parse(time.RFC3339, *req.CompletedAt)
        if err != nil {
            response.BadRequest(c, "invalid completed time")
            return
        }
        completedAt = &parsed
    }

    updated, undoToken, err := h.service.Complete(c.Request.Context(), uuid, completedAt)
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }

    response.Success(c, dto.FromTask(*updated), undoToken)
}

func (h *TaskHandler) Delete(c *gin.Context) {
    uuid := c.Param("uuid")
    undoToken, err := h.service.Delete(c.Request.Context(), uuid)
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }
    response.Success(c, gin.H{"uuid": uuid}, undoToken)
}

func (h *TaskHandler) BulkMove(c *gin.Context) {
    var req dto.BulkMoveRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    status := task.Status(req.Status)
    if !task.IsValidStatus(status) {
        response.BadRequest(c, "invalid status")
        return
    }

    tasks, undoToken, err := h.service.BulkMove(c.Request.Context(), req.IDs, status)
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }

    response.Success(c, dto.TaskListResponse{Items: dto.FromTasks(tasks), Total: int64(len(tasks))}, undoToken)
}

func (h *TaskHandler) BulkComplete(c *gin.Context) {
    var req dto.BulkOperationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    tasks, undoToken, err := h.service.BulkMove(c.Request.Context(), req.IDs, task.StatusHistory)
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }

    response.Success(c, dto.TaskListResponse{Items: dto.FromTasks(tasks), Total: int64(len(tasks))}, undoToken)
}

func (h *TaskHandler) BulkDelete(c *gin.Context) {
    var req dto.BulkOperationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    undoToken, err := h.service.BulkDelete(c.Request.Context(), req.IDs)
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }
    response.Success(c, gin.H{"deleted": req.IDs}, undoToken)
}

func (h *TaskHandler) UpdateOrder(c *gin.Context) {
    var req dto.OrderUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    status := task.Status(req.Status)
    if !task.IsValidStatus(status) {
        response.BadRequest(c, "invalid status")
        return
    }

    undoToken, err := h.service.UpdateOrder(c.Request.Context(), status, req.OrderedIDs)
    if err != nil {
        if err == task.ErrTaskNotFound {
            response.NotFound(c, "task not found")
            return
        }
        response.InternalServerError(c, err.Error())
        return
    }
    response.Success(c, gin.H{"status": status, "orderedIds": req.OrderedIDs}, undoToken)
}

