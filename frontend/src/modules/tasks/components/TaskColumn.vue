<template>
  <section
    class="glass-surface flex h-full flex-col overflow-hidden border border-white/10"
    :class="{ 'ring-2 ring-brand shadow-lg': highlight }"
  >
    <header class="flex items-center justify-between border-b border-white/10 px-5 py-4">
      <div class="flex items-center gap-3">
        <h3 class="text-lg font-semibold text-slate-100">{{ title }}</h3>
        <span class="rounded-full bg-white/10 px-2 py-0.5 text-xs text-slate-300">
          {{ tasks.length }}
        </span>
        <span v-if="capacityWarning" class="text-xs text-amber-300">容量建议：≤ 5</span>
      </div>
      <div class="text-xs text-slate-400">
        <slot name="extra" />
      </div>
    </header>

    <div class="min-h-[400px] flex-1 overflow-auto px-3 py-4 scrollbar-thin">
      <div v-if="loading" class="flex justify-center py-10 text-slate-400">加载中...</div>
      <div v-else-if="!tasks.length" class="flex flex-col items-center gap-3 py-12 text-center text-sm text-slate-400">
        <slot name="empty">
          <span>{{ emptyText }}</span>
        </slot>
      </div>
      <VueDraggableNext
        v-else
        v-model="internalTasks"
        :group="group"
        :animation="200"
        item-key="uuid"
        class="flex flex-col gap-3"
        ghost-class="opacity-50"
        handle="[data-draggable-handle]"
        :data-status="status"
        @end="onDragEnd"
      >
        <template #item="{ element }">
          <TaskCard
            :key="element.uuid"
            :task="element"
            @open="emit('open-task', element)"
            @complete="emit('complete', element.uuid)"
            @delete="emit('delete', element.uuid)"
          />
        </template>
      </VueDraggableNext>
    </div>
  </section>
</template>

<script setup lang="ts">
import { VueDraggableNext } from 'vue-draggable-plus';
import { computed, ref, watch } from 'vue';
import type { SortableEvent } from 'sortablejs';
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
  move: [payload: { uuid: string; from: TaskStatus; to: TaskStatus; orderedIds: string[] }];
  complete: [uuid: string];
  delete: [uuid: string];
}>();

const group = { name: 'tasks', pull: true, put: true };

const internalTasks = ref<TaskDTO[]>([...props.tasks]);

watch(
  () => props.tasks,
  (tasks) => {
    internalTasks.value = [...tasks];
  },
  { deep: true }
);

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

function onDragEnd(evt: SortableEvent) {
  const fromStatus = (evt.from?.dataset?.status ?? props.status) as TaskStatus;
  const toStatus = (evt.to?.dataset?.status ?? props.status) as TaskStatus;
  const moved = evt.newIndex != null ? internalTasks.value[evt.newIndex] : undefined;
  const orderedIds = internalTasks.value.map((task) => task.uuid);

  if (!moved) {
    emit('reorder', props.status, orderedIds);
    return;
  }

  if (fromStatus === toStatus) {
    emit('reorder', toStatus, orderedIds);
  } else {
    emit('move', {
      uuid: moved.uuid,
      from: fromStatus,
      to: toStatus,
      orderedIds
    });
  }
}
</script>

