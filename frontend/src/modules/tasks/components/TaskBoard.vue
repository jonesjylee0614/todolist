<template>
  <div class="task-board-wrapper">
    <TaskColumn
      v-for="column in visibleColumns"
      :key="column.status"
      :status="column.status"
      :title="column.title"
      :tasks="column.tasks"
      :loading="column.loading"
      :highlight="column.status === props.active"
      :capacity-warning="capacityWarning && column.status === 'now'"
      @open-task="openDrawer"
      @reorder="handleReorder"
      @move="handleMoveTask"
      @complete="handleComplete"
      @delete="handleDelete"
      @postpone="handlePostpone"
      @pin="handlePin"
      @add-subtask="handleAddSubtask"
    />
    <ConfirmDialog ref="confirmDialog" />
  </div>
</template>

<script setup lang="ts">
import { computed, watch, onMounted, ref } from 'vue';
import { storeToRefs } from 'pinia';
import TaskColumn from '@/modules/tasks/components/TaskColumn.vue';
import ConfirmDialog from '@/modules/tasks/components/ConfirmDialog.vue';
import { useTasksStore } from '@/stores/tasks';
import { useUiStore } from '@/stores/ui';
import { useDevice } from '@/utils/device';
import { applyFilter, applySort } from '@/utils/filter';
import type { FilterType, SortType } from '@/utils/filter';
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
const { searchKeyword, activeFilter, activeSort } = storeToRefs(uiStore);

const { isDesktop } = useDevice();
const confirmDialog = ref<InstanceType<typeof ConfirmDialog>>();

// 组件挂载时加载当前状态的任务
onMounted(() => {
  tasksStore.load(props.active);
});

// 关键词搜索过滤
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

// 应用筛选和排序
const applyFiltersAndSort = (tasks: TaskDTO[]) => {
  // 1. 先应用关键词搜索
  let filtered = tasks.filter((task) => filterTask(task, searchKeyword.value));
  
  // 2. 应用快速筛选
  if (activeFilter.value && activeFilter.value !== 'all') {
    filtered = applyFilter(filtered, activeFilter.value as FilterType);
  }
  
  // 3. 应用排序（如果不是默认状态）
  if (activeSort.value && activeSort.value !== 'deadline') {
    filtered = applySort(filtered, activeSort.value as SortType);
  }
  
  return filtered;
};

const nowFiltered = computed(() => applyFiltersAndSort(nowTasks.value));
const futureFiltered = computed(() => applyFiltersAndSort(futureTasks.value));
const historyFiltered = computed(() => applyFiltersAndSort(historyTasks.value));

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
  // 只显示当前激活的列
  return columnConfigs.value.filter((column) => column.status === props.active);
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

async function handleMoveTask(uuid: string) {
  // 简单的状态切换逻辑：now <-> future
  const task = tasksStore.tasksById[uuid];
  if (!task) return;
  
  if (task.status === 'now') {
    await tasksStore.updateStatus(uuid, 'future');
  } else if (task.status === 'future') {
    await tasksStore.updateStatus(uuid, 'now');
  }
}

async function handleComplete(uuid: string) {
  await tasksStore.completeTask(uuid);
}

async function handleDelete(uuid: string) {
  if (!confirmDialog.value) return;
  
  const confirmed = await confirmDialog.value.open({
    title: '确认删除',
    message: '确定要删除这个任务吗？此操作无法撤销。',
    confirmText: '删除',
    type: 'danger'
  });

  if (confirmed) {
    await tasksStore.deleteTask(uuid);
  }
}

async function handlePostpone(uuid: string) {
  // 延期任务到明天
  const task = tasksStore.tasksById[uuid];
  if (task) {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const deadlineStr = tomorrow.toISOString().split('T')[0];
    await tasksStore.updateTask(uuid, { deadline: deadlineStr });
  }
}

async function handlePin(uuid: string) {
  // TODO: 实现置顶功能（需要后端支持）
  console.log('置顶任务:', uuid);
}

function handleAddSubtask(task: TaskDTO) {
  uiStore.openModal('create', undefined, undefined, task.uuid);
}

const emit = defineEmits<{
  'add-subtask': [parent: TaskDTO];
}>();
</script>


