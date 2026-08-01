<template>
  <div class="page-container">
    <div class="header-actions">
      <div>
        <h1 class="page-title">Manajemen User & Penilai</h1>
        <p class="page-subtitle">Kelola daftar akun pengguna, pemohon, serta akun Penilai / Verifikator dalam sistem.</p>
      </div>
      <button v-if="auth.isAdmin" class="btn btn-primary" @click="isCreateModalOpen = true">
        + Buat Akun Penilai / User Baru
      </button>
    </div>

    <!-- Filter & Search Toolbar -->
    <div class="glass-card toolbar-card">
      <div class="search-box">
        <input 
          v-model="search" 
          type="text" 
          class="form-input" 
          placeholder="Cari nama, email, atau instansi/perusahaan..." 
          @input="debouncedFetch"
        />
      </div>

      <div class="filter-box">
        <select v-model="roleFilter" class="form-select" @change="fetchUsers">
          <option value="">Semua Role</option>
          <option value="admin">Admin</option>
          <option value="penilai">Penilai / Verifikator</option>
          <option value="pemohon">Pemohon Dokumen</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-box">Memuat data pengguna...</div>

    <!-- Empty State -->
    <div v-else-if="users.length === 0" class="glass-card empty-card">
      <h3>Tidak Ada User</h3>
      <p>Tidak ada data pengguna yang sesuai dengan filter atau pencarian Anda.</p>
    </div>

    <!-- Table View -->
    <div v-else class="glass-card table-card">
      <table class="data-table">
        <thead>
          <tr>
            <th class="text-center w-12">No</th>
            <th>Nama Lengkap</th>
            <th>Email</th>
            <th>Role / Hak Akses</th>
            <th>Perusahaan / Instansi</th>
            <th>Telepon</th>
            <th>Tanggal Buat</th>
            <th v-if="auth.isAdmin" class="text-right">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(user, index) in users" :key="user.id">
            <td class="text-center font-mono text-muted">{{ (meta.page - 1) * meta.per_page + index + 1 }}</td>
            <td class="font-bold">{{ user.name }}</td>
            <td class="font-mono text-accent">{{ user.email }}</td>
            <td>
              <span class="role-badge" :class="`badge-${user.role}`">
                {{ getRoleText(user.role) }}
              </span>
            </td>
            <td>{{ user.company || '-' }}</td>
            <td>{{ user.phone || '-' }}</td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td v-if="auth.isAdmin" class="text-right">
              <button 
                class="btn btn-danger btn-sm" 
                :disabled="user.id === auth.user?.id"
                @click="confirmDelete(user)"
                title="Hapus User"
              >
                Hapus
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

    <!-- Create User Modal (Admin only) -->
    <div v-if="isCreateModalOpen" class="modal-backdrop" @click.self="isCreateModalOpen = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Tambah Akun Penilai / User Baru</h3>
          <button class="btn-close" @click="isCreateModalOpen = false">&times;</button>
        </div>

        <form @submit.prevent="handleCreateUser">
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">Role / Hak Akses *</label>
              <select v-model="form.role" class="form-select" required>
                <option value="penilai">Penilai / Verifikator Dokumen</option>
                <option value="pemohon">Pemohon Dokumen</option>
                <option value="admin">Administrator</option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label">Nama Lengkap *</label>
              <input 
                v-model="form.name" 
                type="text" 
                class="form-input" 
                placeholder="Nama Pengguna" 
                required 
              />
            </div>

            <div class="form-group">
              <label class="form-label">Email Login *</label>
              <input 
                v-model="form.email" 
                type="email" 
                class="form-input" 
                placeholder="email@kelayakan.id" 
                required 
              />
            </div>

            <div class="form-group">
              <label class="form-label">Password * (Min 8 Karakter)</label>
              <input 
                v-model="form.password" 
                type="password" 
                class="form-input" 
                placeholder="••••••••" 
                required 
                minlength="8"
              />
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Nomor Telepon</label>
                <input v-model="form.phone" type="text" class="form-input" placeholder="0812xxxx" />
              </div>
              <div class="form-group">
                <label class="form-label">Instansi / Departemen</label>
                <input v-model="form.company" type="text" class="form-input" placeholder="Contoh: Dinas Lingkungan Hidup" />
              </div>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="isCreateModalOpen = false">Batal</button>
            <button type="submit" class="btn btn-primary" :disabled="creating">
              <span>{{ creating ? 'Menyimpan...' : 'Buat Akun User' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import Pagination from '../components/Pagination.vue'
import { alertSuccess, alertError, confirmDialog } from '../utils/swal'

const auth = useAuthStore()

const loading = ref(true)
const creating = ref(false)
const users = ref([])
const search = ref('')
const roleFilter = ref('')

const isCreateModalOpen = ref(false)

const meta = reactive({
  page: 1,
  per_page: 20,
  total: 0,
  total_pages: 1
})

const form = reactive({
  name: '',
  email: '',
  password: '',
  role: 'penilai',
  phone: '',
  company: 'Dinas Lingkungan Hidup'
})

const resetForm = () => {
  form.name = ''
  form.email = ''
  form.password = ''
  form.role = 'penilai'
  form.phone = ''
  form.company = 'Dinas Lingkungan Hidup'
}

const fetchUsers = async (page = 1) => {
  loading.value = true
  try {
    const params = { page, per_page: meta.per_page }
    if (search.value) params.search = search.value
    if (roleFilter.value) params.role = roleFilter.value

    const res = await apiClient.get('/users', { params })
    if (res.data?.data) {
      users.value = res.data.data
      if (res.data.meta) {
        meta.page = res.data.meta.page
        meta.total = res.data.meta.total
        meta.total_pages = res.data.meta.total_pages
      }
    }
  } catch (err) {
    console.error('Failed to fetch users:', err)
  } finally {
    loading.value = false
  }
}

let timeout = null
const debouncedFetch = () => {
  clearTimeout(timeout)
  timeout = setTimeout(() => fetchUsers(1), 300)
}

const onPageChange = (newPage) => fetchUsers(newPage)

const handleCreateUser = async () => {
  creating.value = true
  try {
    await apiClient.post('/users', form)
    alertSuccess('Akun Dibuat!', `Berhasil membuat akun ${form.role.toUpperCase()} untuk ${form.email}.`)
    isCreateModalOpen.value = false
    resetForm()
    fetchUsers()
  } catch (err) {
    alertError('Gagal Membuat Akun', err.response?.data?.error || err.message)
  } finally {
    creating.value = false
  }
}

const confirmDelete = async (user) => {
  const confirmed = await confirmDialog(
    'Hapus Akun User?',
    `Apakah Anda yakin ingin menghapus akun "${user.name}" (${user.email})?`,
    'Ya, Hapus Akun'
  )

  if (confirmed) {
    try {
      await apiClient.delete(`/users/${user.id}`)
      alertSuccess('Terhapus', 'Akun pengguna berhasil dihapus.')
      fetchUsers(meta.page)
    } catch (err) {
      alertError('Gagal Menghapus', err.response?.data?.error || err.message)
    }
  }
}

const getRoleText = (role) => {
  if (role === 'admin') return 'Admin'
  if (role === 'penilai') return 'Penilai / Verifikator'
  return 'Pemohon'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(() => fetchUsers())
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

.search-box { flex: 1; }
.filter-box select { min-width: 180px; }

.loading-box, .empty-card {
  padding: 4rem;
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

.font-mono { font-family: monospace; font-size: 0.85rem; }
.font-bold { font-weight: 700; }
.text-accent { color: var(--accent-primary); }
.text-muted { color: var(--text-muted); }
.text-center { text-align: center; }
.text-right { text-align: right; }
.w-12 { width: 3rem; }

.role-badge {
  padding: 0.25rem 0.65rem;
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 600;
}
.badge-admin { background: rgba(168, 85, 247, 0.15); color: #c084fc; border: 1px solid rgba(168, 85, 247, 0.3); }
.badge-penilai { background: rgba(59, 130, 246, 0.15); color: #60a5fa; border: 1px solid rgba(59, 130, 246, 0.3); }
.badge-pemohon { background: rgba(16, 185, 129, 0.15); color: #34d399; border: 1px solid rgba(16, 185, 129, 0.3); }

/* Modal */
.btn-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}
</style>
