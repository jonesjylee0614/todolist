<template>
  <Transition name="drawer">
    <div v-if="isDrawerOpen" class="fixed inset-0 z-40 flex">
      <div class="hidden flex-1 bg-black/40 backdrop-blur-sm md:block" @click="close" />
      <div class="relative ml-auto flex h-full w-full max-w-md flex-col bg-slate-900/95 p-6 shadow-2xl">
        <header class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-slate-100">编辑任务</h2>
          <button class="text-slate-400 transition hover:text-white" type="button" @click="close">
            ×
          </button>
        </header>

        <div v-if="task" class="mt-6 flex-1 space-y-4 overflow-auto pr-1">
          <label class="flex flex-col gap-2 text-sm text-slate-200">
            标题
            <input
              v-model="form.title"
              class="rounded-xl border border-white/10 bg-white/10 px-4 py-3 text-base text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
              placeholder="任务标题"
            />
          </label>
          <label class="flex flex-col gap-2 text-sm text-slate-200">
            备注
            <textarea
              v-model="form.notes"
              class="min-h-[120px] rounded-xl border border-white/10 bg-white/10 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
              placeholder="补充信息"
            />
          </label>
          <div class="grid grid-cols-2 gap-4 text-sm text-slate-200">
            <label class="flex flex-col gap-2">
              截止日期
              <input
                v-model="form.deadline"
                type="date"
                class="rounded-xl border border-white/10 bg-white/10 px-3 py-2 text-sm text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
              />
            </label>
            <label class="flex flex-col gap-2">
              分组
              <select
                v-model="form.status"
                class="rounded-xl border border-white/10 bg-white/10 px-3 py-2 text-sm text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
              >
                <option value="now">现在</option>
                <option value="future">未来</option>
                <option value="history">历史</option>
              </select>
            </label>
          </div>
          <div v-if="task.completedAt" class="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-100">
            已于 {{ formattedCompleted }} 完成
          </div>
        </div>
        <div v-else class="mt-6 flex flex-1 items-center justify-center text-slate-400">
          正在加载...
        </div>

        <footer class="mt-6 flex flex-col gap-3 border-t border-white/5 pt-4 text-sm">
          <div class="flex items-center gap-3">
            <button
              class="flex-1 rounded-full bg-brand px-4 py-2 font-semibold text-white transition hover:bg-brand-400"
              type="button"
              @click="save"
              :disabled="saving"
            >
              保存
            </button>
            <button
              v-if="task && task.status !== 'history'"
              class="rounded-full border border-emerald-400/40 px-4 py-2 text-emerald-200 transition hover:border-emerald-400 hover:text-emerald-100"
              type="button"
              @click="complete"
            >
              完成
            </button>
            <button
              class="rounded-full border border-red-400/40 px-4 py-2 text-red-300 transition hover:border-red-500 hover:text-red-200"
              type="button"
              @click="remove"
            >
              删除
            </button>
          </div>
        </footer>
      </div>
      <ConfirmDialog ref="confirmDialog" />
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { reactive, watch, computed, ref } from 'vue';
import { storeToRefs } from 'pinia';
import { useUiStore } from '@/stores/ui';
import { useTasksStore } from '@/stores/tasks';
import { formatCompletedAt } from '@/utils/date';
import type { TaskStatus } from '@/services/types';
import ConfirmDialog from '@/modules/tasks/components/ConfirmDialog.vue';

const uiStore = useUiStore();
const tasksStore = useTasksStore();

const { isDrawerOpen, drawerTaskId } = storeToRefs(uiStore);
const confirmDialog = ref<InstanceType<typeof ConfirmDialog>>();

const task = computed(() => {
  const id = drawerTaskId.value;
  if (!id) {
    return undefined;
  }
  return tasksStore.tasksById[id];
});

const form = reactive({
  title: '',
  notes: '',
  deadline: '',
  status: 'future' as TaskStatus
});

const saving = ref(false);

watch(
  () => task.value,
  (current) => {
    if (!current) {
      return;
    }
    form.title = current.title;
    form.notes = current.notes ?? '';
    form.deadline = current.deadline ?? '';
    form.status = current.status;
  },
  { immediate: true }
);

const formattedCompleted = computed(() => formatCompletedAt(task.value?.completedAt));

function close() {
  uiStore.closeDrawer();
}

async function save() {
  if (!task.value) {
    return;
  }
  if (!form.title.trim()) {
    return;
  }
  saving.value = true;
  try {
    await tasksStore.updateTask(task.value.uuid, {
      title: form.title.trim(),
      notes: form.notes.trim() ? form.notes.trim() : null,
      deadline: form.deadline || null
    });
    if (form.status !== task.value.status) {
      await tasksStore.updateStatus(task.value.uuid, form.status);
    }
    close();
  } finally {
    saving.value = false;
  }
}

async function complete() {
  if (!task.value) {
    return;
  }
  await tasksStore.completeTask(task.value.uuid);
  close();
}

async function remove() {
  if (!task.value || !confirmDialog.value) {
    return;
  }
  
  const confirmed = await confirmDialog.value.open({
    title: '确认删除',
    message: '确定要删除这个任务吗？此操作无法撤销。',
    confirmText: '删除',
    type: 'danger'
  });

  if (confirmed) {
    await tasksStore.deleteTask(task.value.uuid);
    close();
  }
}
</script>

<style scoped>
.drawer-enter-active,
.drawer-leave-active {
  transition: all 0.24s ease;
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>


