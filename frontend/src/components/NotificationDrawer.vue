<template>
  <Teleport to="body">
    <div v-if="isOpen" class="drawer-backdrop" @click.self="closeDrawer">
      <div class="drawer-sheet" :class="{ open: isOpen }">
        <div class="drawer-header">
          <div class="drawer-title-box">
            <span class="bell-icon">🔔</span>
            <div>
              <h3 class="drawer-title">Notifikasi & Aktivitas Terbaru</h3>
              <p class="drawer-subtitle">Daftar pembaruan & permohonan masuk terbaru.</p>
            </div>
          </div>
          <button class="close-btn" @click="closeDrawer">✕</button>
        </div>

        <div class="drawer-body">
          <!-- Sort Filter Controls -->
          <div class="filter-bar" v-if="items.length > 0">
            <span class="filter-label">Urutan Waktu:</span>
            <div class="sort-toggle">
              <button 
                class="sort-btn" 
                :class="{ active: sortOrder === 'DESC' }" 
                @click="sortOrder = 'DESC'"
                title="Urutkan Waktu Terbaru (Descending)"
              >
                ⬇️ Terbaru (DESC)
              </button>
              <button 
                class="sort-btn" 
                :class="{ active: sortOrder === 'ASC' }" 
                @click="sortOrder = 'ASC'"
                title="Urutkan Waktu Terlama (Ascending)"
              >
                ⬆️ Terlama (ASC)
              </button>
            </div>
          </div>

          <div v-if="loading" class="drawer-loading">
            Memuat notifikasi...
          </div>

          <div v-else-if="sortedItems.length === 0" class="drawer-empty">
            <div class="empty-icon">🔕</div>
            <h4>Belum Ada Notifikasi Baru</h4>
            <p>Aktivitas permohonan terbaru akan muncul di sini.</p>
          </div>

          <div v-else class="notification-list">
            <div 
              v-for="item in sortedItems" 
              :key="item.id" 
              class="notification-item" 
              @click="handleItemClick(item)"
            >
              <div class="notification-top">
                <span class="proj-num">{{ item.project_number }}</span>
                <span class="proj-date">{{ formatDate(item.submitted_at || item.updated_at || item.created_at) }}</span>
              </div>
              <h4 class="proj-title">{{ item.title }}</h4>
              <p v-if="!auth.isPemohon" class="proj-meta">
                Pemohon: <strong>{{ item.user?.name || '-' }}</strong> ({{ item.company_name || 'Umum' }})
              </p>
              <div class="notification-bottom">
                <StatusBadge :status="item.status" />
                <span class="action-link">Lihat Detail &rarr;</span>
              </div>
            </div>
          </div>
        </div>

        <div class="drawer-footer">
          <button class="btn btn-secondary w-full" @click="closeDrawer">Tutup Panel Notifikasi</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import StatusBadge from './StatusBadge.vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])
const router = useRouter()
const auth = useAuthStore()

const loading = ref(false)
const items = ref([])
const sortOrder = ref('DESC') // DESC = Terbaru, ASC = Terlama

const sortedItems = computed(() => {
  return [...items.value].sort((a, b) => {
    const timeA = new Date(a.submitted_at || a.updated_at || a.created_at || 0).getTime()
    const timeB = new Date(b.submitted_at || b.updated_at || b.created_at || 0).getTime()
    if (sortOrder.value === 'ASC') {
      return timeA - timeB
    }
    return timeB - timeA
  })
})

const closeDrawer = () => {
  emit('close')
}

const fetchNotifications = async () => {
  loading.value = true
  try {
    const endpoint = auth.isPemohon ? '/projects?per_page=15' : '/reviews/projects?per_page=15'
    const res = await apiClient.get(endpoint)
    if (res.data?.data) {
      items.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to fetch notifications:', err)
  } finally {
    loading.value = false
  }
}

const handleItemClick = (item) => {
  closeDrawer()
  if (auth.isPemohon) {
    router.push(`/projects/${item.id}`)
  } else {
    router.push(`/reviews/${item.id}`)
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    fetchNotifications()
  }
})
</script>

<style scoped>
.drawer-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(6px);
  z-index: 999;
  display: flex;
  justify-content: flex-end;
}

.drawer-sheet {
  width: 420px;
  max-width: 100vw;
  height: 100vh;
  background: #0f172a;
  border-left: 1px solid var(--border-color);
  box-shadow: -10px 0 30px rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  animation: slideInRight 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideInRight {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}

.drawer-header {
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(30, 41, 59, 0.5);
}

.drawer-title-box {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.bell-icon {
  font-size: 1.5rem;
}

.drawer-title {
  font-size: 1.05rem;
  font-weight: 800;
  color: var(--text-main);
}

.drawer-subtitle {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-top: 0.1rem;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0.2rem 0.5rem;
  border-radius: 6px;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-main);
}

.drawer-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.25rem;
}

.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-color);
}

.filter-label {
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
}

.sort-toggle {
  display: flex;
  gap: 0.35rem;
  background: rgba(30, 41, 59, 0.6);
  padding: 0.2rem;
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.sort-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.3rem 0.6rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.sort-btn:hover {
  color: var(--text-main);
}

.sort-btn.active {
  background: var(--accent-primary);
  color: white;
}

.drawer-loading, .drawer-empty {
  padding: 4rem 1.5rem;
  text-align: center;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 2.5rem;
  margin-bottom: 0.75rem;
}

.drawer-empty h4 {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-main);
}

.drawer-empty p {
  font-size: 0.82rem;
  color: var(--text-muted);
  margin-top: 0.25rem;
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.notification-item {
  padding: 1rem;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.notification-item:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: var(--accent-primary);
  transform: translateX(-3px);
}

.notification-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.35rem;
}

.proj-num {
  font-family: monospace;
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--accent-primary);
}

.proj-date {
  font-size: 0.72rem;
  color: var(--text-subtle);
}

.proj-title {
  font-size: 0.92rem;
  font-weight: 700;
  color: var(--text-main);
  line-height: 1.3;
}

.proj-meta {
  font-size: 0.78rem;
  color: var(--text-muted);
  margin-top: 0.25rem;
}

.notification-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 0.75rem;
  padding-top: 0.6rem;
  border-top: 1px dashed var(--border-color);
}

.action-link {
  font-size: 0.78rem;
  color: var(--accent-primary);
  font-weight: 600;
}

.drawer-footer {
  padding: 1rem 1.25rem;
  border-top: 1px solid var(--border-color);
  background: rgba(30, 41, 59, 0.5);
}

.w-full {
  width: 100%;
}
</style>
