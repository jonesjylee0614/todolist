<template>
  <article class="task-card group flex flex-col gap-3 p-4">
    <div class="flex items-start gap-3">
      <button
        type="button"
        data-draggable-handle
        class="mt-1 flex h-8 w-8 items-center justify-center rounded-full border border-white/10 bg-white/10 text-slate-300 opacity-0 transition group-hover:opacity-100"
        title="拖拽移动"
      >
        ☰
      </button>
      <div class="flex-1">
        <button
          class="text-left text-base font-semibold text-slate-100 transition hover:text-white"
          type="button"
          @click="emit('open', task)"
        >
          {{ task.title }}
        </button>
        <p v-if="task.notes" class="mt-1 text-sm text-slate-300">{{ task.notes }}</p>
        <div class="mt-3 flex flex-wrap items-center gap-2 text-xs">
          <span class="inline-flex items-center gap-1 rounded-full bg-white/10 px-3 py-1 text-slate-200">
            <span class="h-2 w-2 rounded-full" :class="statusDot" />
            {{ statusText }}
          </span>
          <span
            v-if="task.deadline"
            class="inline-flex items-center gap-1 rounded-full px-3 py-1"
            :class="deadlineClass"
          >
            <span class="h-2 w-2 rounded-full" :class="deadlineDot" />
            截止：{{ formattedDeadline }}
          </span>
          <span v-if="task.completedAt" class="rounded-full bg-emerald-400/10 px-3 py-1 text-emerald-200">
            完成于 {{ formattedCompleted }}
          </span>
        </div>
      </div>
    </div>
    <footer class="flex items-center justify-between text-xs text-slate-400">
      <span>创建于 {{ createdAt }}</span>
      <div class="flex items-center gap-2">
        <button
          v-if="task.status !== 'history'"
          class="rounded-full border border-emerald-400/40 px-3 py-1 text-emerald-200 transition hover:border-emerald-400 hover:text-emerald-100"
          type="button"
          @click="emit('complete')"
        >
          完成
        </button>
        <button
          class="rounded-full border border-white/10 px-3 py-1 text-slate-300 transition hover:border-red-400 hover:text-red-200"
          type="button"
          @click="emit('delete')"
        >
          删除
        </button>
      </div>
    </footer>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { formatDeadline, formatCompletedAt, isOverdue } from '@/utils/date';
import type { TaskDTO } from '@/services/types';

const props = defineProps<{ task: TaskDTO }>();
const emit = defineEmits<{
  complete: [];
  open: [task: TaskDTO];
  delete: [];
}>();

const formattedDeadline = computed(() => formatDeadline(props.task.deadline));
const formattedCompleted = computed(() => formatCompletedAt(props.task.completedAt));

const statusText = computed(() => {
  switch (props.task.status) {
    case 'now':
      return '正在进行';
    case 'future':
      return '规划中';
    case 'history':
      return '已完成';
    default:
      return props.task.status;
  }
});

const statusDot = computed(() => {
  switch (props.task.status) {
    case 'now':
      return 'bg-emerald-400';
    case 'future':
      return 'bg-brand-300';
    case 'history':
      return 'bg-slate-400';
    default:
      return 'bg-slate-500';
  }
});

const overdue = computed(() => isOverdue(props.task.deadline));

const deadlineClass = computed(() =>
  overdue.value
    ? 'bg-red-500/10 text-red-200'
    : 'bg-brand/10 text-brand-100'
);

const deadlineDot = computed(() => (overdue.value ? 'bg-red-400' : 'bg-brand-300'));

const createdAt = computed(() => new Date(props.task.createdAt).toLocaleString());
</script>

