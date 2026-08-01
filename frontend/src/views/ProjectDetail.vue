<template>
  <div v-if="loading" class="loading-box">Memuat detail permohonan...</div>
  <div v-else-if="!project" class="empty-box">Project tidak ditemukan.</div>
  <div v-else class="detail-container">
    <!-- Header -->
    <div class="glass-card header-card">
      <div class="header-top">
        <button class="btn btn-secondary btn-sm" @click="$router.back()">&larr; Kembali</button>
        <StatusBadge :status="project.status" />
      </div>

      <div class="header-main">
        <span class="project-num">{{ project.project_number }}</span>
        <h1 class="project-title">{{ project.title }}</h1>
        <div class="project-meta">
          <span>Perusahaan: <strong>{{ project.company_name || '-' }}</strong></span>
          <span>&bull;</span>
          <span>Pemohon: <strong>{{ project.user?.name || '-' }}</strong></span>
          <span>&bull;</span>
          <span>Dibuat: <strong>{{ formatDate(project.created_at) }}</strong></span>
        </div>
      </div>

      <!-- Action Banner -->
      <div v-if="auth.isPemohon" class="action-banner">
        <div v-if="project.status === 'DRAFT'" class="banner-box banner-info">
          <span>Project ini masih berstatus <strong>DRAFT</strong>. Pastikan dokumen lampiran sudah diunggah sebelum mengirimkan untuk penilaian.</span>
          <button class="btn btn-primary btn-sm" @click="submitProject">Kirimkan Untuk Penilaian &rarr;</button>
        </div>
        <div v-else-if="project.status === 'REVISION'" class="banner-box banner-warning">
          <span>Penilai meminta <strong>REVISI</strong> pada permohonan ini. Silakan ubah data/dokumen lalu kirim ulang.</span>
          <div class="banner-btns">
            <button class="btn btn-secondary btn-sm" @click="$router.push(`/projects/${project.id}/edit`)">Edit Project</button>
            <button class="btn btn-primary btn-sm" @click="submitProject">Kirim Ulang Permohonan &rarr;</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Details Grid -->
    <div class="grid-2col">
      <!-- Left Column: Description & Documents -->
      <div class="glass-card detail-card">
        <h3 class="section-title">Deskripsi Permohonan</h3>
        <p class="description-text">{{ project.description || 'Tidak ada deskripsi tambahan.' }}</p>

        <!-- Document Type Info -->
        <div v-if="project.document_type" class="doc-type-box">
          <h4 class="doc-type-title">Jenis Dokumen: {{ project.document_type.name }} ({{ project.document_type.code }})</h4>
          <div v-if="project.document_type.requirement" class="req-box">
            <strong>Persyaratan Berkas Wajib:</strong>
            <div class="html-content" v-html="project.document_type.requirement"></div>
          </div>
        </div>

        <h3 class="section-title margin-top">Dokumen Lampiran ({{ documents.length }})</h3>
        <div v-if="documents.length === 0" class="empty-text">Belum ada dokumen lampiran yang diunggah.</div>
        <div v-else class="doc-list">
          <div v-for="doc in documents" :key="doc.id" class="doc-item">
            <div class="doc-meta">
              <span class="doc-name">{{ doc.file_name }}</span>
              <small class="doc-sub">{{ doc.file_type?.toUpperCase() }} &bull; {{ formatSize(doc.file_size) }}</small>
            </div>
            <a :href="getDownloadUrl(doc.id)" target="_blank" class="btn btn-secondary btn-sm">Download</a>
          </div>
        </div>
      </div>

      <!-- Right Column: Review History Timeline -->
      <div class="glass-card detail-card">
        <ReviewTimeline :histories="project.review_histories || []" />
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
const project = ref(null)

const documents = computed(() => project.value?.documents || [])

const fetchProject = async () => {
  loading.value = true
  try {
    const res = await apiClient.get(`/projects/${route.params.id}`)
    if (res.data?.data) {
      project.value = res.data.data
    }
  } catch (err) {
    alertError('Gagal Memuat Detail', err.response?.data?.error || err.message)
    router.push('/projects')
  } finally {
    loading.value = false
  }
}

const submitProject = async () => {
  if (documents.value.length === 0) {
    alertWarning('Dokumen Belum Diunggah', 'Harap unggah minimal 1 berkas lampiran pendukung sebelum mengirimkan permohonan!')
    return
  }

  const confirmed = await confirmDialog(
    'Kirimkan Permohonan?',
    'Kirimkan permohonan dokumen ini kepada Penilai untuk diverifikasi?',
    'Ya, Kirim Sekarang'
  )

  if (confirmed) {
    try {
      await apiClient.post(`/projects/${project.value.id}/submit`)
      alertSuccess('Berhasil Terkirim!', 'Permohonan berhasil dikirimkan untuk proses penilaian.')
      fetchProject()
    } catch (err) {
      alertError('Gagal Mengirimkan', err.response?.data?.error || err.message)
    }
  }
}

const getDownloadUrl = (docId) => {
  return `${apiClient.defaults.baseURL}/documents/${docId}/download`
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
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

.action-banner {
  margin-top: 0.5rem;
}

.banner-box {
  padding: 1rem 1.25rem;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  font-size: 0.88rem;
}

.banner-info {
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  color: #93c5fd;
}

.banner-warning {
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
  color: #fde68a;
}

.banner-btns {
  display: flex;
  gap: 0.5rem;
}

.grid-2col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.detail-card {
  padding: 1.5rem;
}

.section-title {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 0.5rem;
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
  margin-top: 1rem;
  padding: 1rem;
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: var(--radius-md);
}

.doc-type-title {
  font-size: 0.95rem;
  font-weight: 800;
  color: var(--text-main);
}

.req-box {
  margin-top: 0.5rem;
  font-size: 0.83rem;
  color: var(--text-muted);
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

@media (max-width: 1024px) {
  .grid-2col { grid-template-columns: 1fr; }
}
</style>
