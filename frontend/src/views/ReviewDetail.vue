<template>
  <div v-if="loading" class="loading-box">Memuat data penilaian...</div>
  <div v-else-if="!project" class="empty-box">Project tidak ditemukan.</div>
  <div v-else class="detail-container">
    <!-- Header -->
    <div class="glass-card header-card">
      <div class="header-top">
        <button class="btn btn-secondary btn-sm" @click="$router.back()">&larr; Kembali ke Daftar</button>
        <StatusBadge :status="project.status" />
      </div>

      <div class="header-main">
        <span class="project-num">{{ project.project_number }}</span>
        <h1 class="project-title">{{ project.title }}</h1>
        <div class="project-meta">
          <span>Pemohon: <strong>{{ project.user?.name || '-' }}</strong></span>
          <span>&bull;</span>
          <span>Perusahaan: <strong>{{ project.company_name || '-' }}</strong></span>
          <span>&bull;</span>
          <span v-if="project.unit">Unit Kerja: <strong>{{ project.unit }}</strong></span>
          <span v-if="project.unit">&bull;</span>
          <span>Tanggal Submit: <strong>{{ formatDate(project.submitted_at || project.created_at) }}</strong></span>
        </div>
      </div>
    </div>

    <!-- Main Grid -->
    <div class="grid-2col">
      <!-- Left Column: Project Details & Documents -->
      <div class="glass-card detail-card">
        <h3 class="section-title">Detail & Deskripsi Permohonan</h3>
        <p class="description-text">{{ project.description || 'Tidak ada deskripsi tambahan.' }}</p>

        <!-- Document Type Info -->
        <div v-if="project.document_type" class="doc-type-box">
          <div class="doc-type-header">
            <h4 class="doc-type-title">Jenis Dokumen: {{ project.document_type.name }}</h4>
            <span class="doc-type-code">{{ project.document_type.code }}</span>
          </div>
          <div v-if="project.document_type.description" class="doc-type-desc html-content"
            v-html="project.document_type.description"></div>

          <div v-if="project.document_type.requirement" class="req-box">
            <strong class="req-title">📋 Berkas Persyaratan Wajib:</strong>
            <div class="req-text html-content" v-html="project.document_type.requirement"></div>
          </div>
        </div>

        <h3 class="section-title margin-top">Dokumen Lampiran ({{ documents.length }})</h3>
        <div v-if="documents.length === 0" class="empty-text">Pemohon belum mengunggah dokumen lampiran.</div>
        <div v-else class="doc-list">
          <div v-for="doc in documents" :key="doc.id" class="doc-item">
            <div class="doc-meta">
              <span class="doc-name">{{ doc.file_name }}</span>
              <small class="doc-sub">{{ doc.file_type?.toUpperCase() }} &bull; {{ formatSize(doc.file_size) }}</small>
            </div>
            <div class="doc-actions">
              <a :href="getDownloadUrl(doc.id)" target="_blank" class="btn btn-primary btn-sm">Lihat</a>
            </div>
          </div>
        </div>

        <ReviewTimeline :histories="project.review_histories || []" class="margin-top" />
      </div>

      <!-- Right Column: Assessment Form Panel -->
      <div class="glass-card assessment-panel">
        <h3 class="panel-title">Form Keputusan & Penilaian</h3>
        <p class="panel-subtitle">Pilih tindakan penilaian dan berikan catatan evaluasi untuk pemohon.</p>

        <form @submit.prevent="handleAssess" class="assessment-form">
          <div class="form-group">
            <label class="form-label">Keputusan Penilaian *</label>
            <div class="action-selector">
              <label class="action-option option-approve" :class="{ selected: action === 'approve' }">
                <input type="radio" v-model="action" value="approve" />
                <span class="action-icon">✓</span>
                <span>Setujui (APPROVED)</span>
              </label>

              <label class="action-option option-revise" :class="{ selected: action === 'revise' }">
                <input type="radio" v-model="action" value="revise" />
                <span class="action-icon">↺</span>
                <span>Minta Revisi (REVISION)</span>
              </label>

              <label class="action-option option-reject" :class="{ selected: action === 'reject' }">
                <input type="radio" v-model="action" value="reject" />
                <span class="action-icon">✕</span>
                <span>Tolak (REJECTED)</span>
              </label>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Catatan Penilai / Evaluasi *</label>
            <textarea v-model="notes" class="form-textarea"
              placeholder="Tuliskan alasan persetujuan, daftar dokumen yang perlu direvisi, atau alasan penolakan..."
              rows="6" :required="action !== 'approve'"></textarea>
          </div>

          <button type="submit" class="btn btn-primary btn-block" :disabled="submitting">
            <span>{{ submitting ? 'Memproses Keputusan...' : 'Simpan Keputusan Penilaian' }}</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import StatusBadge from '../components/StatusBadge.vue'
import ReviewTimeline from '../components/ReviewTimeline.vue'
import { alertSuccess, alertError, alertWarning, confirmDialog } from '../utils/swal'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const submitting = ref(false)
const project = ref(null)

const action = ref('approve')
const notes = ref('')

const documents = computed(() => project.value?.documents || [])

const fetchProject = async () => {
  loading.value = true
  try {
    const res = await apiClient.get(`/reviews/projects/${route.params.id}`)
    if (res.data?.data) {
      project.value = res.data.data
    }
  } catch (err) {
    alertError('Gagal Memuat Detail', err.response?.data?.error || err.message)
    router.push('/reviews')
  } finally {
    loading.value = false
  }
}

const handleAssess = async () => {
  if (action.value !== 'approve' && !notes.value.trim()) {
    alertWarning('Catatan Wajib Diisi', 'Harap isi catatan penilai untuk memberikan instruksi revisi atau alasan penolakan!')
    return
  }

  const confirmed = await confirmDialog(
    'Simpan Keputusan Penilaian?',
    `Simpan keputusan penilaian (${action.value.toUpperCase()}) untuk permohonan ini?`,
    'Ya, Simpan Keputusan'
  )

  if (confirmed) {
    submitting.value = true
    try {
      await apiClient.post(`/reviews/projects/${project.value.id}/assess`, {
        action: action.value,
        notes: notes.value
      })
      alertSuccess('Penilaian Berhasil!', 'Keputusan penilaian telah berhasil disimpan.')
      router.push('/reviews')
    } catch (err) {
      alertError('Gagal Menyimpan Penilaian', err.response?.data?.error || err.message)
    } finally {
      submitting.value = false
    }
  }
}

const getDownloadUrl = (docId) => {
  return `${apiClient.defaults.baseURL}/documents/${docId}/download?token=${auth.token}`
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(() => {
  fetchProject()
})
</script>

<style scoped>
.detail-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.loading-box,
.empty-box {
  padding: 4rem;
  text-align: center;
  color: var(--text-muted);
}

.header-card {
  padding: 1.75rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.project-num {
  font-family: monospace;
  font-size: 0.85rem;
  color: var(--accent-primary);
  font-weight: 700;
}

.project-title {
  font-size: 1.4rem;
  font-weight: 800;
  color: var(--text-main);
  margin-top: 0.2rem;
}

.project-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.6rem;
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-top: 0.4rem;
}

.grid-2col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.detail-card,
.assessment-panel {
  padding: 1.5rem;
}

.section-title,
.panel-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-main);
}

.panel-subtitle {
  font-size: 0.83rem;
  color: var(--text-muted);
  margin-bottom: 1.25rem;
}

.margin-top {
  margin-top: 1.5rem;
}

.description-text {
  font-size: 0.9rem;
  color: var(--text-muted);
  line-height: 1.6;
  white-space: pre-wrap;
  background: var(--bg-input);
  padding: 1rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
}

.doc-type-box {
  margin-top: 1.25rem;
  padding: 1.25rem;
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: var(--radius-md);
}

.doc-type-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.doc-type-title {
  font-size: 0.98rem;
  font-weight: 800;
  color: var(--text-main);
}

.doc-type-code {
  font-family: monospace;
  font-size: 0.85rem;
  font-weight: 800;
  color: var(--accent-primary);
  background: rgba(99, 102, 241, 0.15);
  padding: 0.2rem 0.65rem;
  border-radius: 9999px;
}

.doc-type-desc {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-bottom: 0.75rem;
}

.req-box {
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px dashed rgba(99, 102, 241, 0.3);
  font-size: 0.88rem;
  color: var(--text-main);
}

.req-title {
  font-weight: 700;
  display: block;
  margin-bottom: 0.35rem;
}

.req-text {
  margin-top: 0.35rem;
  background: rgba(0, 0, 0, 0.25);
  padding: 0.75rem 1rem;
  border-radius: var(--radius-sm);
}

.html-content :deep(p) {
  margin-bottom: 0.35rem;
}

.html-content :deep(ul),
.html-content :deep(ol) {
  padding-left: 1.5rem;
  margin-top: 0.35rem;
  margin-bottom: 0.35rem;
}

.html-content :deep(li) {
  margin-bottom: 0.25rem;
  color: var(--text-muted);
  line-height: 1.5;
}

.doc-list {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
  margin-top: 0.5rem;
}

.doc-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.doc-meta {
  display: flex;
  flex-direction: column;
}

.doc-name {
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--text-main);
}

.doc-sub {
  font-size: 0.75rem;
  color: var(--text-subtle);
}

.doc-actions {
  display: flex;
  gap: 0.35rem;
}

.empty-text {
  font-size: 0.85rem;
  color: var(--text-subtle);
  font-style: italic;
}

.action-selector {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.action-option {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  background: var(--bg-input);
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 600;
  transition: all 0.2s ease;
}

.action-option input {
  display: none;
}

.option-approve.selected {
  background: rgba(16, 185, 129, 0.15);
  border-color: rgba(16, 185, 129, 0.5);
  color: #34d399;
}

.option-revise.selected {
  background: rgba(245, 158, 11, 0.15);
  border-color: rgba(245, 158, 11, 0.5);
  color: #fbbf24;
}

.option-reject.selected {
  background: rgba(244, 63, 94, 0.15);
  border-color: rgba(244, 63, 94, 0.5);
  color: #fb7185;
}

.action-icon {
  font-size: 1.1rem;
  font-weight: 800;
}

.btn-block {
  width: 100%;
  padding: 0.8rem;
  font-size: 0.95rem;
  margin-top: 1rem;
}

@media (max-width: 1024px) {
  .grid-2col {
    grid-template-columns: 1fr;
  }
}
</style>
