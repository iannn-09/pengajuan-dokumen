<template>
  <div class="page-container">
    <div class="header-actions">
      <div>
        <h1 class="page-title">Riwayat Penilaian Saya</h1>
        <p class="page-subtitle">Log audit seluruh keputusan penilaian yang pernah Anda lakukan.</p>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-box">Memuat riwayat penilaian...</div>

    <!-- Empty State -->
    <div v-else-if="histories.length === 0" class="glass-card empty-card">
      <h3>Belum Ada Riwayat Penilaian</h3>
      <p>Anda belum pernah melakukan penilaian permohonan dokumen.</p>
    </div>

    <!-- History Table -->
    <div v-else class="glass-card table-card">
      <table class="data-table">
        <thead>
          <tr>
            <th>Waktu Keputusan</th>
            <th>No. Project</th>
            <th>Judul Project</th>
            <th>Pemohon / Perusahaan</th>
            <th>Perubahan Status</th>
            <th>Catatan Evaluasi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in histories" :key="item.id">
            <td class="font-mono">{{ formatDate(item.created_at) }}</td>
            <td class="font-mono text-accent">{{ item.project?.project_number }}</td>
            <td class="font-bold cursor-pointer" @click="$router.push(`/reviews/${item.project_id}`)">
              {{ item.project?.title }}
            </td>
            <td>
              <div>{{ item.project?.user?.name || '-' }}</div>
              <small class="text-muted">{{ item.project?.company_name || '-' }}</small>
            </td>
            <td>
              <div class="status-change">
                <StatusBadge :status="item.status_from" />
                <span>&rarr;</span>
                <StatusBadge :status="item.status_to" />
              </div>
            </td>
            <td>
              <div class="notes-cell">{{ item.notes || '-' }}</div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <Pagination 
        :page="meta.page" 
        :total-pages="meta.total_pages" 
        :total="meta.total" 
        @change-page="onPageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import apiClient from '../services/api'
import StatusBadge from '../components/StatusBadge.vue'
import Pagination from '../components/Pagination.vue'

const loading = ref(true)
const histories = ref([])

const meta = reactive({
  page: 1,
  per_page: 20,
  total: 0,
  total_pages: 1
})

const fetchHistory = async (page = 1) => {
  loading.value = true
  try {
    const res = await apiClient.get(`/reviews/history?page=${page}&per_page=${meta.per_page}`)
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

.font-mono { font-family: monospace; font-size: 0.82rem; }
.font-bold { font-weight: 700; }
.text-accent { color: var(--accent-primary); }
.text-muted { color: var(--text-muted); }
.cursor-pointer { cursor: pointer; }

.status-change {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.notes-cell {
  max-width: 300px;
  white-space: pre-wrap;
  font-size: 0.85rem;
  color: var(--text-muted);
}
</style>
