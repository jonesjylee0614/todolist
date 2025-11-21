<template>
  <section class="flex h-full flex-col overflow-hidden">
    <div v-if="capacityWarning" class="mb-4 rounded-lg border border-amber-400/30 bg-amber-400/10 px-4 py-2 text-sm text-amber-300">
      💡 建议"现在"列表保持在 5 个任务以内，以保持专注
    </div>

    <div class="min-h-[400px] flex-1 overflow-auto scrollbar-thin">
      <div v-if="loading" class="flex justify-center py-10 text-slate-400">加载中...</div>
      <div v-else-if="!tasks.length" class="flex flex-col items-center gap-3 py-12 text-center text-sm text-slate-400">
        <slot name="empty">
          <span>{{ emptyText }}</span>
        </slot>
      </div>
      <div v-else class="flex flex-col gap-4">
        <TaskCard
          v-for="task in internalTasks"
          :key="task.uuid"
          :task="task"
          @open="emit('open-task', task)"
          @complete="emit('complete', task.uuid)"
          @delete="emit('delete', task.uuid)"
          @postpone="emit('postpone', task.uuid)"
          @move="emit('move', task.uuid)"
          @pin="emit('pin', task.uuid)"
          @add-subtask="(t) => emit('add-subtask', t)"
          @complete-child="(t) => emit('complete', t.uuid)"
          @delete-child="(t) => emit('delete', t.uuid)"
          @postpone-child="(t) => emit('postpone', t.uuid)"
          @move-child="(t) => emit('move', t.uuid)"
        />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import type { TaskDTO, TaskStatus } from '@/services/types';
import TaskCard from './TaskCard.vue';

const props = defineProps<{
  status: TaskStatus;
  title: string;
  tasks: TaskDTO[];
  loading: boolean;
  highlight?: boolean;
  capacityWarning?: boolean;
}>();

const emit = defineEmits<{
  'open-task': [task: TaskDTO];
  reorder: [status: TaskStatus, orderedIds: string[]];
  move: [uuid: string];
  complete: [uuid: string];
  delete: [uuid: string];
  postpone: [uuid: string];
  pin: [uuid: string];
  'add-subtask': [task: TaskDTO];
}>();

const internalTasks = computed(() => props.tasks);

const emptyText = computed(() => {
  switch (props.status) {
    case 'now':
      return '现在列表为空，试着从未来拉一些任务过来';
    case 'future':
      return '未来任务为空，点击上方快速添加';
    case 'history':
      return '暂无完成记录';
    default:
      return '暂无数据';
  }
});
</script>

