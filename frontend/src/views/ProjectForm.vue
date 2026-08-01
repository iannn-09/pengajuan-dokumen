<template>
  <div class="form-container">
    <div class="header-nav">
      <button class="btn btn-secondary btn-sm" @click="$router.back()">&larr; Kembali</button>
      <h1 class="page-title">{{ isEdit ? 'Edit Permohonan Dokumen' : 'Buat Permohonan Dokumen Baru' }}</h1>
    </div>

    <div class="glass-card form-card">
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label class="form-label">Jenis Dokumen Kelayakan *</label>
          <select v-model="form.document_type_id" class="form-select" @change="onDocTypeChange" required>
            <option :value="null">-- Pilih Jenis Dokumen Kelayakan --</option>
            <option v-for="dt in docTypes" :key="dt.id" :value="dt.id">
              {{ dt.code }} - {{ dt.name }}
            </option>
          </select>
        </div>

        <!-- Dynamic Requirements Info Box -->
        <div v-if="selectedDocType" class="doc-type-info-box">
          <div class="info-header">
            <span class="info-code">{{ selectedDocType.code }}</span>
          </div>
          <h4 class="info-title">{{ selectedDocType.name }}</h4>
          <div v-if="selectedDocType.description" class="info-desc html-content" v-html="selectedDocType.description">
          </div>
          <div class="req-box">
            <strong>📋 Berkas Persyaratan Wajib yang Harus Dilampirkan:</strong>
            <div class="req-text html-content" v-html="selectedDocType.requirement"></div>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">Judul Permohonan Dokumen *</label>
          <input v-model="form.title" type="text" class="form-input"
            placeholder="Contoh: Permohonan Dokumen Kelayakan Lingkungan PT XYZ" required />
        </div>

        <div class="form-group">
          <label class="form-label">Nama Perusahaan / Pemohon *</label>
          <input v-model="form.company_name" type="text" class="form-input" placeholder="Nama Resmi PT / CV / Instansi"
            required />
        </div>

        <div class="form-group">
          <label class="form-label">Deskripsi Permohonan & Detail Kelayakan</label>
          <textarea v-model="form.description" class="form-textarea"
            placeholder="Jelaskan ruang lingkup, lokasi, dan latar belakang pengajuan permohonan dokumen kelayakan..."
            rows="4"></textarea>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="$router.back()">Batal</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">
            <span>{{ submitting ? 'Menyimpan...' : (isEdit ? 'Simpan Perubahan' : 'Simpan Sebagai Draft') }}</span>
          </button>
        </div>
      </form>
    </div>

    <!-- Document Attachment Section (Available when project exists) -->
    <div v-if="isEdit && project" class="glass-card upload-card">
      <h3 class="section-title">Dokumen Lampiran (PDF / Gambar / DOCX)</h3>
      <p class="section-subtitle">Unggah dokumen kelayakan pendukung (Maksimal 10MB per file).</p>

      <!-- Upload Form -->
      <div class="file-uploader">
        <input type="file" ref="fileInput" class="file-input-hidden" accept=".pdf,.jpg,.jpeg,.png,.doc,.docx"
          @change="handleFileUpload" />
        <div class="upload-dropzone" @click="$refs.fileInput.click()">
          <span class="upload-icon">📁</span>
          <span class="upload-text">Klik di sini untuk memilih dan mengunggah dokumen</span>
          <small class="upload-hint">Format yang diizinkan: PDF, JPG, PNG, DOC, DOCX</small>
        </div>
      </div>

      <!-- File List -->
      <div v-if="documents.length > 0" class="file-list">
        <div v-for="doc in documents" :key="doc.id" class="file-item">
          <div class="file-info">
            <span class="file-name">{{ doc.file_name }}</span>
            <span class="file-size">{{ formatSize(doc.file_size) }}</span>
          </div>
          <div class="file-actions">
            <a :href="getDownloadUrl(doc.id)" target="_blank" class="btn btn-secondary btn-sm">Lihat</a>
            <button type="button" class="btn btn-danger btn-sm" @click="deleteDoc(doc.id)">Hapus</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import { alertSuccess, alertError, alertWarning, confirmDialog } from '../utils/swal'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const isEdit = computed(() => !!route.params.id)
const project = ref(null)
const documents = ref([])
const docTypes = ref([])
const selectedDocType = ref(null)

const submitting = ref(false)
const fileInput = ref(null)

const form = reactive({
  title: '',
  company_name: '',
  description: '',
  document_type_id: null
})

const fetchDocTypes = async () => {
  try {
    const res = await apiClient.get('/document-types?active_only=true')
    if (res.data?.data) {
      docTypes.value = res.data.data
      if (form.document_type_id) {
        onDocTypeChange()
      }
    }
  } catch (err) {
    console.error('Failed to fetch document types:', err)
  }
}

const onDocTypeChange = () => {
  if (!form.document_type_id) {
    selectedDocType.value = null
    return
  }
  selectedDocType.value = docTypes.value.find(d => d.id === form.document_type_id) || null
}

const fetchProject = async () => {
  if (!isEdit.value) return
  try {
    const res = await apiClient.get(`/projects/${route.params.id}`)
    if (res.data?.data) {
      project.value = res.data.data
      form.title = project.value.title
      form.company_name = project.value.company_name
      form.description = project.value.description
      form.document_type_id = project.value.document_type_id || null
      documents.value = project.value.documents || []
      onDocTypeChange()
    }
  } catch (err) {
    alertError('Gagal Memuat', err.response?.data?.error || err.message)
    router.push('/projects')
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEdit.value) {
      await apiClient.put(`/projects/${route.params.id}`, form)
      alertSuccess('Berhasil!', 'Permohonan berhasil diperbarui.')
    } else {
      const res = await apiClient.post('/projects', form)
      alertSuccess('Draft Dibuat!', 'Draft permohonan berhasil dibuat. Silakan unggah dokumen pendukung.')
      router.push(`/projects/${res.data.data.id}/edit`)
    }
  } catch (err) {
    alertError('Gagal Menyimpan', err.response?.data?.error || err.message)
  } finally {
    submitting.value = false
  }
}

const handleFileUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  try {
    await apiClient.post(`/projects/${route.params.id}/documents`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    alertSuccess('Upload Berhasil', `File "${file.name}" berhasil diunggah.`)
    fetchProject()
    fileInput.value.value = ''
  } catch (err) {
    alertError('Gagal Upload Dokumen', err.response?.data?.error || err.message)
  }
}

const deleteDoc = async (docId) => {
  const confirmed = await confirmDialog('Hapus Dokumen?', 'Apakah Anda yakin ingin menghapus dokumen ini?', 'Ya, Hapus File')
  if (confirmed) {
    try {
      await apiClient.delete(`/projects/${route.params.id}/documents/${docId}`)
      alertSuccess('Terhapus', 'Dokumen lampiran berhasil dihapus.')
      fetchProject()
    } catch (err) {
      alertError('Gagal Menghapus', err.response?.data?.error || err.message)
    }
  }
}

const getDownloadUrl = (docId) => {
  return `${apiClient.defaults.baseURL}/documents/${docId}/download?token=${auth.token}`
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(async () => {
  await fetchDocTypes()
  await fetchProject()
})
</script>

<style scoped>
.form-container {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.header-nav {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.page-title {
  font-size: 1.3rem;
  font-weight: 800;
  color: var(--text-main);
}

.form-card,
.upload-card {
  padding: 1.75rem;
}

.doc-type-info-box {
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: var(--radius-md);
  padding: 1.25rem;
  margin-bottom: 1.5rem;
}

.info-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.info-code {
  font-family: monospace;
  font-size: 0.85rem;
  font-weight: 800;
  color: var(--accent-primary);
}

.info-title {
  font-size: 1.05rem;
  font-weight: 800;
  color: var(--text-main);
}

.info-desc {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-top: 0.25rem;
}

.req-box {
  margin-top: 0.85rem;
  padding-top: 0.75rem;
  border-top: 1px dashed rgba(99, 102, 241, 0.3);
  font-size: 0.85rem;
  color: var(--text-main);
}

.req-text {
  font-size: 0.83rem;
  color: var(--text-muted);
  margin-top: 0.35rem;
  background: rgba(0, 0, 0, 0.2);
  padding: 0.75rem 1rem;
  border-radius: var(--radius-sm);
}

.html-content :deep(p) {
  margin-bottom: 0.35rem;
}

.html-content :deep(ul),
.html-content :deep(ol) {
  padding-left: 1.2rem;
  margin-bottom: 0.35rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
}

.section-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-main);
}

.section-subtitle {
  font-size: 0.83rem;
  color: var(--text-muted);
  margin-bottom: 1.25rem;
}

.file-input-hidden {
  display: none;
}

.upload-dropzone {
  border: 2px dashed var(--border-color);
  border-radius: var(--radius-md);
  padding: 2rem;
  text-align: center;
  cursor: pointer;
  background: var(--bg-input);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.2s ease;
}

.upload-dropzone:hover {
  border-color: var(--accent-primary);
  background: rgba(99, 102, 241, 0.05);
}

.upload-icon {
  font-size: 2rem;
}

.upload-text {
  font-weight: 600;
  color: var(--text-main);
  font-size: 0.9rem;
}

.upload-hint {
  color: var(--text-subtle);
  font-size: 0.78rem;
}

.file-list {
  margin-top: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.file-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.file-info {
  display: flex;
  flex-direction: column;
}

.file-name {
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--text-main);
}

.file-size {
  font-size: 0.75rem;
  color: var(--text-subtle);
}

.file-actions {
  display: flex;
  gap: 0.4rem;
}
</style>
