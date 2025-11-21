import { computed, reactive, ref } from 'vue';
import { defineStore } from 'pinia';
import { taskApi } from '@/services/taskApi';
import type { TaskDTO, TaskStatus } from '@/services/types';
import { sortTasks } from '@/utils/sort';
import { useUiStore } from './ui';

type TaskMap = Record<string, TaskDTO>;

export interface CreateTaskOptions {
  showToast?: boolean;
  toastMessage?: string;
}

const STATUSES: TaskStatus[] = ['now', 'future', 'history'];

function ensureStatus(status?: TaskStatus): TaskStatus {
  if (!status) {
    return 'future';
  }
  if (!STATUSES.includes(status)) {
    return 'future';
  }
  return status;
}

function upsertTask(tasks: TaskMap, task: TaskDTO) {
  tasks[task.uuid] = task;
  if (task.children?.length) {
    task.children.forEach((child) => upsertTask(tasks, child));
  }
}

function removeTask(tasks: TaskMap, uuid: string) {
  const task = tasks[uuid];
  if (task?.children?.length) {
    task.children.forEach((child) => removeTask(tasks, child.uuid));
  }
  delete tasks[uuid];
}

function removeFromList(list: string[], uuid: string) {
  const index = list.indexOf(uuid);
  if (index >= 0) {
    list.splice(index, 1);
  }
}

export const useTasksStore = defineStore('tasks', () => {
  const tasksById = reactive<TaskMap>({});
  const listIds = reactive<Record<TaskStatus, string[]>>({
    now: [],
    future: [],
    history: []
  });

  const loading = reactive<Record<TaskStatus | 'global', boolean>>({
    now: false,
    future: false,
    history: false,
    global: false
  });

  const uiStore = useUiStore();
  const lastUndoToken = ref<string | undefined>();

  const tasksByStatus = (status: TaskStatus) =>
    computed(() => listIds[status].map((id) => tasksById[id]).filter(Boolean));

  const nowTasks = tasksByStatus('now');
  const futureTasks = tasksByStatus('future');
  const historyTasks = tasksByStatus('history');

  function notifyError(error: unknown, fallback: string) {
    const message =
      typeof error === 'object' && error && 'message' in error && typeof (error as any).message === 'string'
        ? (error as any).message
        : fallback;
    uiStore.pushToast({ message, intent: 'error', duration: 4000 });
  }

  function applyList(status: TaskStatus, tasks: TaskDTO[]) {
    const sorted = sortTasks(tasks, status);
    listIds[status] = sorted.map((task) => {
      upsertTask(tasksById, task);
      return task.uuid;
    });
  }

  async function load(status?: TaskStatus) {
    if (status) {
      loading[status] = true;
      try {
        const { tasks } = await taskApi.list({ status });
        applyList(status, tasks);
      } finally {
        loading[status] = false;
      }
      return;
    }

    loading.global = true;
    try {
      await Promise.all(
        STATUSES.map(async (st) => {
          const { tasks } = await taskApi.list({ status: st });
          applyList(st, tasks);
        })
      );
    } finally {
      loading.global = false;
    }
  }

  function insertToList(task: TaskDTO) {
    const status = ensureStatus(task.status);
    const currentList = listIds[status];
    if (!currentList.includes(task.uuid)) {
      currentList.push(task.uuid);
    }
    const sorted = sortTasks(currentList.map((id) => tasksById[id] ?? task), status);
    listIds[status] = sorted.map((item) => item.uuid);
  }

  function syncLists() {
    STATUSES.forEach((status) => {
      const sorted = sortTasks(
        listIds[status].map((id) => tasksById[id]).filter(Boolean),
        status
      );
      listIds[status] = sorted.map((task) => task.uuid);
    });
  }

  function moveBetweenLists(uuid: string, from: TaskStatus, to: TaskStatus) {
    if (from === to) {
      syncLists();
      return;
    }
    removeFromList(listIds[from], uuid);
    if (!listIds[to].includes(uuid)) {
      listIds[to].push(uuid);
    }
    syncLists();
  }

  async function createTask(payload: { title: string; notes?: string | null; deadline?: string | null; status?: TaskStatus; parentUuid?: string | null }, options: CreateTaskOptions = {}) {
    try {
      const { task, undoToken } = await taskApi.create(payload);
      upsertTask(tasksById, task);

      // Child tasks should not be added to the main list
      // They only exist in their parent's children array
      const isChildTask = !!payload.parentUuid;

      if (!isChildTask) {
        // Only add root tasks to the main list
        insertToList(task);
      }

      // If this is a child task, refresh the parent to update its children array
      if (payload.parentUuid && tasksById[payload.parentUuid]) {
        try {
          const updatedParent = await taskApi.get(payload.parentUuid);
          upsertTask(tasksById, updatedParent);
        } catch (error) {
          console.error('Failed to refresh parent task:', error);
        }
      }

      lastUndoToken.value = undoToken;
      if (options.showToast !== false) {
        uiStore.pushUndoToast(options.toastMessage ?? '任务已创建', undoToken);
      }
      return task;
    } catch (error) {
      notifyError(error, '创建任务失败');
      throw error;
    }
  }

  async function updateTask(uuid: string, payload: { title?: string; notes?: string | null; deadline?: string | null }) {
    const existing = tasksById[uuid];
    const parentUuid = existing?.parentUuid;
    const isChildTask = !!parentUuid;

    try {
      const { task, undoToken } = await taskApi.update(uuid, payload);
      upsertTask(tasksById, task);

      // Only update lists for root tasks
      if (!isChildTask) {
        insertToList(task);
      }

      // If this was a child task, refresh the parent to update its children array
      if (parentUuid && tasksById[parentUuid]) {
        try {
          const updatedParent = await taskApi.get(parentUuid);
          upsertTask(tasksById, updatedParent);
        } catch (error) {
          console.error('Failed to refresh parent task:', error);
        }
      }

      lastUndoToken.value = undoToken;
      uiStore.pushUndoToast('任务已更新', undoToken);
      return task;
    } catch (error) {
      notifyError(error, '更新任务失败');
      throw error;
    }
  }

  async function updateStatus(uuid: string, status: TaskStatus, sortWeight?: number) {
    const existing = tasksById[uuid];
    const parentUuid = existing?.parentUuid;
    const isChildTask = !!parentUuid;

    try {
      const { task, undoToken } = await taskApi.updateStatus(uuid, {
        status,
        sortWeight
      });
      upsertTask(tasksById, task);

      // Only update lists for root tasks
      if (!isChildTask) {
        if (existing) {
          moveBetweenLists(uuid, ensureStatus(existing.status), ensureStatus(task.status));
        } else {
          insertToList(task);
        }
      }

      // If this was a child task, refresh the parent to update its children array
      if (parentUuid && tasksById[parentUuid]) {
        try {
          const updatedParent = await taskApi.get(parentUuid);
          upsertTask(tasksById, updatedParent);
        } catch (error) {
          console.error('Failed to refresh parent task:', error);
        }
      }

      lastUndoToken.value = undoToken;
      uiStore.pushUndoToast('任务已移动', undoToken);
      return task;
    } catch (error) {
      notifyError(error, '移动任务失败');
      throw error;
    }
  }

  async function completeTask(uuid: string) {
    const existing = tasksById[uuid];
    const parentUuid = existing?.parentUuid;
    const isChildTask = !!parentUuid;

    try {
      const { task, undoToken } = await taskApi.complete(uuid);
      upsertTask(tasksById, task);

      // Only update lists for root tasks
      if (!isChildTask) {
        const fromStatus = existing ? ensureStatus(existing.status) : ensureStatus(task.status);
        moveBetweenLists(uuid, fromStatus, 'history');
      }

      // If this was a child task, refresh the parent to update its children array
      if (parentUuid && tasksById[parentUuid]) {
        try {
          const updatedParent = await taskApi.get(parentUuid);
          upsertTask(tasksById, updatedParent);
        } catch (error) {
          console.error('Failed to refresh parent task:', error);
        }
      }

      lastUndoToken.value = undoToken;
      uiStore.pushUndoToast('任务已完成', undoToken);
      return task;
    } catch (error) {
      notifyError(error, '完成任务失败');
      throw error;
    }
  }

  async function deleteTask(uuid: string) {
    try {
      const existing = tasksById[uuid];
      const parentUuid = existing?.parentUuid;
      const isChildTask = !!parentUuid;

      const result = await taskApi.remove(uuid);

      if (existing) {
        removeTask(tasksById, uuid);

        // Only remove from lists for root tasks
        if (!isChildTask) {
          removeFromList(listIds[ensureStatus(existing.status)], uuid);
        }
      }

      // If this was a child task, refresh the parent to update its children array
      if (parentUuid && tasksById[parentUuid]) {
        try {
          const updatedParent = await taskApi.get(parentUuid);
          upsertTask(tasksById, updatedParent);
        } catch (error) {
          console.error('Failed to refresh parent task:', error);
        }
      }

      lastUndoToken.value = result.undoToken;
      uiStore.pushUndoToast('任务已删除', result.undoToken);
      return result.uuid;
    } catch (error) {
      notifyError(error, '删除任务失败');
      throw error;
    }
  }

  async function bulkMove(ids: string[], status: TaskStatus) {
    if (!ids.length) {
      return;
    }
    try {
      const { tasks, undoToken } = await taskApi.bulkMove(ids, status);
      tasks.forEach((task) => {
        const previous = tasksById[task.uuid];
        upsertTask(tasksById, task);
        if (previous) {
          moveBetweenLists(task.uuid, ensureStatus(previous.status), ensureStatus(task.status));
        } else {
          insertToList(task);
        }
      });
      lastUndoToken.value = undoToken;
      uiStore.pushUndoToast('批量移动完成', undoToken);
    } catch (error) {
      notifyError(error, '批量移动失败');
      throw error;
    }
  }

  async function bulkComplete(ids: string[]) {
    if (!ids.length) {
      return;
    }
    try {
      const { tasks, undoToken } = await taskApi.bulkComplete(ids);
      tasks.forEach((task) => {
        const previous = tasksById[task.uuid];
        upsertTask(tasksById, task);
        const from = previous ? ensureStatus(previous.status) : 'now';
        moveBetweenLists(task.uuid, from, 'history');
      });
      lastUndoToken.value = undoToken;
      uiStore.pushUndoToast('批量完成任务', undoToken);
    } catch (error) {
      notifyError(error, '批量完成失败');
      throw error;
    }
  }

  async function bulkDelete(ids: string[]) {
    if (!ids.length) {
      return;
    }
    try {
      const { undoToken } = await taskApi.bulkDelete(ids);
      ids.forEach((id) => {
        const existing = tasksById[id];
        if (!existing) {
          return;
        }
        removeTask(tasksById, id);
        removeFromList(listIds[ensureStatus(existing.status)], id);
      });
      lastUndoToken.value = undoToken;
      uiStore.pushUndoToast('任务已删除', undoToken);
    } catch (error) {
      notifyError(error, '批量删除失败');
      throw error;
    }
  }

  async function updateOrder(status: TaskStatus, orderedIds: string[]) {
    const filtered = orderedIds.filter((id) => listIds[status].includes(id));
    listIds[status] = filtered;
    try {
      const { undoToken } = await taskApi.updateOrder({ status, orderedIds: filtered });
      lastUndoToken.value = undoToken;
      uiStore.pushUndoToast('排序已保存', undoToken);
    } catch (error) {
      notifyError(error, '更新排序失败');
      throw error;
    }
  }

  async function undoLast(token?: string) {
    const undoToken = token ?? lastUndoToken.value;
    if (!undoToken) {
      return;
    }
    try {
      const { affectedIds } = await taskApi.undo(undoToken);
      await refreshMany(affectedIds);
      lastUndoToken.value = undefined;
      uiStore.pushToast({ message: '撤销成功', intent: 'success', duration: 2500 });
    } catch (error) {
      notifyError(error, '撤销失败');
      throw error;
    }
  }

  async function refreshMany(ids: string[]) {
    try {
      const results = await Promise.allSettled(ids.map((id) => taskApi.get(id)));
      const missing = new Set(ids);
      results.forEach((result) => {
        if (result.status === 'fulfilled') {
          const task = result.value;
          missing.delete(task.uuid);
          upsertTask(tasksById, task);
          if (!listIds[ensureStatus(task.status)].includes(task.uuid)) {
            listIds[ensureStatus(task.status)].push(task.uuid);
          }
        }
      });
      missing.forEach((id) => {
        const existing = tasksById[id];
        if (existing) {
          removeFromList(listIds[ensureStatus(existing.status)], id);
        }
        removeTask(tasksById, id);
      });
      syncLists();
    } catch (error) {
      notifyError(error, '刷新任务失败');
      throw error;
    }
  }

  const capacityWarning = computed(() => nowTasks.value.length > 5);

  return {
    tasksById,
    listIds,
    nowTasks,
    futureTasks,
    historyTasks,
    loading,
    capacityWarning,
    load,
    createTask,
    updateTask,
    updateStatus,
    completeTask,
    deleteTask,
    bulkMove,
    bulkComplete,
    bulkDelete,
    updateOrder,
    undoLast,
    refreshMany
  };
});

