import { defineStore } from 'pinia';
import { reactive } from 'vue';
import { nanoid } from 'nanoid/non-secure';

export type ToastIntent = 'success' | 'error' | 'warning' | 'info';

export interface ToastOptions {
  id?: string;
  message: string;
  intent?: ToastIntent;
  undoToken?: string;
  actionLabel?: string;
  duration?: number;
}

export interface Toast extends ToastOptions {
  id: string;
  createdAt: number;
}

export const useUiStore = defineStore('ui', () => {
  const state = reactive({
    toasts: [] as Toast[],
    drawerTaskId: null as string | null,
    isDrawerOpen: false,
    isQuickAddOpen: false,
    searchKeyword: '',
    mailbox: [] as Array<() => void>
  });

  function pushToast(options: ToastOptions) {
    const id = options.id ?? nanoid();
    const toast: Toast = {
      id,
      intent: 'info',
      duration: 4000,
      createdAt: Date.now(),
      ...options
    };
    state.toasts.push(toast);
    if (toast.duration && toast.duration > 0 && typeof window !== 'undefined') {
      window.setTimeout(() => dismissToast(id), toast.duration);
    }
    return id;
  }

  function pushUndoToast(message: string, undoToken?: string) {
    if (!undoToken) {
      pushToast({ message, intent: 'success' });
      return;
    }
    pushToast({
      message,
      intent: 'success',
      undoToken,
      actionLabel: '撤销',
      duration: 5000
    });
  }

  function dismissToast(id: string) {
    const index = state.toasts.findIndex((toast) => toast.id === id);
    if (index >= 0) {
      state.toasts.splice(index, 1);
    }
  }

  function openDrawer(taskId: string | null) {
    state.drawerTaskId = taskId;
    state.isDrawerOpen = Boolean(taskId);
  }

  function closeDrawer() {
    state.isDrawerOpen = false;
    state.drawerTaskId = null;
  }

  function setSearch(keyword: string) {
    state.searchKeyword = keyword;
  }

  return {
    ...state,
    pushToast,
    pushUndoToast,
    dismissToast,
    openDrawer,
    closeDrawer,
    setSearch
  };
});

