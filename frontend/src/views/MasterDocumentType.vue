<template>
  <div class="page-container">
    <div class="header-actions">
      <div>
        <h1 class="page-title">Master Data Jenis Dokumen Kelayakan</h1>
        <p class="page-subtitle">Kelola opsi jenis dokumen, deskripsi target pengerjaan, serta daftar persyaratan berkas wajib.</p>
      </div>
      <button class="btn btn-primary" @click="openCreateModal">
        + Tambah Jenis Dokumen Baru
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-box">Memuat data master jenis dokumen...</div>

    <!-- Empty State -->
    <div v-else-if="docTypes.length === 0" class="glass-card empty-card">
      <h3>Belum Ada Data Master Jenis Dokumen</h3>
      <p>Klik tombol di atas untuk membuat jenis dokumen kelayakan pertama.</p>
    </div>

    <!-- Table View -->
    <div v-else class="glass-card table-card">
      <table class="data-table">
        <thead>
          <tr>
            <th class="text-center w-12">No</th>
            <th>Kode</th>
            <th>Nama Jenis Dokumen</th>
            <th>Deskripsi & Target Pengerjaan</th>
            <th>Daftar Persyaratan Wajib</th>
            <th>Status</th>
            <th class="text-right">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(dt, index) in docTypes" :key="dt.id">
            <td class="text-center font-mono text-muted">{{ index + 1 }}</td>
            <td class="font-mono text-accent font-bold">{{ dt.code }}</td>
            <td class="font-bold">{{ dt.name }}</td>
            <td>
              <div class="desc-preview html-content" v-html="dt.description || '-'"></div>
            </td>
            <td>
              <div class="req-preview html-content" v-html="dt.requirement || '-'"></div>
            </td>
            <td>
              <span class="status-badge" :class="dt.is_active ? 'badge-active' : 'badge-inactive'">
                {{ dt.is_active ? 'Aktif' : 'Non-Aktif' }}
              </span>
            </td>
            <td class="text-right">
              <div class="action-buttons">
                <button class="btn btn-secondary btn-sm" @click="openEditModal(dt)">Edit</button>
                <button class="btn btn-danger btn-sm" @click="confirmDelete(dt)">Hapus</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Form (Create / Edit) -->
    <div v-if="isModalOpen" class="modal-backdrop" @click.self="isModalOpen = false">
      <div class="modal-content modal-lg">
        <div class="modal-header">
          <h3 class="modal-title">{{ isEdit ? 'Edit Master Jenis Dokumen' : 'Tambah Master Jenis Dokumen Baru' }}</h3>
          <button class="btn-close" @click="isModalOpen = false">&times;</button>
        </div>

        <form @submit.prevent="handleSubmit">
          <div class="modal-body">
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Kode Singkat *</label>
                <input 
                  v-model="form.code" 
                  type="text" 
                  class="form-input" 
                  placeholder="Contoh: AMDAL, SLF, UKL-UPL" 
                  required 
                />
              </div>
              <div class="form-group">
                <label class="form-label">Nama Jenis Dokumen Lengkap *</label>
                <input 
                  v-model="form.name" 
                  type="text" 
                  class="form-input" 
                  placeholder="Contoh: Analisis Mengenai Dampak Lingkungan (AMDAL)" 
                  required 
                />
              </div>
            </div>

            <!-- Tiptap Rich Text Editor for Requirements -->
            <div class="form-group">
              <label class="form-label">Daftar Persyaratan Berkas Wajib *</label>
              <TiptapEditor 
                v-model="form.requirement" 
                placeholder="Tuliskan rincian dokumen wajib (gunakan Bullet/Numbered list)..." 
              />
            </div>

            <!-- Tiptap Rich Text Editor for Description -->
            <div class="form-group">
              <label class="form-label">Deskripsi & Target Pengerjaan</label>
              <TiptapEditor 
                v-model="form.description" 
                placeholder="Penjelasan umum mengenai ruang lingkup & estimasi target waktu pengerjaan verifikasi..." 
              />
            </div>

            <div class="form-group checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="form.is_active" />
                <span>Status Aktif (Tampilkan di pilihan form pengajuan pemohon)</span>
              </label>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="isModalOpen = false">Batal</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span>{{ submitting ? 'Menyimpan...' : 'Simpan Master Data' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import apiClient from '../services/api'
import TiptapEditor from '../components/TiptapEditor.vue'
import { alertSuccess, alertError, alertWarning, confirmDialog } from '../utils/swal'

const loading = ref(true)
const submitting = ref(false)
const docTypes = ref([])

const isModalOpen = ref(false)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({
  code: '',
  name: '',
  requirement: '',
  description: '',
  is_active: true
})

const resetForm = () => {
  form.code = ''
  form.name = ''
  form.requirement = ''
  form.description = ''
  form.is_active = true
  isEdit.value = false
  currentId.value = null
}

const fetchDocTypes = async () => {
  loading.value = true
  try {
    const res = await apiClient.get('/document-types')
    if (res.data?.data) {
      docTypes.value = res.data.data
    }
  } catch (err) {
    console.error('Failed to fetch document types:', err)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  resetForm()
  isModalOpen.value = true
}

const openEditModal = (dt) => {
  isEdit.value = true
  currentId.value = dt.id
  form.code = dt.code
  form.name = dt.name
  form.requirement = dt.requirement
  form.description = dt.description
  form.is_active = dt.is_active
  isModalOpen.value = true
}

const handleSubmit = async () => {
  if (!form.requirement || form.requirement === '<p></p>') {
    alertWarning('Form Tidak Lengkap', 'Harap isi rincian berkas persyaratan wajib!')
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await apiClient.put(`/document-types/${currentId.value}`, form)
      alertSuccess('Berhasil!', 'Master jenis dokumen berhasil diperbarui!')
    } else {
      await apiClient.post('/document-types', form)
      alertSuccess('Berhasil!', 'Master jenis dokumen baru berhasil ditambahkan!')
    }
    isModalOpen.value = false
    fetchDocTypes()
  } catch (err) {
    alertError('Gagal Menyimpan', err.response?.data?.error || err.message)
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async (dt) => {
  const confirmed = await confirmDialog(
    'Hapus Master Jenis Dokumen?',
    `Apakah Anda yakin ingin menghapus "${dt.name}"?`,
    'Ya, Hapus Data'
  )

  if (confirmed) {
    try {
      await apiClient.delete(`/document-types/${dt.id}`)
      alertSuccess('Terhapus', 'Jenis dokumen berhasil dihapus.')
      fetchDocTypes()
    } catch (err) {
      alertError('Gagal Menghapus', err.response?.data?.error || err.message)
    }
  }
}

onMounted(() => fetchDocTypes())
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

.desc-preview {
  max-width: 260px;
  font-size: 0.82rem;
  color: var(--text-main);
}

.req-preview {
  max-width: 320px;
  font-size: 0.82rem;
  color: var(--text-muted);
}

.html-content :deep(p) { margin-bottom: 0.35rem; }
.html-content :deep(ul), .html-content :deep(ol) { padding-left: 1.2rem; margin-bottom: 0.35rem; }

.status-badge {
  padding: 0.25rem 0.65rem;
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 600;
}
.badge-active { background: rgba(16, 185, 129, 0.15); color: #34d399; border: 1px solid rgba(16, 185, 129, 0.3); }
.badge-inactive { background: rgba(148, 163, 184, 0.15); color: #94a3b8; border: 1px solid rgba(148, 163, 184, 0.3); }

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 0.35rem;
}

/* Modal */
.modal-lg { max-width: 720px; }

.btn-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 1rem;
}

.checkbox-group { margin-top: 0.5rem; }
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.88rem;
  color: var(--text-main);
  cursor: pointer;
}
</style>
