import http from './http';
import type { ApiResponse, TaskDTO, TaskStatus } from './types';

export interface ListTasksParams {
  status?: TaskStatus;
  keyword?: string;
  page?: number;
  pageSize?: number;
}

export interface CreateTaskPayload {
  title: string;
  notes?: string | null;
  deadline?: string | null;
  status?: TaskStatus;
  sortWeight?: number;
}

export interface UpdateTaskPayload {
  title?: string;
  notes?: string | null;
  deadline?: string | null;
}

export interface UpdateStatusPayload {
  status: TaskStatus;
  sortWeight?: number;
  completedAt?: string | null;
}

export interface OrderUpdatePayload {
  status: TaskStatus;
  orderedIds: string[];
}

export interface BulkOperationPayload {
  ids: string[];
}

const unwrapTask = (response: ApiResponse<TaskDTO>) => ({
  task: response.data,
  undoToken: response.undoToken ?? undefined
});

const unwrapList = (response: ApiResponse<{ items: TaskDTO[]; total: number }>) => ({
  tasks: response.data.items,
  total: response.data.total,
  undoToken: response.undoToken ?? undefined
});

const unwrapUndo = (response: ApiResponse<{ affectedIds: string[] }>) => ({
  affectedIds: response.data.affectedIds,
  undoToken: response.undoToken ?? undefined
});

export const taskApi = {
  async list(params: ListTasksParams = {}): Promise<{ tasks: TaskDTO[]; total: number }> {
    const { data } = await http.get<ApiResponse<{ items: TaskDTO[]; total: number }>>(
      '/tasks',
      { params }
    );
    return {
      tasks: data.data.items,
      total: data.data.total
    };
  },
  async get(uuid: string): Promise<TaskDTO> {
    const { data } = await http.get<ApiResponse<TaskDTO>>(`/tasks/${uuid}`);
    return data.data;
  },
  async create(payload: CreateTaskPayload) {
    const { data } = await http.post<ApiResponse<TaskDTO>>('/tasks', payload);
    return unwrapTask(data);
  },
  async update(uuid: string, payload: UpdateTaskPayload) {
    const { data } = await http.patch<ApiResponse<TaskDTO>>(`/tasks/${uuid}`, payload);
    return unwrapTask(data);
  },
  async updateStatus(uuid: string, payload: UpdateStatusPayload) {
    const { data } = await http.patch<ApiResponse<TaskDTO>>(`/tasks/${uuid}/status`, payload);
    return unwrapTask(data);
  },
  async complete(uuid: string, completedAt?: string | null) {
    const { data } = await http.post<ApiResponse<TaskDTO>>(`/tasks/${uuid}/complete`, {
      completedAt
    });
    return unwrapTask(data);
  },
  async remove(uuid: string) {
    const { data } = await http.delete<ApiResponse<{ uuid: string }>>(`/tasks/${uuid}`);
    return {
      uuid: data.data.uuid,
      undoToken: data.undoToken ?? undefined
    };
  },
  async bulkMove(ids: string[], status: TaskStatus) {
    const { data } = await http.post<ApiResponse<{ items: TaskDTO[]; total: number }>>(
      '/tasks/bulk/move',
      {
        ids,
        status
      }
    );
    return unwrapList(data);
  },
  async bulkComplete(ids: string[]) {
    const { data } = await http.post<ApiResponse<{ items: TaskDTO[]; total: number }>>(
      '/tasks/bulk/complete',
      {
        ids
      }
    );
    return unwrapList(data);
  },
  async bulkDelete(ids: string[]) {
    const { data } = await http.post<ApiResponse<{ deleted: string[] }>>(
      '/tasks/bulk/delete',
      {
        ids
      }
    );
    return {
      deleted: data.data.deleted,
      undoToken: data.undoToken ?? undefined
    };
  },
  async updateOrder(payload: OrderUpdatePayload) {
    const { data } = await http.post<
      ApiResponse<{ status: TaskStatus; orderedIds: string[] }>
    >('/tasks/order', payload);
    return {
      status: data.data.status,
      orderedIds: data.data.orderedIds,
      undoToken: data.undoToken ?? undefined
    };
  },
  async undo(token: string) {
    const { data } = await http.post<ApiResponse<{ affectedIds: string[] }>>(
      '/undo',
      { token }
    );
    return unwrapUndo(data);
  }
};

export type TaskApi = typeof taskApi;

