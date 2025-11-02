export type TaskStatus = 'now' | 'future' | 'history';

export interface TaskDTO {
  uuid: string;
  title: string;
  notes?: string | null;
  deadline?: string | null;
  status: TaskStatus;
  sortWeight: number;
  createdAt: string;
  updatedAt: string;
  completedAt?: string | null;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
  undoToken?: string | null;
}

