import type { TaskDTO } from '@/services/types';
import { isOverdue, isToday, isThisWeek } from './date';

export type FilterType = 'all' | 'with-deadline' | 'overdue' | 'today' | 'this-week' | 'no-deadline';
export type SortType = 'deadline' | 'created' | 'title';

export interface FilterOption {
  label: string;
  value: FilterType;
  icon?: string;
  filter: (task: TaskDTO) => boolean;
}

export interface SortOption {
  label: string;
  value: SortType;
  compare: (a: TaskDTO, b: TaskDTO) => number;
}

export const quickFilters: FilterOption[] = [
  {
    label: '全部',
    value: 'all',
    filter: () => true
  },
  {
    label: '有截止日期',
    value: 'with-deadline',
    icon: '📅',
    filter: (task) => !!task.deadline
  },
  {
    label: '逾期',
    value: 'overdue',
    icon: '⚠️',
    filter: (task) => task.deadline ? isOverdue(task.deadline) : false
  },
  {
    label: '今天',
    value: 'today',
    icon: '📌',
    filter: (task) => task.deadline ? isToday(new Date(task.deadline)) : false
  },
  {
    label: '本周',
    value: 'this-week',
    icon: '📆',
    filter: (task) => task.deadline ? isThisWeek(new Date(task.deadline)) : false
  },
  {
    label: '无日期',
    value: 'no-deadline',
    filter: (task) => !task.deadline
  }
];

function compareByDeadline(a: TaskDTO, b: TaskDTO): number {
  if (!a.deadline && !b.deadline) return 0;
  if (!a.deadline) return 1;
  if (!b.deadline) return -1;
  return new Date(a.deadline).getTime() - new Date(b.deadline).getTime();
}

function compareByCreated(a: TaskDTO, b: TaskDTO): number {
  return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
}

function compareByTitle(a: TaskDTO, b: TaskDTO): number {
  return a.title.localeCompare(b.title, 'zh-CN');
}

export const sortOptions: SortOption[] = [
  {
    label: '截止日期',
    value: 'deadline',
    compare: compareByDeadline
  },
  {
    label: '创建时间',
    value: 'created',
    compare: compareByCreated
  },
  {
    label: '标题',
    value: 'title',
    compare: compareByTitle
  }
];

export function applyFilter(tasks: TaskDTO[], filterType: FilterType): TaskDTO[] {
  const filterOption = quickFilters.find(f => f.value === filterType);
  if (!filterOption) {
    return tasks;
  }
  return tasks.filter(filterOption.filter);
}

export function applySort(tasks: TaskDTO[], sortType: SortType): TaskDTO[] {
  const sortOption = sortOptions.find(s => s.value === sortType);
  if (!sortOption) {
    return [...tasks];
  }
  return [...tasks].sort(sortOption.compare);
}

