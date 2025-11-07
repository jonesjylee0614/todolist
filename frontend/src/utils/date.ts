import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import timezone from 'dayjs/plugin/timezone';
import utc from 'dayjs/plugin/utc';
import 'dayjs/locale/zh-cn';

dayjs.locale('zh-cn');
dayjs.extend(relativeTime);
dayjs.extend(utc);
dayjs.extend(timezone);

export const DEADLINE_FORMAT = 'YYYY-MM-DD';

export function formatDeadline(deadline?: string | null): string {
  if (!deadline) {
    return '无截止日期';
  }

  const today = dayjs().startOf('day');
  const date = dayjs(deadline);
  if (!date.isValid()) {
    return deadline;
  }

  if (date.isSame(today, 'day')) {
    return '今天';
  }
  if (date.isSame(today.add(1, 'day'), 'day')) {
    return '明天';
  }
  if (date.isBefore(today)) {
    return `${date.format(DEADLINE_FORMAT)}（逾期）`;
  }
  return date.format(DEADLINE_FORMAT);
}

export function isOverdue(deadline?: string | null): boolean {
  if (!deadline) {
    return false;
  }
  const date = dayjs(deadline);
  if (!date.isValid()) {
    return false;
  }
  return date.isBefore(dayjs().startOf('day'), 'day');
}

export function formatCompletedAt(completedAt?: string | null): string {
  if (!completedAt) {
    return '';
  }
  const date = dayjs(completedAt);
  if (!date.isValid()) {
    return completedAt;
  }
  return date.fromNow();
}

export function currentISO(): string {
  return dayjs().toISOString();
}

export function isToday(date: Date): boolean {
  const today = dayjs().startOf('day');
  return dayjs(date).isSame(today, 'day');
}

export function isThisWeek(date: Date): boolean {
  const today = dayjs();
  const weekStart = today.startOf('week');
  const weekEnd = today.endOf('week');
  const targetDate = dayjs(date);
  return targetDate.isAfter(weekStart) && targetDate.isBefore(weekEnd) || targetDate.isSame(weekStart, 'day') || targetDate.isSame(weekEnd, 'day');
}

export function isSameDay(date1: Date, date2: Date): boolean {
  return dayjs(date1).isSame(dayjs(date2), 'day');
}

