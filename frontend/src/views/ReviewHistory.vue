<template>
  <div class="page-container">
    <div class="header-actions">
      <div>
        <h1 class="page-title">{{ auth.isAdmin ? 'Seluruh Penilaian' : 'Riwayat Penilaian Saya' }}</h1>
        <p class="page-subtitle">
          {{ auth.isAdmin
            ? 'Daftar permohonan yang telah diverifikasi riwayat penilaiannya.'
            : 'Daftar permohonan yang pernah Anda verifikasi dan riwayat penilaiannya.'
          }}
        </p>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-box">Memuat riwayat penilaian...</div>

    <!-- Empty State -->
    <div v-else-if="groupedProjects.length === 0" class="glass-card empty-card">
      <h3>Belum Ada Riwayat Penilaian</h3>
      <p>Belum ada aktivitas penilaian permohonan dokumen.</p>
    </div>

    <!-- History Table (Grouped by Project) -->
    <div v-else class="glass-card table-card">
      <table class="data-table">
        <thead>
          <tr>
            <th class="text-center w-12">No</th>
            <th>Waktu Terakhir Verifikasi</th>
            <th>No. Project</th>
            <th>Judul Project</th>
            <th>Pemohon / Perusahaan</th>
            <th>Penilai Terakhir</th>
            <th>Status Terakhir</th>
            <th class="text-center">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(proj, index) in groupedProjects" :key="proj.id">
            <td class="text-center font-mono text-muted">{{ (meta.page - 1) * meta.per_page + index + 1 }}</td>
            <td class="font-mono">{{ formatDate(proj.latest_history?.created_at) }}</td>
            <td class="font-mono text-accent">{{ proj.project_number }}</td>
            <td class="font-bold cursor-pointer" @click="openLogModal(proj)">
              {{ proj.title }}
            </td>
            <td>
              <div>{{ proj.user?.name || '-' }}</div>
              <small class="text-muted">{{ proj.company_name || '-' }}</small>
            </td>
            <td>
              <div class="reviewer-info">
                <div class="font-bold text-main">{{ proj.latest_history?.reviewer?.name || 'Penilai' }}</div>
                <small class="text-subtle">{{ proj.latest_history?.reviewer?.email || '-' }}</small>
              </div>
            </td>
            <td>
              <StatusBadge :status="proj.status" />
            </td>
            <td class="text-center">
              <button class="btn btn-sm btn-secondary" @click="openLogModal(proj)" title="Lihat Riwayat Penilaian">
                Detail ({{ proj.logs.length }})
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <Pagination :page="meta.page" :total-pages="meta.total_pages" :total="meta.total" @change-page="onPageChange" />
    </div>

    <!-- Modal Riwayat Penilaian & Timeline -->
    <div v-if="showModal && selectedProject" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <div>
            <h3 class="modal-title">Riwayat Penilaian</h3>
            <span class="font-mono text-accent text-sm">{{ selectedProject.project_number }} — {{ selectedProject.title }}</span>
          </div>
          <button class="close-btn" @click="showModal = false">✕</button>
        </div>

        <div class="modal-body">
          <div class="project-summary-bar">
            <div>
              <span class="text-subtle block">Pemohon / Instansi:</span>
              <strong>{{ selectedProject.user?.name || '-' }}</strong> ({{ selectedProject.company_name || '-' }})
            </div>
            <div>
              <span class="text-subtle block">Status Terakhir:</span>
              <StatusBadge :status="selectedProject.status" />
            </div>
          </div>

          <div class="timeline-section">
            <ReviewTimeline :histories="selectedProject.logs" />
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showModal = false">Tutup</button>
          <button class="btn btn-primary" @click="$router.push(`/reviews/${selectedProject.id}`)">
            Buka Detail Project Lengkap →
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import StatusBadge from '../components/StatusBadge.vue'
import Pagination from '../components/Pagination.vue'
import ReviewTimeline from '../components/ReviewTimeline.vue'

const auth = useAuthStore()
const loading = ref(true)
const histories = ref([])
const showModal = ref(false)
const selectedProject = ref(null)

const meta = reactive({
  page: 1,
  per_page: 20,
  total: 0,
  total_pages: 1
})

const groupedProjects = computed(() => {
  const map = new Map()

  for (const item of histories.value) {
    if (!item.project_id || !item.project) continue

    if (!map.has(item.project_id)) {
      map.set(item.project_id, {
        id: item.project.id,
        project_number: item.project.project_number,
        title: item.project.title,
        company_name: item.project.company_name,
        user: item.project.user,
        status: item.project.status,
        latest_history: item,
        logs: [item]
      })
    } else {
      const existing = map.get(item.project_id)
      existing.logs.push(item)
    }
  }

  return Array.from(map.values())
})

const fetchHistory = async (page = 1) => {
  loading.value = true
  try {
    const endpoint = auth.isAdmin ? '/reviews/all-history' : '/reviews/history'
    const res = await apiClient.get(`${endpoint}?page=${page}&per_page=${meta.per_page}`)
    if (res.data?.data) {
      histories.value = res.data.data
      if (res.data.meta) {
        meta.page = res.data.meta.page
        meta.total = res.data.meta.total
        meta.total_pages = res.data.meta.total_pages
      }
    }
  } catch (err) {
    console.error('Failed to fetch review history:', err)
  } finally {
    loading.value = false
  }
}

const openLogModal = (proj) => {
  selectedProject.value = proj
  showModal.value = true
}

const onPageChange = (newPage) => {
  fetchHistory(newPage)
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchHistory()
})
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.header-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--text-main);
}

.page-subtitle {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin-top: 0.2rem;
}

.loading-box {
  padding: 3rem;
  text-align: center;
  color: var(--text-muted);
}

.empty-card {
  padding: 4rem 2rem;
  text-align: center;
  color: var(--text-muted);
}

.table-card {
  padding: 1.25rem;
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.9rem;
}

.data-table th {
  padding: 0.85rem 1rem;
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
  border-bottom: 1px solid var(--border-color);
  background: rgba(15, 23, 42, 0.4);
}

.data-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-main);
  vertical-align: middle;
}

.data-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.02);
}

.font-mono {
  font-family: monospace;
  font-size: 0.82rem;
}

.font-bold {
  font-weight: 700;
}

.text-main {
  color: var(--text-main);
}

.text-subtle {
  color: var(--text-subtle);
  display: block;
  font-size: 0.78rem;
}

.text-accent {
  color: var(--accent-primary);
}

.text-muted {
  color: var(--text-muted);
}

.text-center {
  text-align: center;
}

.cursor-pointer {
  cursor: pointer;
}

.block {
  display: block;
}

.text-sm {
  font-size: 0.82rem;
}

.w-12 {
  width: 3rem;
}

.reviewer-info {
  display: flex;
  flex-direction: column;
}

/* Modal styles */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(8px);
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.modal-card {
  background: #0f172a;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 650px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.5);
}

.modal-header {
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: var(--text-main);
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

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.project-summary-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.85rem 1rem;
  background: rgba(255, 255, 255, 0.03);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  font-size: 0.88rem;
}

.timeline-section {
  margin-top: 0.5rem;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.75rem;
}

.btn-sm {
  padding: 0.35rem 0.75rem;
  font-size: 0.8rem;
}
</style>
