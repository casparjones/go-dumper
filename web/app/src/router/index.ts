import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: () => import('@/pages/Dashboard.vue')
    },
    {
      path: '/configure',
      name: 'Configure',
      component: () => import('@/pages/Targets.vue') // Reuse existing targets page
    },
    {
      path: '/backup',
      name: 'Backup',
      component: () => import('@/pages/Backup.vue')
    },
    {
      path: '/schedule-jobs',
      name: 'ScheduleJobs',
      component: () => import('@/pages/ScheduleJobs.vue')
    },
    {
      path: '/restore',
      name: 'Restore',
      component: () => import('@/pages/Restore.vue')
    },
    {
      path: '/manage',
      name: 'Manage',
      component: () => import('@/pages/Manage.vue')
    },
    {
      path: '/log',
      name: 'Log',
      component: () => import('@/pages/Log.vue')
    },
    {
      path: '/help',
      name: 'Help',
      component: () => import('@/pages/Help.vue')
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('@/pages/Settings.vue')
    },
    // Legacy routes for compatibility
    {
      path: '/targets',
      redirect: '/configure'
    },
    {
      path: '/targets/new',
      name: 'NewTarget',
      component: () => import('@/pages/TargetEdit.vue')
    },
    {
      path: '/targets/:id/edit',
      name: 'EditTarget',
      component: () => import('@/pages/TargetEdit.vue')
    },
    {
      path: '/targets/:id/backups',
      name: 'TargetBackups',
      component: () => import('@/pages/Backups.vue')
    }
  ]
})

export default router