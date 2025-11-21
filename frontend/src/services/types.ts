export type TaskStatus = 'now' | 'future' | 'history';

export interface TaskDTO {
  uuid: string;
  parentUuid?: string;
  children?: TaskDTO[];
  title: string;
  notes?: string;
  deadline?: string;
  status: 'now' | 'future' | 'history';
  sortWeight: number;
  createdAt: string;
  updatedAt: string;
  completedAt?: string;
}

export interface CreateTaskPayload {
  title: string;
  notes?: string;
  deadline?: string;
  status?: 'now' | 'future' | 'history';
  sortWeight?: number;
  parentUuid?: string;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
  undoToken?: string | null;
}
