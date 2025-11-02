import type { TaskDTO, TaskStatus } from '@/services/types';

const statusWeight: Record<TaskStatus, number> = {
  now: 0,
  future: 1,
  history: 2
};

export function sortTasks(tasks: TaskDTO[], status: TaskStatus): TaskDTO[] {
  if (status === 'history') {
    return [...tasks].sort((a, b) => {
      const aCompleted = a.completedAt ? new Date(a.completedAt).getTime() : 0;
      const bCompleted = b.completedAt ? new Date(b.completedAt).getTime() : 0;
      return bCompleted - aCompleted;
    });
  }

  return [...tasks].sort((a, b) => {
    const aDeadline = a.deadline ? new Date(a.deadline).getTime() : Infinity;
    const bDeadline = b.deadline ? new Date(b.deadline).getTime() : Infinity;

    if (a.deadline && b.deadline) {
      if (aDeadline !== bDeadline) {
        return aDeadline - bDeadline;
      }
    }
    if (a.deadline && !b.deadline) {
      return -1;
    }
    if (!a.deadline && b.deadline) {
      return 1;
    }
    if (a.sortWeight !== b.sortWeight) {
      return a.sortWeight - b.sortWeight;
    }
    return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
  });
}

export function compareStatus(a: TaskStatus, b: TaskStatus) {
  return statusWeight[a] - statusWeight[b];
}

