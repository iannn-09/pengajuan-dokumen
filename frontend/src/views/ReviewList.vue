<template>
  <div class="page-container">
    <div class="header-actions">
      <div>
        <h1 class="page-title">Penilaian Permohonan Dokumen</h1>
        <p class="page-subtitle">Daftar permohonan yang membutuhkan proses verifikasi & penilaian akhir.</p>
      </div>
    </div>

    <!-- Filter Toolbar -->
    <div class="glass-card toolbar-card">
      <div class="search-box">
        <input 
          v-model="search" 
          type="text" 
          class="form-input" 
          placeholder="Cari berdasarkan judul, nomor project, atau perusahaan..." 
          @input="debouncedFetch"
        />
      </div>

      <div class="filter-box">
        <select v-model="statusFilter" class="form-select" @change="fetchProjects">
          <option value="">Status Perlu Penilaian (Submitted & Review)</option>
          <option value="SUBMITTED">Telah Dikirim (SUBMITTED)</option>
          <option value="UNDER_REVIEW">Sedang Dinilai (UNDER_REVIEW)</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-box">Memuat data antrean penilaian...</div>

    <!-- Empty State -->
    <div v-else-if="projects.length === 0" class="glass-card empty-card">
      <h3>Antrean Penilaian Kosong</h3>
      <p>Tidak ada permohonan dokumen yang sedang menunggu verifikasi saat ini.</p>
    </div>

    <!-- Table View -->
    <div v-else class="glass-card table-card">
      <table class="data-table">
        <thead>
          <tr>
            <th class="text-center w-12">No</th>
            <th>No. Project</th>
            <th>Judul Permohonan</th>
            <th>Pemohon / Perusahaan</th>
            <th>Status</th>
            <th>Tanggal Dikirim</th>
            <th class="text-right">Aksi Penilaian</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in projects" :key="item.id">
            <td class="text-center font-mono text-muted">{{ (meta.page - 1) * meta.per_page + index + 1 }}</td>
            <td class="font-mono text-accent">{{ item.project_number }}</td>
            <td class="font-bold cursor-pointer" @click="$router.push(`/reviews/${item.id}`)">
              {{ item.title }}
            </td>
            <td>
              <div>{{ item.user?.name || '-' }}</div>
              <small class="text-muted">{{ item.company_name || '-' }}</small>
            </td>
            <td><StatusBadge :status="item.status" /></td>
            <td>{{ formatDate(item.submitted_at || item.created_at) }}</td>
            <td class="text-right">
              <button class="btn btn-primary btn-sm" @click="$router.push(`/reviews/${item.id}`)">
                Proses Penilaian &rarr;
              </button>
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
const projects = ref([])
const search = ref('')
const statusFilter = ref('')

const meta = reactive({
  page: 1,
  per_page: 20,
  total: 0,
  total_pages: 1
})

const fetchProjects = async (page = 1) => {
  loading.value = true
  try {
    const params = {
      page,
      per_page: meta.per_page
    }
    if (search.value) params.search = search.value
    if (statusFilter.value) params.status = statusFilter.value

    const res = await apiClient.get('/reviews/projects', { params })
    if (res.data?.data) {
      projects.value = res.data.data
      if (res.data.meta) {
        meta.page = res.data.meta.page
        meta.total = res.data.meta.total
        meta.total_pages = res.data.meta.total_pages
      }
    }
  } catch (err) {
    console.error('Failed to fetch reviews queue:', err)
  } finally {
    loading.value = false
  }
}

let timeout = null
const debouncedFetch = () => {
  clearTimeout(timeout)
  timeout = setTimeout(() => {
    fetchProjects(1)
  }, 300)
}

const onPageChange = (newPage) => {
  fetchProjects(newPage)
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchProjects()
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

.toolbar-card {
  padding: 1rem 1.25rem;
  display: flex;
  gap: 1rem;
}

.search-box {
  flex: 1;
}

.filter-box select {
  min-width: 240px;
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

.font-mono { font-family: monospace; }
.font-bold { font-weight: 700; }
.text-accent { color: var(--accent-primary); }
.text-muted { color: var(--text-muted); }
.text-center { text-align: center; }
.text-right { text-align: right; }
.cursor-pointer { cursor: pointer; }
.w-12 { width: 3rem; }

@media (max-width: 640px) {
  .toolbar-card {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
