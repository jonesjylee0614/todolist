<template>
  <div class="glass-surface flex flex-col gap-4 rounded-xl px-5 py-4 lg:flex-row lg:items-center lg:justify-between">
    <!-- 左侧：快速筛选 -->
    <div class="flex flex-wrap items-center gap-2">
      <span class="text-xs text-slate-400">快速筛选:</span>
      <button
        v-for="filter in quickFilters"
        :key="filter.value"
        class="filter-chip rounded-full border px-3 py-1.5 text-xs transition-all duration-200"
        :class="
          activeFilter === filter.value
            ? 'border-brand bg-brand/20 font-semibold text-brand-200'
            : 'border-white/10 text-slate-300 hover:border-white/30 hover:text-white'
        "
        @click="setFilter(filter.value)"
      >
        <span v-if="filter.icon" class="mr-1">{{ filter.icon }}</span>
        {{ filter.label }}
      </button>
    </div>

    <!-- 右侧：排序 + 清空 -->
    <div class="flex items-center gap-3">
      <div class="flex items-center gap-2">
        <span class="text-xs text-slate-400">排序:</span>
        <select
          v-model="activeSort"
          class="sort-select rounded-lg border border-white/10 bg-white/10 px-3 py-1.5 text-xs text-slate-200 outline-none transition focus:border-brand-300 focus:bg-white/20"
          @change="handleSortChange"
        >
          <option v-for="sort in sortOptions" :key="sort.value" :value="sort.value">
            {{ sort.label }}
          </option>
        </select>
      </div>
      <button
        v-if="hasActiveFilter"
        @click="clearFilter"
        class="clear-btn rounded-lg border border-white/10 px-3 py-1.5 text-xs text-slate-300 transition hover:border-red-400 hover:text-red-400"
      >
        清空筛选
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useUiStore } from '@/stores/ui';
import { quickFilters, sortOptions } from '@/utils/filter';

const uiStore = useUiStore();

const activeFilter = computed(() => uiStore.activeFilter);
const activeSort = computed({
  get: () => uiStore.activeSort,
  set: (value) => uiStore.setSort(value)
});

const hasActiveFilter = computed(() => activeFilter.value !== 'all');

function setFilter(filter: string) {
  uiStore.setFilter(filter);
}

function handleSortChange() {
  // Sort change is handled by v-model
}

function clearFilter() {
  uiStore.setFilter('all');
}
</script>

<style scoped>
.filter-chip {
  white-space: nowrap;
}

.sort-select option {
  background-color: rgb(15 23 42);
  color: rgb(226 232 240);
}
</style>

