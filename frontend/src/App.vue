<template>
  <div class="min-h-screen bg-slate-950/90 pb-24 text-slate-100">
    <header class="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-8 lg:py-12">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h1 class="text-3xl font-bold tracking-tight lg:text-4xl">Todo Flow</h1>
          <p class="mt-2 max-w-xl text-sm text-slate-300">掌控现在、规划未来、回顾历史的轻量化任务管理器。</p>
        </div>
        <TaskSearch class="lg:max-w-xs" />
      </div>
      <nav class="glass-surface flex items-center gap-3 rounded-full px-3 py-2 text-sm shadow-lg">
        <RouterLink
          v-for="tab in tabs"
          :key="tab.to"
          :to="tab.to"
          class="flex-1 rounded-full px-4 py-2 text-center transition"
          :class="
            tab.active
              ? 'bg-white/20 font-semibold text-white shadow-inner'
              : 'text-slate-300 hover:bg-white/10 hover:text-white'
          "
        >
          {{ tab.label }}
          <span v-if="tab.shortcut" class="ml-2 hidden rounded-full border border-white/20 px-2 py-0.5 text-[10px] text-slate-200 lg:inline">{{ tab.shortcut }}</span>
        </RouterLink>
      </nav>
    </header>

    <main class="mx-auto flex max-w-6xl flex-col gap-8 px-6">
      <TaskQuickAdd />
      <RouterView />
    </main>

    <nav class="fixed bottom-4 left-1/2 z-30 flex w-[90%] max-w-md -translate-x-1/2 rounded-full border border-white/10 bg-slate-900/90 px-2 py-2 shadow-xl lg:hidden">
      <RouterLink
        v-for="tab in tabs"
        :key="tab.to"
        :to="tab.to"
        class="flex-1 rounded-full px-4 py-2 text-center text-sm"
        :class="tab.active ? 'bg-white/20 text-white' : 'text-slate-300'"
      >
        {{ tab.label }}
      </RouterLink>
    </nav>

    <ToastStack />
    <TaskDrawer />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue';
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router';
import { useMagicKeys } from '@vueuse/core';
import TaskQuickAdd from '@/modules/tasks/components/TaskQuickAdd.vue';
import TaskSearch from '@/modules/tasks/components/TaskSearch.vue';
import ToastStack from '@/modules/tasks/components/ToastStack.vue';
import TaskDrawer from '@/modules/tasks/components/TaskDrawer.vue';
import { useTasksStore } from '@/stores/tasks';

const tasksStore = useTasksStore();
const route = useRoute();
const router = useRouter();

onMounted(() => {
  tasksStore.load();
});

const tabs = computed(() => [
  { label: '现在', to: { name: 'now' }, active: route.name === 'now', shortcut: '1' },
  { label: '未来', to: { name: 'future' }, active: route.name === 'future', shortcut: '2' },
  { label: '历史', to: { name: 'history' }, active: route.name === 'history', shortcut: '3' }
]);

const { N, shift_N, Digit1, Digit2, Digit3 } = useMagicKeys();

watchKey(N, () => router.push({ name: 'future' }));
watchKey(shift_N, () => router.push({ name: 'now' }));
watchKey(Digit1, () => router.push({ name: 'now' }));
watchKey(Digit2, () => router.push({ name: 'future' }));
watchKey(Digit3, () => router.push({ name: 'history' }));

function watchKey(source: any, callback: () => void) {
  watch(source, (pressed) => {
    if (pressed) {
      callback();
    }
  });
}
</script>

