<template>
  <form class="glass-surface flex flex-col gap-3 p-4" @submit.prevent="handleSubmit('future')">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-semibold text-slate-100">快速添加</h2>
      <div class="text-xs text-slate-400">
        <span class="mr-3">Enter → 未来</span>
        <span>Shift+Enter → 现在</span>
      </div>
    </div>
    <div class="flex flex-col gap-3">
      <input
        ref="titleInput"
        v-model="title"
        class="w-full rounded-xl border border-white/10 bg-white/10 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
        placeholder="任务标题"
        @keydown.enter.prevent="handleKeySubmit"
      />
      <div class="flex flex-col gap-2 text-xs text-slate-400">
        <label class="flex flex-col gap-1">
          <span>备注</span>
          <textarea
            v-model="notes"
            class="min-h-[72px] rounded-xl border border-white/10 bg-white/10 px-3 py-2 text-sm text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
            placeholder="补充描述（可选）"
          />
        </label>
        <label class="flex flex-col gap-1">
          <span>截止日期</span>
          <input
            v-model="deadline"
            type="date"
            class="rounded-xl border border-white/10 bg-white/10 px-3 py-2 text-sm text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
          />
        </label>
      </div>
    </div>
    <div class="flex items-center justify-end gap-2">
      <button
        class="rounded-full border border-white/10 px-4 py-2 text-sm text-slate-200 transition hover:border-white/30 hover:text-white"
        type="button"
        @click="reset"
      >
        清空
      </button>
      <button
        class="rounded-full bg-brand px-4 py-2 text-sm font-semibold text-white shadow-lg transition hover:bg-brand-400"
        type="submit"
      >
        添加到未来
      </button>
      <button
        class="rounded-full border border-brand/30 px-4 py-2 text-sm font-semibold text-brand-100 transition hover:border-brand hover:text-white"
        type="button"
        @click="handleSubmit('now')"
      >
        添加到现在
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useTasksStore } from '@/stores/tasks';

const tasksStore = useTasksStore();

const title = ref('');
const notes = ref('');
const deadline = ref('');
const titleInput = ref<HTMLInputElement | null>(null);

function reset() {
  title.value = '';
  notes.value = '';
  deadline.value = '';
  titleInput.value?.focus();
}

async function handleSubmit(status: 'now' | 'future') {
  if (!title.value.trim()) {
    return;
  }
  await tasksStore.createTask(
    {
      title: title.value.trim(),
      notes: notes.value.trim() ? notes.value.trim() : null,
      deadline: deadline.value || null,
      status
    },
    {
      toastMessage: status === 'now' ? '任务已加入现在' : '任务已加入未来'
    }
  );
  reset();
}

function handleKeySubmit(event: KeyboardEvent) {
  if (event.shiftKey) {
    handleSubmit('now');
  } else {
    handleSubmit('future');
  }
}
</script>

