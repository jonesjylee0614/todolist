<template>
  <div class="glass-surface rounded-2xl border border-white/10 p-5 lg:sticky lg:top-8">
    <h2 class="mb-4 text-lg font-bold text-slate-50 tracking-wide">📊 统计概览</h2>
    <div class="flex flex-col gap-2">
      <!-- 完成率 -->
      <div class="glass-surface rounded-xl py-2.5 pl-3 pr-4 border border-emerald-600/30 bg-emerald-500/5 mb-1">
        <div class="flex items-center gap-3">
          <!-- 图标 -->
          <span class="text-lg shrink-0">🎯</span>
          
          <!-- 文字 + 进度条 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between mb-1">
              <div class="text-sm text-slate-300">完成率</div>
              <div class="text-xs text-slate-400">{{ historyCount }}/{{ totalTasksCount }}</div>
            </div>
            <!-- 进度条 -->
            <div class="h-1.5 bg-slate-700/50 rounded-full overflow-hidden">
              <div 
                class="h-full bg-gradient-to-r from-emerald-500 to-emerald-400 transition-all duration-500 rounded-full"
                :style="{ width: `${completionRate}%` }"
              ></div>
            </div>
          </div>
          
          <!-- 数字 -->
          <div class="text-2xl font-bold leading-tight text-emerald-400 tabular-nums shrink-0">
            {{ completionRate }}%
          </div>
        </div>
      </div>

      <!-- 总计 -->
      <StatCard
        icon="📊"
        label="总计"
        :value="totalCount"
        color="slate"
      />

      <!-- 现在 -->
      <StatCard
        icon="⏰"
        label="现在"
        :value="nowCount"
        color="emerald"
        :warning="capacityWarning"
        :clickable="true"
        @click="navigateTo('now')"
      />

      <!-- 未来 -->
      <StatCard
        icon="📅"
        label="未来"
        :value="futureCount"
        color="blue"
        :clickable="true"
        @click="navigateTo('future')"
      />

      <!-- 已完成 -->
      <StatCard
        icon="✅"
        label="已完成"
        :value="historyCount"
        color="violet"
        :clickable="true"
        @click="navigateTo('history')"
      />

      <!-- 逾期 -->
      <StatCard
        icon="⚠️"
        label="逾期"
        :value="overdueCount"
        color="red"
        :pulse="overdueCount > 0"
        :clickable="true"
        @click="filterOverdue"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useTasksStore } from '@/stores/tasks';
import { useUiStore } from '@/stores/ui';
import { isOverdue } from '@/utils/date';
import StatCard from './StatCard.vue';

const router = useRouter();
const tasksStore = useTasksStore();
const uiStore = useUiStore();

const nowCount = computed(() => tasksStore.nowTasks.length);
const futureCount = computed(() => tasksStore.futureTasks.length);
const historyCount = computed(() => tasksStore.historyTasks.length);

const totalTasksCount = computed(() => {
  return nowCount.value + futureCount.value + historyCount.value;
});

const totalCount = computed(() => nowCount.value + futureCount.value);

const completionRate = computed(() => {
  if (totalTasksCount.value === 0) return 0;
  return Math.round((historyCount.value / totalTasksCount.value) * 100);
});

const overdueCount = computed(() => {
  const nowOverdue = tasksStore.nowTasks.filter(task => task.deadline && isOverdue(task.deadline));
  const futureOverdue = tasksStore.futureTasks.filter(task => task.deadline && isOverdue(task.deadline));
  return nowOverdue.length + futureOverdue.length;
});

const capacityWarning = computed(() => tasksStore.capacityWarning);

function navigateTo(route: string) {
  router.push({ name: route });
}

function filterOverdue() {
  uiStore.setFilter('overdue');
}
</script>

