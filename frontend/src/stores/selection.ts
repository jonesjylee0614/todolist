import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

export const useSelectionStore = defineStore('selection', () => {
  const selectedIds = ref<string[]>([]);
  const lastSelectedId = ref<string | null>(null);

  function setSelected(ids: string[]) {
    selectedIds.value = Array.from(new Set(ids));
  }

  function toggle(id: string) {
    if (selectedIds.value.includes(id)) {
      selectedIds.value = selectedIds.value.filter((item) => item !== id);
    } else {
      selectedIds.value = [...selectedIds.value, id];
      lastSelectedId.value = id;
    }
  }

  function select(id: string, append = false) {
    if (!append) {
      selectedIds.value = [id];
    } else if (!selectedIds.value.includes(id)) {
      selectedIds.value = [...selectedIds.value, id];
    }
    lastSelectedId.value = id;
  }

  function selectRange(ids: string[], anchorId: string, targetId: string) {
    const start = ids.indexOf(anchorId);
    const end = ids.indexOf(targetId);
    if (start === -1 || end === -1) {
      select(targetId);
      return;
    }
    const [from, to] = start < end ? [start, end] : [end, start];
    const range = ids.slice(from, to + 1);
    selectedIds.value = Array.from(new Set([...selectedIds.value, ...range]));
    lastSelectedId.value = targetId;
  }

  function clear() {
    selectedIds.value = [];
    lastSelectedId.value = null;
  }

  const hasSelection = computed(() => selectedIds.value.length > 0);

  return {
    selectedIds,
    lastSelectedId,
    hasSelection,
    setSelected,
    toggle,
    select,
    selectRange,
    clear
  };
});

