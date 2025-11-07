<template>
  <button
    class="floating-add-btn group fixed z-30 flex items-center justify-center rounded-full shadow-2xl transition-all duration-300 hover:scale-110"
    @click="openModal"
    aria-label="添加任务"
    title="添加任务 (快捷键: N)"
  >
    <svg
      class="h-6 w-6 text-white transition-transform duration-300 group-hover:rotate-90"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
    </svg>
  </button>
</template>

<script setup lang="ts">
import { useUiStore } from '@/stores/ui';

const uiStore = useUiStore();

function openModal() {
  uiStore.openModal('create');
}
</script>

<style scoped>
.floating-add-btn {
  bottom: 5rem; /* 移动端避开底部导航栏 */
  right: 1.5rem;
  width: 3.5rem;
  height: 3.5rem;
  background: linear-gradient(135deg, #22d3ee, #3b82f6);
}

.floating-add-btn:hover {
  background: linear-gradient(135deg, #06b6d4, #2563eb);
  box-shadow: 0 20px 40px rgba(34, 211, 238, 0.4);
}

.floating-add-btn:active {
  transform: scale(0.95);
}

@media (min-width: 1024px) {
  .floating-add-btn {
    bottom: 2rem;
    right: 2rem;
    width: 4rem;
    height: 4rem;
  }
}
</style>

