<template>
  <div class="relative">
    <input
      ref="inputRef"
      v-model="keyword"
      class="w-full rounded-full border border-white/10 bg-white/10 px-5 py-3 text-sm text-slate-100 shadow-lg outline-none transition focus:border-brand-300 focus:bg-white/20"
      placeholder="搜索任务 / 快捷键 /"
      type="search"
      @keyup="handleKeyup"
    />
    <span
      class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 rounded-full border border-white/20 px-2 py-1 text-[10px] uppercase tracking-wide text-slate-300"
    >
      /
    </span>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useMagicKeys } from '@vueuse/core';
import { storeToRefs } from 'pinia';
import { useUiStore } from '@/stores/ui';

const uiStore = useUiStore();
const { searchKeyword } = storeToRefs(uiStore);

const inputRef = ref<HTMLInputElement | null>(null);
const keyword = ref(searchKeyword.value);

watch(keyword, (value) => {
  uiStore.setSearch(value.trim());
});

watch(searchKeyword, (value) => {
  if (value !== keyword.value) {
    keyword.value = value;
  }
});

const { slash } = useMagicKeys();

watch(slash, (pressed) => {
  if (pressed) {
    inputRef.value?.focus();
  }
});

function handleKeyup(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    keyword.value = '';
    inputRef.value?.blur();
  }
}
</script>

