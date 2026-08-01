<template>
  <!-- Landing Page Full-Width Layout -->
  <div v-if="route.name === 'LandingPage'" class="landing-layout">
    <router-view />
  </div>

  <!-- Authenticated App Layout with Sidebar -->
  <div v-else-if="auth.isAuthenticated" class="app-layout">
    <Sidebar />
    <div class="main-area">
      <header class="topbar">
        <div class="topbar-left">
          <h2 class="page-heading">{{ currentPageTitle }}</h2>
        </div>
        <div class="topbar-right">
          <div class="status-dot" :class="{ online: isOnline }"></div>
          <span class="role-badge">{{ getRoleBadgeText(auth.userRole) }}</span>
        </div>
      </header>
      <main class="main-content">
        <router-view />
      </main>
    </div>
  </div>

  <!-- Auth Layout (Login / Register) -->
  <div v-else class="auth-layout">
    <router-view />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import Sidebar from './components/Sidebar.vue'
import apiClient from './services/api'

const auth = useAuthStore()
const route = useRoute()
const isOnline = ref(false)

const pageNameMap = {
  Dashboard: 'Dashboard',
  MasterDocumentType: 'Master Jenis Dokumen Kelayakan',
  UserManagement: 'Manajemen User & Penilai',
  ProjectList: 'Daftar Project Permohonan',
  ProjectCreate: 'Buat Pengajuan Baru',
  ProjectEdit: 'Edit Pengajuan',
  ProjectDetail: 'Detail Pengajuan',
  ReviewList: 'Penilaian Dokumen',
  ReviewDetail: 'Detail Penilaian',
  ReviewHistory: 'Riwayat Penilaian',
  Login: 'Login',
  Register: 'Register',
  LandingPage: 'Landing Page'
}

const currentPageTitle = computed(() => pageNameMap[route.name] || 'Halaman')

const getRoleBadgeText = (role) => {
  if (role === 'admin') return 'Admin'
  if (role === 'penilai') return 'Penilai'
  return 'Pemohon'
}

const checkHealth = async () => {
  try {
    const res = await apiClient.get('/health')
    isOnline.value = res.data?.status === 'healthy'
  } catch {
    isOnline.value = false
  }
}

onMounted(() => {
  if (auth.isAuthenticated) {
    checkHealth()
    setInterval(checkHealth, 30000)
  }
})
</script>

<style scoped>
.landing-layout {
  min-height: 100vh;
  width: 100%;
}

.app-layout {
  display: flex;
  min-height: 100vh;
}

.main-area {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  transition: margin-left 0.2s ease;
}

.topbar {
  height: 64px;
  background: rgba(15, 23, 42, 0.8);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  position: sticky;
  top: 0;
  z-index: 30;
}

.page-heading {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-main);
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ef4444;
}

.status-dot.online {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
}

.role-badge {
  font-size: 0.78rem;
  font-weight: 600;
  padding: 0.3rem 0.75rem;
  border-radius: 9999px;
  background: rgba(99, 102, 241, 0.15);
  color: #818cf8;
  border: 1px solid rgba(99, 102, 241, 0.3);
  text-transform: capitalize;
}

.main-content {
  flex: 1;
  padding: 1.75rem 2rem;
}

.auth-layout {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
