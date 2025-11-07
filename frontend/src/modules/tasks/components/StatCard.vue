<template>
  <div
    class="glass-surface rounded-xl py-2.5 pl-3 pr-4 transition-all duration-300"
    :class="[
      clickable ? 'cursor-pointer hover:-translate-y-0.5 hover:shadow-md hover:bg-white/10' : '',
      pulse ? 'animate-pulse' : '',
      colorClasses
    ]"
    @click="handleClick"
  >
    <div class="flex items-center gap-3">
      <!-- 图标区域 ~15% -->
      <span class="text-lg shrink-0">{{ icon }}</span>
      
      <!-- 文字区域 ~45% -->
      <div class="text-sm text-slate-300 flex-1 min-w-0">{{ label }}</div>
      
      <!-- 数字区域 ~40% -->
      <div class="flex items-center gap-2 shrink-0">
        <div class="text-2xl font-bold leading-tight transition-all duration-500 tabular-nums" :class="textColorClass">
          {{ displayValue }}
        </div>
        <div v-if="warning" class="text-base text-yellow-400">
          ⚠️
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';

interface Props {
  icon: string;
  label: string;
  value: number;
  color?: 'slate' | 'emerald' | 'blue' | 'violet' | 'red';
  warning?: boolean;
  pulse?: boolean;
  clickable?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  color: 'slate',
  warning: false,
  pulse: false,
  clickable: false
});

const emit = defineEmits<{
  click: [];
}>();

// 动画显示的值
const displayValue = ref(props.value);

// 监听 value 变化，实现递增动画
watch(() => props.value, (newValue, oldValue) => {
  const diff = newValue - oldValue;
  if (diff === 0) return;
  
  const duration = 300; // 动画持续时间（毫秒）
  const steps = 15; // 动画步数
  const stepDuration = duration / steps;
  const stepValue = diff / steps;
  
  let currentStep = 0;
  const timer = setInterval(() => {
    currentStep++;
    displayValue.value = Math.round(oldValue + stepValue * currentStep);
    
    if (currentStep >= steps) {
      displayValue.value = newValue;
      clearInterval(timer);
    }
  }, stepDuration);
});

const colorClasses = computed(() => {
  const baseClasses: Record<string, string> = {
    slate: 'border-slate-700/50',
    emerald: 'border-emerald-600/50',
    blue: 'border-blue-600/50',
    violet: 'border-violet-600/50',
    red: 'border-red-600/50'
  };
  return baseClasses[props.color];
});

const textColorClass = computed(() => {
  const textClasses: Record<string, string> = {
    slate: 'text-slate-50',
    emerald: 'text-emerald-400',
    blue: 'text-blue-400',
    violet: 'text-violet-400',
    red: 'text-red-400'
  };
  return textClasses[props.color];
});

function handleClick() {
  if (props.clickable) {
    emit('click');
  }
}
</script>

