<template>
  <div class="page-container">
    <div class="header-actions">
      <div>
        <h1 class="page-title">Daftar Project Permohonan Saya</h1>
        <p class="page-subtitle">Kelola seluruh draf dan status pengajuan dokumen kelayakan Anda.</p>
      </div>
      <button class="btn btn-primary" @click="openCreateModal">
        + Buat Permohonan Baru
      </button>
    </div>

    <!-- Filter & Search Toolbar -->
    <div class="glass-card toolbar-card">
      <div class="search-box">
        <input 
          v-model="search" 
          type="text" 
          class="form-input" 
          placeholder="Cari berdasarkan judul, nomor project, unit kerja, atau perusahaan..." 
          @input="debouncedFetch"
        />
      </div>

      <div class="filter-box">
        <select v-model="statusFilter" class="form-select" @change="fetchProjects">
          <option value="">Semua Status</option>
          <option value="DRAFT">Draft</option>
          <option value="SUBMITTED">Telah Dikirim</option>
          <option value="UNDER_REVIEW">Sedang Dinilai</option>
          <option value="REVISION">Perlu Revisi</option>
          <option value="APPROVED">Disetujui</option>
          <option value="REJECTED">Ditolak</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-box">Memuat data project...</div>

    <!-- Empty State -->
    <div v-else-if="projects.length === 0" class="glass-card empty-card">
      <h3>Belum Ada Permohonan</h3>
      <p>Tidak ada permohonan yang sesuai dengan filter atau pencarian Anda.</p>
      <button class="btn btn-primary" @click="openCreateModal">Buat Permohonan Pertama</button>
    </div>

    <!-- Table View -->
    <div v-else class="glass-card table-card">
      <table class="data-table">
        <thead>
          <tr>
            <th class="text-center w-12">No</th>
            <th>No. Project</th>
            <th>Judul Permohonan</th>
            <th>Perusahaan & Unit Kerja</th>
            <th>Status</th>
            <th>Tanggal Buat</th>
            <th class="text-right">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in projects" :key="item.id">
            <td class="text-center font-mono text-muted">{{ (meta.page - 1) * meta.per_page + index + 1 }}</td>
            <td class="font-mono text-accent">{{ item.project_number }}</td>
            <td class="font-bold cursor-pointer" @click="$router.push(`/projects/${item.id}`)">
              <div>{{ item.title }}</div>
              <small v-if="item.document_type" class="text-muted font-normal">
                {{ item.document_type.code }} - {{ item.document_type.name }}
              </small>
            </td>
            <td>
              <div>{{ item.company_name || '-' }}</div>
              <small v-if="item.unit" class="text-accent">Unit: {{ item.unit }}</small>
            </td>
            <td><StatusBadge :status="item.status" /></td>
            <td>{{ formatDate(item.created_at) }}</td>
            <td class="text-right">
              <div class="action-buttons">
                <button class="btn btn-secondary btn-sm" @click="$router.push(`/projects/${item.id}`)">
                  Detail
                </button>
                <button 
                  v-if="item.status === 'DRAFT' || item.status === 'REVISION'" 
                  class="btn btn-secondary btn-sm" 
                  @click="openEditModal(item.id)"
                >
                  Edit
                </button>
                <button 
                  v-if="item.status === 'DRAFT'" 
                  class="btn btn-danger btn-sm" 
                  @click="confirmDelete(item)"
                >
                  Hapus
                </button>
              </div>
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

    <!-- 2-Step Wizard Create/Edit Modal Dialog -->
    <CreateProjectModal 
      :is-open="isModalOpen" 
      :edit-project-id="editingProjectId" 
      @close="isModalOpen = false" 
      @created="fetchProjects"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import apiClient from '../services/api'
import StatusBadge from '../components/StatusBadge.vue'
import Pagination from '../components/Pagination.vue'
import CreateProjectModal from '../components/CreateProjectModal.vue'
import { alertSuccess, alertError, confirmDialog } from '../utils/swal'

const route = useRoute()
const loading = ref(true)
const projects = ref([])
const search = ref('')
const statusFilter = ref('')

const isModalOpen = ref(false)
const editingProjectId = ref(null)

const meta = reactive({
  page: 1,
  per_page: 20,
  total: 0,
  total_pages: 1
})

const openCreateModal = () => {
  editingProjectId.value = null
  isModalOpen.value = true
}

const openEditModal = (id) => {
  editingProjectId.value = id
  isModalOpen.value = true
}

const fetchProjects = async (page = 1) => {
  loading.value = true
  try {
    const params = {
      page,
      per_page: meta.per_page
    }
    if (search.value) params.search = search.value
    if (statusFilter.value) params.status = statusFilter.value

    const res = await apiClient.get('/projects', { params })
    if (res.data?.data) {
      projects.value = res.data.data
      if (res.data.meta) {
        meta.page = res.data.meta.page
        meta.total = res.data.meta.total
        meta.total_pages = res.data.meta.total_pages
      }
    }
  } catch (err) {
    console.error('Failed to fetch projects:', err)
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

const confirmDelete = async (project) => {
  const confirmed = await confirmDialog(
    'Hapus Permohonan?',
    `Apakah Anda yakin ingin menghapus draft permohonan "${project.title}"?`,
    'Ya, Hapus Draft'
  )

  if (confirmed) {
    try {
      await apiClient.delete(`/projects/${project.id}`)
      alertSuccess('Terhapus', 'Draft permohonan berhasil dihapus.')
      fetchProjects(meta.page)
    } catch (err) {
      alertError('Gagal Menghapus', err.response?.data?.error || err.message)
    }
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

watch(() => route.query.create, (newVal) => {
  if (newVal === 'true') {
    openCreateModal()
  }
})

onMounted(() => {
  fetchProjects()
  if (route.query.create === 'true') {
    openCreateModal()
  }
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
  min-width: 180px;
}

.loading-box {
  padding: 3rem;
  text-align: center;
  color: var(--text-muted);
}

.empty-card {
  padding: 4rem 2rem;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.85rem;
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
.font-normal { font-weight: 400; }
.text-accent { color: var(--accent-primary); }
.text-muted { color: var(--text-muted); }
.text-center { text-align: center; }
.text-right { text-align: right; }
.cursor-pointer { cursor: pointer; }
.w-12 { width: 3rem; }

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 0.35rem;
}

@media (max-width: 640px) {
  .header-actions, .toolbar-card {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
