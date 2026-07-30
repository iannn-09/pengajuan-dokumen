<template>
  <header class="navbar-wrapper">
    <div class="navbar-container">
      <div class="brand">
        <div class="brand-icon">
          <FileText :size="24" class="icon-pulse" />
        </div>
        <div class="brand-text">
          <span class="brand-title">Pengajuan Dokumen</span>
          <span class="brand-subtitle">Portal Pelayanan Layanan Digital</span>
        </div>
      </div>

      <div class="nav-actions">
        <div class="status-indicator">
          <span class="dot" :class="{ 'connected': isOnline }"></span>
          <span class="status-text">{{ isOnline ? 'System Online (Go REST API)' : 'Connecting...' }}</span>
        </div>
        <button class="btn btn-primary" @click="$emit('open-modal')">
          <PlusCircle :size="18" />
          <span>Buat Pengajuan</span>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { FileText, PlusCircle } from 'lucide-vue-next'
import { DocumentService } from '../services/api'

defineEmits(['open-modal'])

const isOnline = ref(false)

const checkBackend = async () => {
  try {
    const res = await DocumentService.checkHealth()
    if (res.data && res.data.status === 'healthy') {
      isOnline.value = true
    }
  } catch (err) {
    isOnline.value = false
  }
}

onMounted(() => {
  checkBackend()
  setInterval(checkBackend, 15000)
})
</script>

<style scoped>
.navbar-wrapper {
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 50;
}

.navbar-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.85rem;
}

.brand-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(79, 70, 229, 0.4));
  border: 1px solid rgba(99, 102, 241, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-primary);
  box-shadow: 0 0 15px rgba(99, 102, 241, 0.2);
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-title {
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--text-main);
  letter-spacing: -0.02em;
}

.brand-subtitle {
  font-size: 0.75rem;
  color: var(--text-muted);
  font-weight: 500;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 1.25rem;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(30, 41, 59, 0.6);
  padding: 0.4rem 0.85rem;
  border-radius: 9999px;
  border: 1px solid var(--border-color);
  font-size: 0.8rem;
  color: var(--text-muted);
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #ef4444;
  transition: all 0.3s ease;
}

.dot.connected {
  background-color: #10b981;
  box-shadow: 0 0 8px #10b981;
}

@media (max-width: 640px) {
  .brand-subtitle, .status-text {
    display: none;
  }
}
</style>
