<template>
  <div class="min-h-screen bg-slate-950/90 pb-24 text-slate-100">
    <div class="mx-auto flex max-w-[1600px] gap-6 px-6">
      <!-- 左侧主内容区 -->
      <div class="flex-1 min-w-0">
        <header class="flex flex-col gap-4 py-8 lg:py-10">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <h1 class="text-3xl font-bold tracking-tight lg:text-4xl">🎯 Todo Flow</h1>
              <p class="mt-1 max-w-xl text-sm text-slate-300">掌控现在、规划未来、回顾历史的轻量化任务管理器。</p>
            </div>
            <TaskSearch class="lg:max-w-xs" />
          </div>

          <!-- 筛选栏 -->
          <TaskFilter />
        </header>

        <!-- 主任务区域 -->
        <main class="pb-8">
          <!-- 导航标签页 - 和任务列表连在一起 -->
          <div class="glass-surface overflow-hidden rounded-2xl border border-white/10 shadow-lg">
            <nav class="flex border-b border-white/10 bg-white/5">
              <RouterLink
                v-for="tab in tabs"
                :key="tab.to"
                :to="tab.to"
                class="flex-1 border-b-2 px-6 py-4 text-center transition"
                :class="
                  tab.active
                    ? 'border-brand bg-white/10 font-semibold text-white'
                    : 'border-transparent text-slate-300 hover:bg-white/5 hover:text-white'
                "
              >
                {{ tab.label }}
                <span v-if="tab.count !== undefined" class="ml-2 text-xs opacity-70">({{ tab.count }})</span>
              </RouterLink>
            </nav>

            <!-- 任务列表内容 -->
            <div class="p-6">
              <RouterView v-slot="{ Component, route: childRoute }">
                <Transition name="fade-slide" mode="out-in">
                  <component :is="Component" :key="childRoute.name" />
                </Transition>
              </RouterView>
            </div>
          </div>
        </main>
      </div>
      
      <!-- 右侧统计面板 - 桌面端显示 -->
      <aside class="hidden lg:block lg:w-64 lg:flex-shrink-0 pt-8">
        <TaskStats />
      </aside>
    </div>

    <!-- 移动端统计面板 -->
    <div class="lg:hidden px-6 pb-8">
      <TaskStats />
    </div>

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
    <FloatingAddButton />
    <TaskFormModal />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue';
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router';
import { useMagicKeys } from '@vueuse/core';
import TaskStats from '@/modules/tasks/components/TaskStats.vue';
import TaskFilter from '@/modules/tasks/components/TaskFilter.vue';
import TaskSearch from '@/modules/tasks/components/TaskSearch.vue';
import ToastStack from '@/modules/tasks/components/ToastStack.vue';
import TaskDrawer from '@/modules/tasks/components/TaskDrawer.vue';
import FloatingAddButton from '@/modules/tasks/components/FloatingAddButton.vue';
import TaskFormModal from '@/modules/tasks/components/TaskFormModal.vue';
import { useTasksStore } from '@/stores/tasks';
import { useUiStore } from '@/stores/ui';

const tasksStore = useTasksStore();
const uiStore = useUiStore();
const route = useRoute();
const router = useRouter();

onMounted(() => {
  tasksStore.load();
});

const tabs = computed(() => [
  { label: '现在', to: { name: 'now' }, active: route.name === 'now', count: tasksStore.nowTasks.length },
  { label: '未来', to: { name: 'future' }, active: route.name === 'future', count: tasksStore.futureTasks.length },
  { label: '历史', to: { name: 'history' }, active: route.name === 'history', count: tasksStore.historyTasks.length }
]);

const { N, shift_N, Digit1, Digit2, Digit3 } = useMagicKeys();

// 修改 N 和 Shift+N 快捷键，现在打开模态框而不是跳转
watchKey(N, () => uiStore.openModal('create', undefined, 'future'));
watchKey(shift_N, () => uiStore.openModal('create', undefined, 'now'));
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

