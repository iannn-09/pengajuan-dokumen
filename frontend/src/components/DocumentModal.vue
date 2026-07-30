<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="close">
    <div class="modal-content">
      <div class="modal-header">
        <h3 class="modal-title">Form Pengajuan Dokumen Baru</h3>
        <button class="btn-close" @click="close">
          <X :size="20" />
        </button>
      </div>

      <form @submit.prevent="handleSubmit">
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Judul Pengajuan *</label>
            <input 
              v-model="form.title" 
              type="text" 
              class="form-input" 
              placeholder="Contoh: Pengajuan Surat Keterangan Aktif Kuliah" 
              required 
            />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Nama Pemohon *</label>
              <input 
                v-model="form.applicant_name" 
                type="text" 
                class="form-input" 
                placeholder="Nama Lengkap" 
                required 
              />
            </div>
            <div class="form-group">
              <label class="form-label">Departemen / Jurusan</label>
              <input 
                v-model="form.department" 
                type="text" 
                class="form-input" 
                placeholder="Contoh: Teknik Informatika" 
              />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Jenis Dokumen *</label>
            <select v-model="form.document_type" class="form-select" required>
              <option value="" disabled>-- Pilih Jenis Dokumen --</option>
              <option value="Surat Keterangan">Surat Keterangan</option>
              <option value="Surat Izin">Surat Izin</option>
              <option value="Legalisis">Legalisir Ijazah / Transkrip</option>
              <option value="Rekomendasi">Surat Rekomendasi</option>
              <option value="Lainnya">Dokumen Lainnya</option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">Keterangan / Alasan Pengajuan</label>
            <textarea 
              v-model="form.description" 
              class="form-textarea" 
              placeholder="Jelaskan detail permohonan atau dokumen pendukung yang dibutuhkan..."
            ></textarea>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="close">Batal</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">
            <Loader2 v-if="submitting" class="spin" :size="18" />
            <Send v-else :size="18" />
            <span>{{ submitting ? 'Mengirim...' : 'Kirim Pengajuan' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { X, Send, Loader2 } from 'lucide-vue-next'
import { DocumentService } from '../services/api'

const props = defineProps({
  isOpen: Boolean
})

const emit = defineEmits(['close', 'created'])

const submitting = ref(false)
const form = reactive({
  title: '',
  applicant_name: '',
  department: '',
  document_type: '',
  description: ''
})

const resetForm = () => {
  form.title = ''
  form.applicant_name = ''
  form.department = ''
  form.document_type = ''
  form.description = ''
}

const close = () => {
  resetForm()
  emit('close')
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const res = await DocumentService.createDocument(form)
    if (res.data && res.data.status === 'success') {
      emit('created', res.data.data)
      close()
    }
  } catch (err) {
    alert('Gagal membuat pengajuan dokumen: ' + (err.response?.data?.error || err.message))
  } finally {
    submitting.value = false
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
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  color: var(--text-main);
  background: rgba(255, 255, 255, 0.1);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  100% { transform: rotate(360deg); }
}

@media (max-width: 640px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
