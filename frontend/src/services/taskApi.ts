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
  parentUuid?: string | null;
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

// Generic request helper to handle unwrapping
async function request<T>(
  method: 'get' | 'post' | 'patch' | 'delete',
  url: string,
  data?: any,
  config?: any
) {
  const response = await http[method]<ApiResponse<T>>(url, data, config);
  return {
    data: response.data.data,
    undoToken: response.data.undoToken ?? undefined
  };
}

// Specialized helpers for specific response types
const unwrapList = (data: { items: TaskDTO[]; total: number }, undoToken?: string) => ({
  tasks: data.items,
  total: data.total,
  undoToken
});

export const taskApi = {
  async list(params: ListTasksParams = {}) {
    const { data } = await request<{ items: TaskDTO[]; total: number }>(
      'get',
      '/tasks',
      { params }
    );
    return unwrapList(data);
  },

  async get(uuid: string) {
    const { data } = await request<TaskDTO>('get', `/tasks/${uuid}`);
    return data;
  },

  async create(payload: CreateTaskPayload) {
    const { data, undoToken } = await request<TaskDTO>('post', '/tasks', payload);
    return { task: data, undoToken };
  },

  async update(uuid: string, payload: UpdateTaskPayload) {
    const { data, undoToken } = await request<TaskDTO>('patch', `/tasks/${uuid}`, payload);
    return { task: data, undoToken };
  },

  async updateStatus(uuid: string, payload: UpdateStatusPayload) {
    const { data, undoToken } = await request<TaskDTO>(
      'patch',
      `/tasks/${uuid}/status`,
      payload
    );
    return { task: data, undoToken };
  },

  async complete(uuid: string, completedAt?: string | null) {
    const { data, undoToken } = await request<TaskDTO>('post', `/tasks/${uuid}/complete`, {
      completedAt
    });
    return { task: data, undoToken };
  },

  async remove(uuid: string) {
    const { data, undoToken } = await request<{ uuid: string }>('delete', `/tasks/${uuid}`);
    return { uuid: data.uuid, undoToken };
  },

  async bulkMove(ids: string[], status: TaskStatus) {
    const { data, undoToken } = await request<{ items: TaskDTO[]; total: number }>(
      'post',
      '/tasks/bulk/move',
      { ids, status }
    );
    return unwrapList(data, undoToken);
  },

  async bulkComplete(ids: string[]) {
    const { data, undoToken } = await request<{ items: TaskDTO[]; total: number }>(
      'post',
      '/tasks/bulk/complete',
      { ids }
    );
    return unwrapList(data, undoToken);
  },

  async bulkDelete(ids: string[]) {
    const { data, undoToken } = await request<{ deleted: string[] }>(
      'post',
      '/tasks/bulk/delete',
      { ids }
    );
    return { deleted: data.deleted, undoToken };
  },

  async updateOrder(payload: OrderUpdatePayload) {
    const { data, undoToken } = await request<{ status: TaskStatus; orderedIds: string[] }>(
      'post',
      '/tasks/order',
      payload
    );
    return { ...data, undoToken };
  },

  async undo(token: string) {
    const { data, undoToken } = await request<{ affectedIds: string[] }>(
      'post',
      '/undo',
      { token }
    );
    return { affectedIds: data.affectedIds, undoToken };
  }
};

export type TaskApi = typeof taskApi;

