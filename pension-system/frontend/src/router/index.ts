import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/views/Home.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/DataSummary.vue'),
        meta: { title: '数据汇总' }
      },
      {
        path: 'data-manage',
        name: 'DataManage',
        component: () => import('@/views/DataManage.vue'),
        meta: { title: '数据管理', requiresAdmin: true }
      },
      {
        path: 'user-manage',
        name: 'UserManage',
        component: () => import('@/views/UserManage.vue'),
        meta: { title: '用户管理', requiresAdmin: true }
      },
      {
        path: 'data-form',
        name: 'DataForm',
        component: () => import('@/views/DataForm.vue'),
        meta: { title: '数据录入', requiresAdmin: true }
      },
      {
        path: 'my-data',
        name: 'MyData',
        component: () => import('@/views/MyData.vue'),
        meta: { title: '我的数据' }
      },
      {
        path: 'survey',
        name: 'Survey',
        component: () => import('@/views/Survey.vue'),
        meta: { title: '民意调查', requiresAdmin: true }
      },
      {
        path: 'issue',
        name: 'Issue',
        component: () => import('@/views/Issue.vue'),
        meta: { title: '问题收集' }
      },
      {
        path: 'account',
        name: 'Account',
        component: () => import('@/views/Account.vue'),
        meta: { title: '账户设置' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// Navigation guard
router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  const user = userStore.getUser()

  if (to.meta.requiresAuth !== false && !user) {
    next('/login')
  } else if (to.path === '/login' && user) {
    next('/dashboard')
  } else if (to.meta.requiresAdmin && !userStore.isAdmin()) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
