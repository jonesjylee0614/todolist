<template>
  <div
    class="fixed top-6 right-6 z-50 flex w-80 flex-col gap-3"
    role="status"
    aria-live="polite"
  >
    <TransitionGroup name="toast" tag="div" class="flex flex-col gap-3">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="glass-surface border-l-4 p-4 shadow-2xl"
        :class="intentClass(toast.intent)"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="flex-1 text-sm text-slate-100">{{ toast.message }}</div>
          <button
            v-if="toast.actionLabel && toast.undoToken"
            class="rounded-full border border-white/20 bg-white/10 px-3 py-1 text-xs text-white transition hover:bg-white/20"
            type="button"
            @click="handleUndo(toast.undoToken, toast.id)"
          >
            {{ toast.actionLabel }}
          </button>
          <button
            class="text-slate-400 transition hover:text-slate-100"
            type="button"
            @click="dismiss(toast.id)"
          >
            ×
          </button>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { computed } from 'vue';
import { useUiStore } from '@/stores/ui';
import { useTasksStore } from '@/stores/tasks';

const uiStore = useUiStore();
const tasksStore = useTasksStore();

const { toasts } = storeToRefs(uiStore);

const intentClass = computed(
  () =>
    (intent?: string) => {
      switch (intent) {
        case 'success':
          return 'border-brand-400';
        case 'error':
          return 'border-red-400';
        case 'warning':
          return 'border-amber-400';
        default:
          return 'border-slate-500/60';
      }
    }
);

function dismiss(id: string) {
  uiStore.dismissToast(id);
}

async function handleUndo(token: string, toastId: string) {
  await tasksStore.undoLast(token);
  dismiss(toastId);
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.18s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>

