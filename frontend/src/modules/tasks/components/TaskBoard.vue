<template>
  <div class="flex flex-col gap-6">
    <div class="hidden items-center justify-between gap-3 text-sm text-slate-300 lg:flex">
      <div class="flex items-center gap-2">
        <span class="inline-flex h-2 w-2 rounded-full bg-emerald-400" />
        <span>现在：{{ nowTasks.length }} 项</span>
        <span
          v-if="capacityWarning"
          class="ml-3 inline-flex items-center gap-1 rounded-full bg-amber-400/10 px-3 py-1 text-xs text-amber-300"
        >
          <span class="h-2 w-2 rounded-full bg-amber-400" /> 建议控制 5 项以内
        </span>
      </div>
      <div class="flex items-center gap-4">
        <span>{{ nowTasks.length + futureTasks.length + historyTasks.length }} 条任务</span>
        <button
          type="button"
          class="rounded-full border border-white/10 px-4 py-2 text-xs text-slate-200 transition hover:border-white/30 hover:text-white"
          @click="refreshAll"
        >
          刷新
        </button>
      </div>
    </div>

    <div class="grid gap-4 lg:grid-cols-3">
      <TaskColumn
        v-for="column in visibleColumns"
        :key="column.status"
        :status="column.status"
        :title="column.title"
        :tasks="column.tasks"
        :loading="column.loading"
        :highlight="column.status === active"
        :capacity-warning="capacityWarning && column.status === 'now'"
        @open-task="openDrawer"
        @reorder="handleReorder"
        @move="handleMove"
        @complete="handleComplete"
        @delete="handleDelete"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { storeToRefs } from 'pinia';
import TaskColumn from '@/modules/tasks/components/TaskColumn.vue';
import { useTasksStore } from '@/stores/tasks';
import { useUiStore } from '@/stores/ui';
import { useDevice } from '@/utils/device';
import type { TaskDTO, TaskStatus } from '@/services/types';

interface MovePayload {
  uuid: string;
  from: TaskStatus;
  to: TaskStatus;
  orderedIds: string[];
}

const props = defineProps<{ active: TaskStatus }>();

const tasksStore = useTasksStore();
const uiStore = useUiStore();
const { nowTasks, futureTasks, historyTasks, capacityWarning } = storeToRefs(tasksStore);
const { searchKeyword } = storeToRefs(uiStore);

const { isDesktop } = useDevice();

const active = computed(() => props.active);

const filterTask = (task: TaskDTO, keyword: string) => {
  if (!keyword) {
    return true;
  }
  const lower = keyword.toLowerCase();
  return (
    task.title.toLowerCase().includes(lower) ||
    (task.notes ?? '').toLowerCase().includes(lower)
  );
};

const nowFiltered = computed(() => nowTasks.value.filter((task) => filterTask(task, searchKeyword.value)));
const futureFiltered = computed(() => futureTasks.value.filter((task) => filterTask(task, searchKeyword.value)));
const historyFiltered = computed(() => historyTasks.value.filter((task) => filterTask(task, searchKeyword.value)));

const columnConfigs = computed(() => [
  {
    status: 'now' as TaskStatus,
    title: '现在',
    tasks: nowFiltered.value,
    loading: tasksStore.loading.now
  },
  {
    status: 'future' as TaskStatus,
    title: '未来',
    tasks: futureFiltered.value,
    loading: tasksStore.loading.future
  },
  {
    status: 'history' as TaskStatus,
    title: '历史',
    tasks: historyFiltered.value,
    loading: tasksStore.loading.history
  }
]);

const visibleColumns = computed(() => {
  if (isDesktop.value) {
    return columnConfigs.value;
  }
  return columnConfigs.value.filter((column) => column.status === active.value);
});

function openDrawer(task: TaskDTO) {
  uiStore.openDrawer(task.uuid);
}

async function handleReorder(status: TaskStatus, orderedIds: string[]) {
  await tasksStore.updateOrder(status, orderedIds);
}

async function handleMove(payload: MovePayload) {
  if (payload.from === payload.to) {
    await tasksStore.updateOrder(payload.to, payload.orderedIds);
  } else {
    await tasksStore.updateStatus(payload.uuid, payload.to);
    await tasksStore.updateOrder(payload.to, payload.orderedIds);
  }
}

async function handleComplete(uuid: string) {
  await tasksStore.completeTask(uuid);
}

async function handleDelete(uuid: string) {
  await tasksStore.deleteTask(uuid);
}

async function refreshAll() {
  await tasksStore.load();
}
</script>

