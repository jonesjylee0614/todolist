<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen"
        class="modal-overlay fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
        @click.self="close"
        @keydown.esc="close"
      >
        <div class="modal-container glass-surface w-full max-w-lg rounded-2xl shadow-2xl">
          <!-- Header -->
          <header class="modal-header flex items-center justify-between border-b border-white/10 px-6 py-4">
            <h2 class="text-xl font-bold text-slate-100">
              {{ isEdit ? '编辑任务' : '添加任务' }}
            </h2>
            <button
              @click="close"
              class="close-btn flex h-8 w-8 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
              aria-label="关闭"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </header>

          <!-- Form -->
          <form @submit.prevent="handleSubmit">
            <div class="modal-body space-y-4 px-6 py-5">
              <!-- 标题 -->
              <div class="form-group">
                <label class="mb-2 block text-sm font-medium text-slate-300">任务标题 <span class="text-red-400">*</span></label>
                <input
                  v-model="form.title"
                  ref="titleInput"
                  type="text"
                  placeholder="输入任务标题"
                  class="w-full rounded-xl border border-white/10 bg-white/10 px-4 py-3 text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
                  required
                />
              </div>

              <!-- 备注 -->
              <div class="form-group">
                <label class="mb-2 block text-sm font-medium text-slate-300">备注</label>
                <textarea
                  v-model="form.notes"
                  placeholder="添加补充描述（可选）"
                  rows="4"
                  class="w-full rounded-xl border border-white/10 bg-white/10 px-4 py-3 text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
                />
              </div>

              <!-- 截止日期 -->
              <div class="form-group">
                <label class="mb-2 block text-sm font-medium text-slate-300">截止日期</label>
                <div class="flex gap-2">
                  <input
                    v-model="form.deadline"
                    type="date"
                    :min="todayString"
                    class="flex-1 rounded-xl border border-white/10 bg-white/10 px-4 py-3 text-slate-100 outline-none transition focus:border-brand-300 focus:bg-white/20"
                  />
                  <button
                    v-if="form.deadline"
                    type="button"
                    @click="form.deadline = ''"
                    class="rounded-xl border border-white/10 bg-white/10 px-3 text-slate-300 transition hover:bg-white/20 hover:text-white"
                    title="清除日期"
                  >
                    ✕
                  </button>
                </div>
                <div class="mt-2 flex gap-2 text-xs">
                  <button
                    type="button"
                    @click="setDeadlineToday"
                    class="rounded-lg border border-white/10 px-2 py-1 text-slate-300 transition hover:border-brand-300 hover:text-brand-300"
                  >
                    今天
                  </button>
                  <button
                    type="button"
                    @click="setDeadlineTomorrow"
                    class="rounded-lg border border-white/10 px-2 py-1 text-slate-300 transition hover:border-brand-300 hover:text-brand-300"
                  >
                    明天
                  </button>
                  <button
                    type="button"
                    @click="setDeadlineNextWeek"
                    class="rounded-lg border border-white/10 px-2 py-1 text-slate-300 transition hover:border-brand-300 hover:text-brand-300"
                  >
                    下周
                  </button>
                </div>
              </div>

              <!-- 添加到 -->
              <div v-if="!isEdit" class="form-group">
                <label class="mb-2 block text-sm font-medium text-slate-300">添加到</label>
                <div class="radio-group flex gap-3">
                  <button
                    type="button"
                    class="group relative flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-xl border px-4 py-3 transition-all duration-200"
                    :class="[
                      form.status === 'future' 
                        ? 'border-brand bg-brand/20 text-white shadow-[0_0_15px_rgba(74,95,255,0.15)]' 
                        : 'border-white/10 bg-white/5 text-slate-400 hover:border-brand-300/50 hover:bg-white/10 hover:text-slate-200'
                    ]"
                    @click="form.status = 'future'"
                  >
                    <div v-if="form.status === 'future'" class="absolute right-2 top-2 h-1.5 w-1.5 rounded-full bg-brand shadow-[0_0_5px_rgba(74,95,255,0.8)]"></div>
                    <span>📅</span>
                    <span class="text-sm font-medium">未来 <span class="text-xs opacity-60 font-normal">(稍后)</span></span>
                  </button>
                  
                  <button
                    type="button"
                    class="group relative flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-xl border px-4 py-3 transition-all duration-200"
                    :class="[
                      form.status === 'now' 
                        ? 'border-emerald-500 bg-emerald-500/20 text-white shadow-[0_0_15px_rgba(16,185,129,0.15)]' 
                        : 'border-white/10 bg-white/5 text-slate-400 hover:border-emerald-400/50 hover:bg-white/10 hover:text-slate-200'
                    ]"
                    @click="form.status = 'now'"
                  >
                    <div v-if="form.status === 'now'" class="absolute right-2 top-2 h-1.5 w-1.5 rounded-full bg-emerald-500 shadow-[0_0_5px_rgba(16,185,129,0.8)]"></div>
                    <span>⏰</span>
                    <span class="text-sm font-medium">现在 <span class="text-xs opacity-60 font-normal">(立即)</span></span>
                  </button>
                </div>
              </div>
            </div>

            <!-- Footer -->
            <footer class="modal-footer flex items-center justify-end gap-3 border-t border-white/10 px-6 py-4">
              <button
                type="button"
                @click="close"
                class="btn-secondary rounded-full border border-white/10 px-5 py-2 text-sm text-slate-200 transition hover:border-white/30 hover:text-white"
              >
                取消
              </button>
              <button
                type="submit"
                class="btn-primary rounded-full bg-brand px-5 py-2 text-sm font-semibold text-white shadow-lg transition hover:bg-brand-400"
              >
                {{ isEdit ? '保存' : '创建任务' }}
              </button>
            </footer>
          </form>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue';
import dayjs from 'dayjs';
import { useUiStore } from '@/stores/ui';
import { useTasksStore } from '@/stores/tasks';
import { DEADLINE_FORMAT } from '@/utils/date';

const uiStore = useUiStore();
const tasksStore = useTasksStore();

const isOpen = computed(() => uiStore.isModalOpen);
const isEdit = computed(() => uiStore.modalMode === 'edit');
const editingTaskId = computed(() => uiStore.editingTaskId);
const modalDefaultStatus = computed(() => uiStore.modalDefaultStatus);

const titleInput = ref<HTMLInputElement | null>(null);

const form = ref({
  title: '',
  notes: '',
  deadline: '',
  status: 'future' as 'now' | 'future'
});

const todayString = computed(() => dayjs().format(DEADLINE_FORMAT));

// 当模态框打开时，自动聚焦到标题输入框并加载编辑数据
watch(isOpen, async (open) => {
  if (open) {
    await nextTick();
    titleInput.value?.focus();
    
    if (isEdit.value && editingTaskId.value) {
      const task = tasksStore.tasksById[editingTaskId.value];
      if (task) {
        form.value = {
          title: task.title,
          notes: task.notes || '',
          deadline: task.deadline || '',
          status: task.status === 'now' ? 'now' : 'future'
        };
      }
    } else {
      resetForm();
      // 使用默认状态
      form.value.status = modalDefaultStatus.value;
      // 设置默认日期为当天（未来视图则为明天）
      if (modalDefaultStatus.value === 'now') {
        form.value.deadline = dayjs().format(DEADLINE_FORMAT);
      } else if (modalDefaultStatus.value === 'future') {
        form.value.deadline = dayjs().add(1, 'day').format(DEADLINE_FORMAT);
      }
    }
  }
});

// 监听 ESC 键关闭模态框
watch(isOpen, (open) => {
  if (open) {
    document.addEventListener('keydown', handleEscape);
  } else {
    document.removeEventListener('keydown', handleEscape);
  }
});

function handleEscape(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    close();
  }
}

function resetForm() {
  form.value = {
    title: '',
    notes: '',
    deadline: '',
    status: 'future'
  };
}

async function handleSubmit() {
  if (!form.value.title.trim()) {
    return;
  }

  try {
    if (isEdit.value && editingTaskId.value) {
      await tasksStore.updateTask(editingTaskId.value, {
        title: form.value.title.trim(),
        notes: form.value.notes.trim() || null,
        deadline: form.value.deadline || null
      });
    } else {
      await tasksStore.createTask(
        {
          title: form.value.title.trim(),
          notes: form.value.notes.trim() || null,
          deadline: form.value.deadline || null,
          status: form.value.status,
          parentUuid: uiStore.modalParentUuid
        },
        {
          toastMessage: form.value.status === 'now' ? '任务已加入现在' : '任务已加入未来'
        }
      );
    }
    close();
  } catch (error) {
    console.error('提交任务失败:', error);
  }
}

function close() {
  uiStore.closeModal();
}

function setDeadlineToday() {
  form.value.deadline = dayjs().format(DEADLINE_FORMAT);
}

function setDeadlineTomorrow() {
  form.value.deadline = dayjs().add(1, 'day').format(DEADLINE_FORMAT);
}

function setDeadlineNextWeek() {
  form.value.deadline = dayjs().add(7, 'day').format(DEADLINE_FORMAT);
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.3s ease;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9) translateY(-20px);
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
</style>

