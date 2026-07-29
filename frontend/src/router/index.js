import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'LandingPage',
    component: () => import('../views/LandingPage.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects',
    name: 'ProjectList',
    component: () => import('../views/ProjectList.vue'),
    meta: { requiresAuth: true, role: 'pemohon' }
  },
  {
    path: '/projects/create',
    name: 'ProjectCreate',
    component: () => import('../views/ProjectForm.vue'),
    meta: { requiresAuth: true, role: 'pemohon' }
  },
  {
    path: '/projects/:id/edit',
    name: 'ProjectEdit',
    component: () => import('../views/ProjectForm.vue'),
    meta: { requiresAuth: true, role: 'pemohon' }
  },
  {
    path: '/projects/:id',
    name: 'ProjectDetail',
    component: () => import('../views/ProjectDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/reviews',
    name: 'ReviewList',
    component: () => import('../views/ReviewList.vue'),
    meta: { requiresAuth: true, role: 'penilai' }
  },
  {
    path: '/reviews/:id',
    name: 'ReviewDetail',
    component: () => import('../views/ReviewDetail.vue'),
    meta: { requiresAuth: true, role: 'penilai' }
  },
  {
    path: '/reviews/history',
    name: 'ReviewHistory',
    component: () => import('../views/ReviewHistory.vue'),
    meta: { requiresAuth: true, role: 'penilai' }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const auth = useAuthStore()

  // Protected routes require authentication
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next('/login')
  }

  // Guest-only routes (login/register) — redirect to dashboard if already logged in
  if (to.meta.guest && auth.isAuthenticated) {
    return next('/dashboard')
  }

  // Role-based guard
  if (to.meta.role && auth.userRole !== to.meta.role) {
    return next('/dashboard')
  }

  next()
})

export default router
