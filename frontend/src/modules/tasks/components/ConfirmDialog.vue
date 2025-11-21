<template>
  <Transition name="dialog">
    <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/60 backdrop-blur-sm transition-opacity" @click="cancel"></div>
      
      <!-- Dialog Panel -->
      <div class="relative w-full max-w-sm overflow-hidden rounded-2xl bg-slate-900 border border-white/10 shadow-2xl p-6 transform transition-all">
        <h3 class="text-lg font-semibold text-white mb-2">{{ title }}</h3>
        <p class="text-slate-300 text-sm mb-6 leading-relaxed">{{ message }}</p>
        
        <div class="flex justify-end gap-3">
          <button 
            class="px-4 py-2 rounded-lg text-sm font-medium text-slate-300 hover:bg-white/5 transition-colors"
            @click="cancel"
          >
            {{ cancelText }}
          </button>
          <button 
            class="px-4 py-2 rounded-lg text-sm font-medium text-white shadow-lg transition-all transform active:scale-95"
            :class="confirmBtnClass"
            @click="confirm"
          >
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';

const isOpen = ref(false);
const title = ref('');
const message = ref('');
const confirmText = ref('确认');
const cancelText = ref('取消');
const type = ref<'danger' | 'info'>('info');

let resolvePromise: (value: boolean) => void;

const confirmBtnClass = computed(() => {
  return type.value === 'danger' 
    ? 'bg-red-500 hover:bg-red-600 shadow-red-500/20' 
    : 'bg-brand hover:bg-brand-400 shadow-brand/20';
});

function open(options: { 
  title: string; 
  message: string; 
  confirmText?: string; 
  cancelText?: string;
  type?: 'danger' | 'info';
}) {
  title.value = options.title;
  message.value = options.message;
  confirmText.value = options.confirmText || '确认';
  cancelText.value = options.cancelText || '取消';
  type.value = options.type || 'info';
  isOpen.value = true;
  
  return new Promise<boolean>((resolve) => {
    resolvePromise = resolve;
  });
}

function confirm() {
  isOpen.value = false;
  resolvePromise?.(true);
}

function cancel() {
  isOpen.value = false;
  resolvePromise?.(false);
}



function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && isOpen.value) {
    cancel();
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeydown);
});

onUnmounted(() => {
  document.removeEventListener('keydown', onKeydown);
});

defineExpose({ open });
</script>

<style scoped>
.dialog-enter-active,
.dialog-leave-active {
  transition: all 0.2s ease-out;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-enter-from .relative,
.dialog-leave-to .relative {
  opacity: 0;
  transform: scale(0.95) translateY(10px);
}

.dialog-enter-active .relative,
.dialog-leave-active .relative {
  transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
