<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="close">
    <div class="modal-content">
      <div class="modal-header">
        <h3 class="modal-title">Ubah Status Pengajuan</h3>
        <button class="btn-close" @click="close">
          <X :size="20" />
        </button>
      </div>

      <div class="modal-body">
        <div v-if="document" class="doc-summary">
          <div class="doc-num">{{ document.document_number }}</div>
          <div class="doc-title">{{ document.title }}</div>
          <div class="doc-applicant">Pemohon: {{ document.applicant_name }} ({{ document.department || 'Umum' }})</div>
        </div>

        <div class="form-group">
          <label class="form-label">Pilih Status Baru *</label>
          <div class="status-options">
            <label class="status-option option-approved" :class="{ selected: selectedStatus === 'APPROVED' }">
              <input type="radio" v-model="selectedStatus" value="APPROVED" />
              <CheckCircle2 :size="20" />
              <span>Setujui (APPROVED)</span>
            </label>
            <label class="status-option option-rejected" :class="{ selected: selectedStatus === 'REJECTED' }">
              <input type="radio" v-model="selectedStatus" value="REJECTED" />
              <XCircle :size="20" />
              <span>Tolak (REJECTED)</span>
            </label>
            <label class="status-option option-pending" :class="{ selected: selectedStatus === 'PENDING' }">
              <input type="radio" v-model="selectedStatus" value="PENDING" />
              <Clock :size="20" />
              <span>Menunggu (PENDING)</span>
            </label>
          </div>
        </div>

        <div v-if="selectedStatus === 'REJECTED'" class="form-group animate-fade">
          <label class="form-label">Alasan Penolakan *</label>
          <textarea 
            v-model="rejectionReason" 
            class="form-textarea" 
            placeholder="Jelaskan alasan mengapa pengajuan ini ditolak..."
            required
          ></textarea>
        </div>
      </div>

      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" @click="close">Batal</button>
        <button 
          type="button" 
          class="btn" 
          :class="buttonClass" 
          :disabled="updating"
          @click="handleSave"
        >
          <Loader2 v-if="updating" class="spin" :size="18" />
          <span>{{ updating ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { X, CheckCircle2, XCircle, Clock, Loader2 } from 'lucide-vue-next'
import { DocumentService } from '../services/api'

const props = defineProps({
  isOpen: Boolean,
  document: Object
})

const emit = defineEmits(['close', 'updated'])

const selectedStatus = ref('APPROVED')
const rejectionReason = ref('')
const updating = ref(false)

watch(() => props.document, (newDoc) => {
  if (newDoc) {
    selectedStatus.value = newDoc.status || 'APPROVED'
    rejectionReason.value = newDoc.rejection_reason || ''
  }
}, { immediate: true })

const close = () => {
  emit('close')
}

const buttonClass = computed(() => {
  if (selectedStatus.value === 'APPROVED') return 'btn-success'
  if (selectedStatus.value === 'REJECTED') return 'btn-danger'
  return 'btn-primary'
})

const handleSave = async () => {
  if (selectedStatus.value === 'REJECTED' && !rejectionReason.value.trim()) {
    alert('Harap isi alasan penolakan!')
    return
  }

  updating.value = true
  try {
    const res = await DocumentService.updateStatus(props.document.id, {
      status: selectedStatus.value,
      rejection_reason: selectedStatus.value === 'REJECTED' ? rejectionReason.value : ''
    })

    if (res.data && res.data.status === 'success') {
      emit('updated', res.data.data)
      close()
    }
  } catch (err) {
    alert('Gagal mengupdate status: ' + (err.response?.data?.error || err.message))
  } finally {
    updating.value = false
  }
}
</script>

<style scoped>
.btn-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 0.25rem;
  border-radius: var(--radius-sm);
}

.doc-summary {
  background: rgba(15, 23, 42, 0.5);
  padding: 0.85rem 1rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  margin-bottom: 1.25rem;
}

.doc-num {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--accent-primary);
}

.doc-title {
  font-weight: 700;
  color: var(--text-main);
  margin: 0.15rem 0;
}

.doc-applicant {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.status-options {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.status-option {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  background: var(--bg-input);
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9rem;
  font-weight: 600;
}

.status-option input {
  display: none;
}

.option-approved.selected {
  background: rgba(16, 185, 129, 0.15);
  border-color: rgba(16, 185, 129, 0.5);
  color: #34d399;
}

.option-rejected.selected {
  background: rgba(244, 63, 94, 0.15);
  border-color: rgba(244, 63, 94, 0.5);
  color: #fb7185;
}

.option-pending.selected {
  background: rgba(245, 158, 11, 0.15);
  border-color: rgba(245, 158, 11, 0.5);
  color: #fbbf24;
}

.animate-fade {
  animation: fadeIn 0.2s ease;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  100% { transform: rotate(360deg); }
}
</style>
