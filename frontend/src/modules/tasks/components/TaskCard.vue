<template>
  <article 
    class="task-card group rounded-xl bg-white/5 p-5 transition-all hover:bg-white/10 hover:shadow-lg hover:scale-[1.01]"
    :class="{ 'opacity-60': task.status === 'history' }"
  >
    <!-- 单行布局：左侧内容 + 右侧操作 -->
    <div class="flex items-center justify-between gap-4">
      <!-- 左侧：标题、状态、日期 -->
      <div class="flex-1 min-w-0 flex items-center gap-3">
        <!-- 标题区域 -->
        <button
          class="text-left text-lg font-bold text-emerald-400 transition hover:text-emerald-300 flex-1 min-w-0 leading-snug"
          type="button"
          @click="emit('open', task)"
        >
          {{ task.title }}
        </button>
        
        <!-- 状态和日期标签 -->
        <div class="flex items-center gap-2 shrink-0">
          <!-- 状态标签 -->
          <span class="inline-flex items-center gap-1.5 rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-slate-300">
            <span class="h-2 w-2 rounded-full" :class="statusDot" />
            {{ statusText }}
          </span>
          
          <!-- 截止日期标签 -->
          <span
            v-if="task.deadline"
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium"
            :class="deadlineClass"
          >
            <span class="h-2 w-2 rounded-full" :class="deadlineDot" />
            {{ formattedDeadline }}
          </span>
          
          <!-- 完成时间（仅历史任务） -->
          <span v-if="task.status === 'history' && task.completed_at" class="text-xs text-slate-500 px-2">
            完成于 {{ formatCompletedAt(task.completed_at) }}
          </span>
        </div>
      </div>

      <!-- 右侧：操作按钮组 -->
      <div class="flex items-center gap-2 shrink-0">
        <button
          v-if="task.status !== 'history'"
          class="action-btn h-9 rounded-lg border border-emerald-500/40 bg-emerald-500/10 px-4 py-2 text-sm font-semibold text-emerald-200 transition hover:border-emerald-400 hover:bg-emerald-500/20 hover:text-emerald-100 hover:shadow-lg"
          type="button"
          @click="emit('complete')"
          title="完成任务"
        >
          ✓ 完成
        </button>
        <button
          v-if="task.status !== 'history'"
          class="action-btn h-9 rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-slate-400 transition hover:border-brand-300 hover:bg-white/10 hover:text-brand-200"
          type="button"
          @click="emit('open', task)"
          title="编辑任务"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <!-- 更多操作按钮 -->
        <div class="relative">
          <button
            ref="menuButton"
            class="action-btn h-9 rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-slate-400 transition hover:border-white/20 hover:bg-white/10 hover:text-white"
            type="button"
            @click.stop="toggleMenu"
            title="更多操作"
          >
            <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"/>
            </svg>
          </button>
          <!-- 下拉菜单 -->
          <Transition name="menu">
            <div
              v-if="showMenu"
              v-click-outside="closeMenu"
              class="absolute right-0 top-full z-20 mt-2 min-w-[160px] rounded-lg border border-white/30 bg-slate-900/95 py-1.5 shadow-2xl backdrop-blur-md"
            >
              <button
                v-if="task.status !== 'history'"
                class="menu-item flex w-full items-center gap-2.5 px-4 py-2.5 text-left text-sm text-slate-300 transition hover:bg-white/10 hover:text-white"
                @click="handlePostpone"
              >
                <span>🕓</span> 延期
              </button>
              <button
                v-if="task.status !== 'history'"
                class="menu-item flex w-full items-center gap-2.5 px-4 py-2.5 text-left text-sm text-slate-300 transition hover:bg-white/10 hover:text-white"
                @click="handleMove"
              >
                <span>🔁</span> 
                <span>{{ task.status === 'now' ? '移到未来' : '移到现在' }}</span>
              </button>
              <!-- 置顶功能需要后端支持，暂时隐藏 -->
              <!-- <button
                v-if="task.status !== 'history'"
                class="menu-item flex w-full items-center gap-2.5 px-4 py-2.5 text-left text-sm text-slate-300 transition hover:bg-white/10 hover:text-white"
                @click="handlePin"
              >
                <span>⭐</span> 置顶
              </button> -->
              <div class="my-1 h-px bg-white/10"></div>
              <button
                class="menu-item flex w-full items-center gap-2.5 px-4 py-2.5 text-left text-sm text-red-300 transition hover:bg-red-500/10 hover:text-red-200"
                @click="handleDelete"
              >
                <span>🗑</span> 删除
              </button>
            </div>
          </Transition>
        </div>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { formatDeadline, formatCompletedAt, isOverdue } from '@/utils/date';
import type { TaskDTO } from '@/services/types';

const props = defineProps<{ task: TaskDTO }>();
const emit = defineEmits<{
  complete: [];
  open: [task: TaskDTO];
  delete: [];
  postpone: [];
  move: [];
  pin: [];
}>();

const showMenu = ref(false);
const menuButton = ref<HTMLButtonElement>();

const formattedDeadline = computed(() => formatDeadline(props.task.deadline));

const statusText = computed(() => {
  switch (props.task.status) {
    case 'now':
      return '正在进行';
    case 'future':
      return '规划中';
    case 'history':
      return '已完成';
    default:
      return props.task.status;
  }
});

const statusDot = computed(() => {
  switch (props.task.status) {
    case 'now':
      return 'bg-emerald-400';
    case 'future':
      return 'bg-brand-300';
    case 'history':
      return 'bg-slate-400';
    default:
      return 'bg-slate-500';
  }
});

const overdue = computed(() => isOverdue(props.task.deadline));

const deadlineClass = computed(() =>
  overdue.value
    ? 'bg-red-500/10 text-red-200'
    : 'bg-brand/10 text-brand-100'
);

const deadlineDot = computed(() => (overdue.value ? 'bg-red-400' : 'bg-brand-300'));

function toggleMenu() {
  showMenu.value = !showMenu.value;
}

function closeMenu() {
  showMenu.value = false;
}

function handlePostpone() {
  emit('postpone');
  closeMenu();
}

function handleMove() {
  emit('move');
  closeMenu();
}

function handlePin() {
  emit('pin');
  closeMenu();
}

function handleDelete() {
  emit('delete');
  closeMenu();
}

// v-click-outside 指令
interface ClickOutsideElement extends HTMLElement {
  clickOutsideEvent?: (event: Event) => void;
}

const vClickOutside = {
  mounted(el: ClickOutsideElement, binding: { value: (event: Event) => void }) {
    el.clickOutsideEvent = (event: Event) => {
      if (!(el === event.target || el.contains(event.target as Node))) {
        binding.value(event);
      }
    };
    document.addEventListener('click', el.clickOutsideEvent);
  },
  unmounted(el: ClickOutsideElement) {
    if (el.clickOutsideEvent) {
      document.removeEventListener('click', el.clickOutsideEvent);
    }
  }
};
</script>

<style scoped>
.menu-enter-active,
.menu-leave-active {
  transition: all 0.2s ease;
}

.menu-enter-from,
.menu-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.action-btn {
  transition: all 0.2s ease;
}

.action-btn:active {
  transform: scale(0.96);
}

.task-card {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

/* 确保下拉菜单不被卡片遮挡 */
.task-card:hover {
  z-index: 1;
}
</style>

