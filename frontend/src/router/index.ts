import { createRouter, createWebHistory } from 'vue-router';

const routes = [
  {
    path: '/',
    redirect: { name: 'now' }
  },
  {
    path: '/now',
    name: 'now',
    component: () => import('@/modules/tasks/pages/NowView.vue')
  },
  {
    path: '/future',
    name: 'future',
    component: () => import('@/modules/tasks/pages/FutureView.vue')
  },
  {
    path: '/history',
    name: 'history',
    component: () => import('@/modules/tasks/pages/HistoryView.vue')
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;

