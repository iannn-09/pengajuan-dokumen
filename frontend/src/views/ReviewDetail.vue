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

        <h3 class="section-title margin-top">Dokumen Lampiran ({{ documents.length }})</h3>
        <div v-if="documents.length === 0" class="empty-text">Pemohon belum mengunggah dokumen lampiran.</div>
        <div v-else class="doc-list">
          <div v-for="doc in documents" :key="doc.id" class="doc-item">
            <div class="doc-meta">
              <span class="doc-name">{{ doc.file_name }}</span>
              <small class="doc-sub">{{ doc.file_type?.toUpperCase() }} &bull; {{ formatSize(doc.file_size) }}</small>
            </div>
            <a :href="getDownloadUrl(doc.id)" target="_blank" class="btn btn-primary btn-sm">Download File</a>
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
            <textarea 
              v-model="notes" 
              class="form-textarea" 
              placeholder="Tuliskan alasan persetujuan, daftar dokumen yang perlu direvisi, atau alasan penolakan..."
              rows="6"
              :required="action !== 'approve'"
            ></textarea>
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
import apiClient from '../services/api'
import StatusBadge from '../components/StatusBadge.vue'
import ReviewTimeline from '../components/ReviewTimeline.vue'

const route = useRoute()
const router = useRouter()

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
    alert('Gagal memuat detail permohonan: ' + (err.response?.data?.error || err.message))
    router.push('/reviews')
  } finally {
    loading.value = false
  }
}

const handleAssess = async () => {
  if (action.value !== 'approve' && !notes.value.trim()) {
    alert('Harap isi catatan penilai untuk memberikan revisi atau penolakan!')
    return
  }

  if (confirm(`Simpan keputusan penilaian (${action.value.toUpperCase()}) untuk permohonan ini?`)) {
    submitting.value = true
    try {
      await apiClient.post(`/reviews/projects/${project.value.id}/assess`, {
        action: action.value,
        notes: notes.value
      })
      alert('Penilaian berhasil disimpan!')
      router.push('/reviews')
    } catch (err) {
      alert('Gagal menyimpan penilaian: ' + (err.response?.data?.error || err.message))
    } finally {
      submitting.value = false
    }
  }
}

const getDownloadUrl = (docId) => {
  return `${apiClient.defaults.baseURL}/documents/${docId}/download`
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

.loading-box, .empty-box {
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

.detail-card, .assessment-panel {
  padding: 1.5rem;
}

.section-title, .panel-title {
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

.action-option input { display: none; }

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
  .grid-2col { grid-template-columns: 1fr; }
}
</style>
