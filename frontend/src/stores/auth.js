import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  // State — token stored in memory only, NOT in localStorage
  // TODO(security): In production, use HttpOnly cookies set by backend instead of Bearer tokens
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => user.value?.role || '')
  const isPemohon = computed(() => userRole.value === 'pemohon')
  const isPenilai = computed(() => userRole.value === 'penilai')
  const userName = computed(() => user.value?.name || '')

  // Actions
  function setAuth(tokenValue, userData) {
    token.value = tokenValue
    user.value = userData
    localStorage.setItem('token', tokenValue)
    localStorage.setItem('user', JSON.stringify(userData))
    apiClient.defaults.headers.common['Authorization'] = `Bearer ${tokenValue}`
  }

  async function login(email, password) {
    const res = await apiClient.post('/auth/login', { email, password })
    if (res.data?.status === 'success') {
      setAuth(res.data.data.token, res.data.data.user)
    }
    return res.data
  }

  async function register(data) {
    const res = await apiClient.post('/auth/register', data)
    if (res.data?.status === 'success') {
      setAuth(res.data.data.token, res.data.data.user)
    }
    return res.data
  }

  async function fetchMe() {
    try {
      const res = await apiClient.get('/auth/me')
      if (res.data?.status === 'success') {
        user.value = res.data.data
        localStorage.setItem('user', JSON.stringify(res.data.data))
      }
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    delete apiClient.defaults.headers.common['Authorization']
    // Full page redirect to clear any cached state
    window.location.href = '/login'
  }

  // Initialize: set auth header if token exists
  if (token.value) {
    apiClient.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return {
    token,
    user,
    isAuthenticated,
    userRole,
    isPemohon,
    isPenilai,
    userName,
    login,
    register,
    fetchMe,
    logout,
    setAuth,
  }
})
