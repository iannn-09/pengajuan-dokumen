<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="close">
    <div class="modal-content">
      <div class="modal-header">
        <h3 class="modal-title">Detail Pengajuan Dokumen</h3>
        <button class="btn-close" @click="close">
          <X :size="20" />
        </button>
      </div>

      <div class="modal-body" v-if="document">
        <div class="detail-top">
          <span class="detail-number">{{ document.document_number }}</span>
          <StatusBadge :status="document.status" />
        </div>

        <h2 class="detail-title">{{ document.title }}</h2>

        <div class="detail-grid">
          <div class="detail-item">
            <span class="label">Nama Pemohon</span>
            <span class="value">{{ document.applicant_name }}</span>
          </div>
          <div class="detail-item">
            <span class="label">Departemen / Jurusan</span>
            <span class="value">{{ document.department || '-' }}</span>
          </div>
          <div class="detail-item">
            <span class="label">Jenis Dokumen</span>
            <span class="value">{{ document.document_type }}</span>
          </div>
          <div class="detail-item">
            <span class="label">Tanggal Pengajuan</span>
            <span class="value">{{ formatDate(document.created_at) }}</span>
          </div>
        </div>

        <div class="detail-section">
          <h4 class="section-title">Keterangan / Detail Permohonan</h4>
          <p class="section-text">{{ document.description || 'Tidak ada keterangan tambahan.' }}</p>
        </div>

        <div v-if="document.status === 'REJECTED'" class="detail-section rejection-box">
          <h4 class="section-title text-danger">Alasan Penolakan</h4>
          <p class="section-text">{{ document.rejection_reason || '-' }}</p>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn btn-secondary" @click="close">Tutup</button>
        <button class="btn btn-primary" @click="$emit('change-status', document)">
          <Edit3 :size="16" />
          <span>Ubah Status</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { X, Edit3 } from 'lucide-vue-next'
import StatusBadge from './StatusBadge.vue'

defineProps({
  isOpen: Boolean,
  document: Object
})

const emit = defineEmits(['close', 'change-status'])

const close = () => emit('close')

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.btn-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
}

.detail-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.detail-number {
  font-family: monospace;
  font-size: 0.85rem;
  color: var(--accent-primary);
  font-weight: 700;
}

.detail-title {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 1.25rem;
  line-height: 1.3;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  background: rgba(15, 23, 42, 0.4);
  padding: 1rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  margin-bottom: 1.25rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
}

.label {
  font-size: 0.75rem;
  color: var(--text-muted);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.value {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-main);
  margin-top: 0.2rem;
}

.detail-section {
  margin-bottom: 1.25rem;
}

.section-title {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-muted);
  margin-bottom: 0.4rem;
}

.section-text {
  font-size: 0.92rem;
  color: var(--text-main);
  background: var(--bg-input);
  padding: 0.85rem 1rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  white-space: pre-wrap;
}

.rejection-box {
  background: rgba(244, 63, 94, 0.1);
  border-radius: var(--radius-md);
  padding: 0.85rem 1rem;
  border: 1px solid rgba(244, 63, 94, 0.3);
}

.text-danger {
  color: #fb7185;
}
</style>
