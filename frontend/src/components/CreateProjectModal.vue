<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="closeModal">
    <div class="modal-content modal-lg">
      <div class="modal-header">
        <div class="header-titles">
          <h3 class="modal-title">{{ isEdit ? 'Edit Draft Permohonan Dokumen' : 'Buat Permohonan Dokumen Baru' }}</h3>
          <p class="modal-subtitle">Pengajuan dokumen kelayakan 2-Step Wizard</p>
        </div>
        <button class="btn-close" @click="closeModal">&times;</button>
      </div>

      <!-- Stepper Progress Bar -->
      <div class="stepper-bar">
        <div class="step-item" :class="{ active: currentStep === 1, completed: currentStep > 1 }" @click="goToStep(1)">
          <div class="step-badge">1</div>
          <span class="step-label">1. Data Permohonan</span>
        </div>
        <div class="step-line" :class="{ active: currentStep > 1 }"></div>
        <div class="step-item" :class="{ active: currentStep === 2, completed: createdProjectId && currentStep === 2 }"
          @click="goToStep(2)">
          <div class="step-badge">2</div>
          <span class="step-label">2. Upload Lampiran Berkas</span>
        </div>
      </div>

      <div class="modal-body">
        <!-- ─── STEP 1: Form Data Permohonan ─────────────────────── -->
        <div v-if="currentStep === 1" class="step-content">
          <form @submit.prevent="proceedToStep2">
            <div class="form-group">
              <label class="form-label">Jenis Dokumen Kelayakan *</label>
              <select v-model="form.document_type_id" class="form-select" @change="onDocTypeChange" required>
                <option :value="null">-- Pilih Jenis Dokumen Kelayakan --</option>
                <option v-for="dt in docTypes" :key="dt.id" :value="dt.id">
                  {{ dt.code }} - {{ dt.name }}
                </option>
              </select>
            </div>

            <!-- Dynamic Requirements Box -->
            <div v-if="selectedDocType" class="doc-type-info-box">
              <div class="info-header">
                <span class="info-code">{{ selectedDocType.code }}</span>
              </div>
              <h4 class="info-title">{{ selectedDocType.name }}</h4>
              <div v-if="selectedDocType.description" class="info-desc html-content"
                v-html="selectedDocType.description"></div>
              <div class="req-box">
                <strong>📋 Berkas Persyaratan Wajib:</strong>
                <div class="req-text html-content" v-html="selectedDocType.requirement"></div>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">Nama Perusahaan / PT (Otomatis)</label>
                <input :value="form.company_name" type="text" class="form-input form-disabled" disabled
                  placeholder="Terisi Otomatis Dari Akun Register" />
                <small class="form-hint">🔒 Nama PT diambil otomatis dari profil akun Anda.</small>
              </div>

              <div class="form-group">
                <label class="form-label">Unit Kerja / Divisi / Pabrik *</label>
                <input v-model="form.unit" type="text" class="form-input"
                  placeholder="Contoh: Unit Pengolahan 2 / Divisi Ops B" required />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Judul Permohonan Dokumen *</label>
              <input v-model="form.title" type="text" class="form-input"
                placeholder="Contoh: Permohonan Dokumen Kelayakan Lingkungan Pabrik Sektor A" required />
            </div>

            <div class="form-group">
              <label class="form-label">Deskripsi & Ruang Lingkup Permohonan</label>
              <textarea v-model="form.description" class="form-textarea"
                placeholder="Jelaskan ruang lingkup, lokasi, dan latar belakang permohonan..." rows="3"></textarea>
            </div>

            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" @click="closeModal">Batal</button>
              <button type="submit" class="btn btn-primary" :disabled="submitting">
                <span>{{ submitting ? 'Memproses...' : 'Lanjut ke Step 2: Upload Lampiran →' }}</span>
              </button>
            </div>
          </form>
        </div>

        <!-- ─── STEP 2: Upload Lampiran Dokumen ───────────────────── -->
        <div v-else-if="currentStep === 2" class="step-content">
          <!-- Requirements Reminder Box in Step 2 -->
          <div v-if="selectedDocType" class="doc-type-info-box">
            <div class="info-header">
              <span class="info-code">{{ selectedDocType.code }}</span>
            </div>
            <h4 class="info-title">{{ selectedDocType.name }}</h4>
            <div v-if="selectedDocType.description" class="info-desc html-content" v-html="selectedDocType.description">
            </div>
            <div class="req-box">
              <strong>📋 Berkas Persyaratan Wajib yang Harus Diunggah:</strong>
              <div class="req-text html-content" v-html="selectedDocType.requirement"></div>
            </div>
          </div>

          <div class="upload-section">
            <h4 class="section-title">Unggah Dokumen Lampiran Wajib</h4>
            <p class="section-subtitle">Silakan unggah dokumen persyaratan di atas (PDF, Gambar, atau DOCX - Maks 10MB
              per file).</p>

            <!-- Upload Dropzone -->
            <div class="file-uploader">
              <input type="file" ref="fileInput" class="file-input-hidden" accept=".pdf,.jpg,.jpeg,.png,.doc,.docx"
                @change="handleFileUpload" />
              <div class="upload-dropzone" @click="$refs.fileInput.click()">
                <span class="upload-icon">📁</span>
                <span class="upload-text">Klik di sini untuk memilih dan mengunggah dokumen</span>
                <small class="upload-hint">Format yang diizinkan: PDF, JPG, PNG, DOC, DOCX</small>
              </div>
            </div>

            <!-- Uploaded File List -->
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
            <div v-else class="empty-upload-notice">
              ⚠️ Belum ada dokumen lampiran yang diunggah. Disarankan mengunggah minimal 1 dokumen pendukung.
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="currentStep = 1">&larr; Kembali ke Step 1</button>
            <button type="button" class="btn btn-secondary" @click="saveDraftAndClose">Simpan Sebagai Draft</button>
            <button type="button" class="btn btn-primary" :disabled="submitting" @click="submitFinalProject">
              <span>{{ submitting ? 'Mengirimkan...' : 'Kirim Permohonan Sekarang ✓' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import { alertSuccess, alertError, alertWarning, confirmDialog } from '../utils/swal'

const props = defineProps({
  isOpen: { type: Boolean, default: false },
  editProjectId: { type: [Number, String], default: null }
})

const emit = defineEmits(['close', 'created', 'updated'])

const auth = useAuthStore()
const currentStep = ref(1)
const submitting = ref(false)
const createdProjectId = ref(null)
const docTypes = ref([])
const selectedDocType = ref(null)
const documents = ref([])
const fileInput = ref(null)

const isEdit = computed(() => !!props.editProjectId)

const form = reactive({
  title: '',
  company_name: '',
  unit: '',
  description: '',
  document_type_id: null
})

const resetForm = () => {
  currentStep.value = 1
  createdProjectId.value = null
  form.title = ''
  form.company_name = auth.userCompany || auth.user?.company || ''
  form.unit = ''
  form.description = ''
  form.document_type_id = null
  selectedDocType.value = null
  documents.value = []
}

const fetchDocTypes = async () => {
  try {
    const res = await apiClient.get('/document-types?active_only=true')
    if (res.data?.data) {
      docTypes.value = res.data.data
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

const fetchProjectDetail = async (id) => {
  try {
    const res = await apiClient.get(`/projects/${id}`)
    if (res.data?.data) {
      const proj = res.data.data
      createdProjectId.value = proj.id
      form.title = proj.title
      form.company_name = proj.company_name || auth.userCompany || ''
      form.unit = proj.unit || ''
      form.description = proj.description || ''
      form.document_type_id = proj.document_type_id || null
      documents.value = proj.documents || []
      onDocTypeChange()
    }
  } catch (err) {
    alertError('Gagal Memuat Detail', err.response?.data?.error || err.message)
  }
}

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    fetchDocTypes()
    if (props.editProjectId) {
      fetchProjectDetail(props.editProjectId)
    } else {
      resetForm()
    }
  }
})

const goToStep = (step) => {
  if (step === 2 && !createdProjectId.value) {
    alertWarning('Lengkapi Step 1', 'Silakan isi data permohonan terlebih dahulu sebelum lanjut ke upload lampiran.')
    return
  }
  currentStep.value = step
}

const proceedToStep2 = async () => {
  submitting.value = true
  try {
    if (createdProjectId.value) {
      // Update existing draft project
      await apiClient.put(`/projects/${createdProjectId.value}`, form)
    } else {
      // Create new draft project
      const res = await apiClient.post('/projects', form)
      createdProjectId.value = res.data.data.id
    }
    currentStep.value = 2
  } catch (err) {
    alertError('Gagal Menyimpan', err.response?.data?.error || err.message)
  } finally {
    submitting.value = false
  }
}

const handleFileUpload = async (event) => {
  const file = event.target.files[0]
  if (!file || !createdProjectId.value) return

  const formData = new FormData()
  formData.append('file', file)

  try {
    await apiClient.post(`/projects/${createdProjectId.value}/documents`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    alertSuccess('Upload Berhasil', `File "${file.name}" berhasil diunggah.`)
    fetchProjectDetail(createdProjectId.value)
    if (fileInput.value) fileInput.value.value = ''
  } catch (err) {
    alertError('Gagal Upload Dokumen', err.response?.data?.error || err.message)
  }
}

const deleteDoc = async (docId) => {
  const confirmed = await confirmDialog('Hapus Dokumen?', 'Apakah Anda yakin ingin menghapus dokumen lampiran ini?', 'Ya, Hapus File')
  if (confirmed && createdProjectId.value) {
    try {
      await apiClient.delete(`/projects/${createdProjectId.value}/documents/${docId}`)
      alertSuccess('Terhapus', 'Dokumen lampiran berhasil dihapus.')
      fetchProjectDetail(createdProjectId.value)
    } catch (err) {
      alertError('Gagal Menghapus', err.response?.data?.error || err.message)
    }
  }
}

const saveDraftAndClose = () => {
  alertSuccess('Draft Tersimpan', 'Permohonan Anda telah disimpan sebagai draft.')
  emit('created')
  closeModal()
}

const submitFinalProject = async () => {
  if (!createdProjectId.value) return

  if (documents.value.length === 0) {
    alertWarning('Lampiran Belum Diunggah', 'Disarankan mengunggah minimal 1 dokumen lampiran pendukung!')
  }

  const confirmed = await confirmDialog(
    'Kirimkan Permohonan Sekarang?',
    'Kirimkan permohonan dokumen ini kepada Penilai untuk langsung diverifikasi?',
    'Ya, Kirim Sekarang'
  )

  if (confirmed) {
    submitting.value = true
    try {
      await apiClient.post(`/projects/${createdProjectId.value}/submit`)
      alertSuccess('Berhasil Terkirim!', 'Permohonan Anda berhasil dikirimkan ke Penilai.')
      emit('created')
      closeModal()
    } catch (err) {
      alertError('Gagal Mengirimkan', err.response?.data?.error || err.message)
    } finally {
      submitting.value = false
    }
  }
}

const closeModal = () => {
  emit('close')
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
</script>

<style scoped>
.modal-lg {
  max-width: 760px;
}

.header-titles {
  display: flex;
  flex-direction: column;
}

.modal-subtitle {
  font-size: 0.82rem;
  color: var(--text-muted);
  margin-top: 0.15rem;
}

.btn-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
}

/* Stepper Bar */
.stepper-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.85rem 1.5rem;
  background: rgba(15, 23, 42, 0.5);
  border-bottom: 1px solid var(--border-color);
}

.step-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  opacity: 0.5;
  transition: all 0.2s ease;
}

.step-item.active,
.step-item.completed {
  opacity: 1;
}

.step-badge {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  color: var(--text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 700;
}

.step-item.active .step-badge {
  background: var(--accent-primary);
  border-color: var(--accent-primary);
  color: white;
}

.step-item.completed .step-badge {
  background: #10b981;
  border-color: #10b981;
  color: white;
}

.step-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-main);
}

.step-line {
  height: 2px;
  width: 50px;
  background: var(--border-color);
  transition: all 0.2s ease;
}

.step-line.active {
  background: var(--accent-primary);
}

/* Form Styles */
.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-disabled {
  background: rgba(255, 255, 255, 0.04);
  color: var(--text-muted);
  cursor: not-allowed;
}

.form-hint {
  font-size: 0.75rem;
  color: var(--text-subtle);
  margin-top: 0.25rem;
  display: block;
}

.doc-type-info-box {
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: var(--radius-md);
  padding: 1.1rem;
  margin-bottom: 1.25rem;
}

.info-header {
  display: flex;
  justify-content: flex-end;
}

.info-code {
  font-family: monospace;
  font-size: 0.85rem;
  font-weight: 800;
  color: var(--accent-primary);
}

.info-title {
  font-size: 1rem;
  font-weight: 800;
  color: var(--text-main);
}

.info-desc {
  font-size: 0.83rem;
  color: var(--text-muted);
  margin-top: 0.2rem;
}

.req-box {
  margin-top: 0.75rem;
  padding-top: 0.65rem;
  border-top: 1px dashed rgba(99, 102, 241, 0.3);
  font-size: 0.83rem;
  color: var(--text-main);
}

.req-text {
  font-size: 0.82rem;
  color: var(--text-muted);
  margin-top: 0.35rem;
  background: rgba(0, 0, 0, 0.2);
  padding: 0.65rem 0.85rem;
  border-radius: var(--radius-sm);
}

.html-content :deep(p) {
  margin-bottom: 0.3rem;
}

.html-content :deep(ul),
.html-content :deep(ol) {
  padding-left: 1.2rem;
  margin-bottom: 0.3rem;
}

/* Step 2 Upload Styles */
.section-title {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-main);
}

.section-subtitle {
  font-size: 0.82rem;
  color: var(--text-muted);
  margin-bottom: 1rem;
}

.file-input-hidden {
  display: none;
}

.upload-dropzone {
  border: 2px dashed var(--border-color);
  border-radius: var(--radius-md);
  padding: 1.75rem;
  text-align: center;
  cursor: pointer;
  background: var(--bg-input);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.4rem;
  transition: all 0.2s ease;
}

.upload-dropzone:hover {
  border-color: var(--accent-primary);
  background: rgba(99, 102, 241, 0.05);
}

.upload-icon {
  font-size: 1.8rem;
}

.upload-text {
  font-weight: 600;
  color: var(--text-main);
  font-size: 0.88rem;
}

.upload-hint {
  color: var(--text-subtle);
  font-size: 0.78rem;
}

.file-list {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 200px;
  overflow-y: auto;
}

.file-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.65rem 0.85rem;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.file-info {
  display: flex;
  flex-direction: column;
}

.file-name {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-main);
}

.file-size {
  font-size: 0.75rem;
  color: var(--text-subtle);
}

.file-actions {
  display: flex;
  gap: 0.35rem;
}

.empty-upload-notice {
  margin-top: 1rem;
  padding: 0.75rem 1rem;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
  border-radius: var(--radius-md);
  font-size: 0.82rem;
  color: #fde68a;
}
</style>
