<template>
  <div v-if="isOpen && doc" class="modal-backdrop" @click.self="close">
    <div class="modal-content modal-xl">
      <div class="modal-header">
        <div class="header-info">
          <h3 class="modal-title">Pratinjau Dokumen: {{ doc.file_name }}</h3>
          <span class="file-meta-tag">{{ doc.file_type?.toUpperCase() }} &bull; {{ formatSize(doc.file_size) }}</span>
        </div>
        <div class="header-actions">
          <a :href="getDownloadUrl(doc.id)" target="_blank" class="btn btn-secondary btn-sm">Download File</a>
          <button class="btn-close" @click="close">&times;</button>
        </div>
      </div>

      <div class="modal-body preview-body">
        <!-- Image Viewer -->
        <div v-if="isImage" class="image-preview-container">
          <img :src="getDownloadUrl(doc.id)" :alt="doc.file_name" class="preview-img" />
        </div>

        <!-- PDF Viewer (iframe) -->
        <div v-else-if="isPdf" class="pdf-preview-container">
          <iframe 
            :src="getDownloadUrl(doc.id)" 
            class="preview-iframe"
            title="PDF Preview"
          ></iframe>
        </div>

        <!-- Fallback Viewer (DOC / DOCX / Other) -->
        <div v-else class="fallback-preview">
          <div class="fallback-icon">📄</div>
          <h4 class="fallback-title">Pratinjau Tidak Tersedia Langsung</h4>
          <p class="fallback-desc">File berformat <strong>.{{ doc.file_type }}</strong> tidak mendukung pratinjau browser secara langsung.</p>
          <a :href="getDownloadUrl(doc.id)" target="_blank" class="btn btn-primary">
            Download & Buka File ↗
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import apiClient from '../services/api'

const props = defineProps({
  isOpen: { type: Boolean, default: false },
  doc: { type: Object, default: null }
})

const emit = defineEmits(['close'])

const close = () => emit('close')

const isImage = computed(() => {
  if (!props.doc?.file_type) return false
  const type = props.doc.file_type.toLowerCase()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(type)
})

const isPdf = computed(() => {
  if (!props.doc?.file_type) return false
  return props.doc.file_type.toLowerCase() === 'pdf'
})

const getDownloadUrl = (docId) => {
  if (!docId) return ''
  return `${apiClient.defaults.baseURL}/documents/${docId}/download`
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
.modal-xl {
  max-width: 900px;
  width: 90%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.header-info {
  display: flex;
  flex-direction: column;
}

.file-meta-tag {
  font-size: 0.78rem;
  color: var(--accent-primary);
  font-weight: 700;
  margin-top: 0.15rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.btn-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
}

.preview-body {
  padding: 1.25rem;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 450px;
  max-height: 75vh;
  overflow: hidden;
  background: rgba(10, 15, 30, 0.6);
}

.image-preview-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;
}

.preview-img {
  max-width: 100%;
  max-height: 65vh;
  object-fit: contain;
  border-radius: var(--radius-md);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
}

.pdf-preview-container {
  width: 100%;
  height: 65vh;
}

.preview-iframe {
  width: 100%;
  height: 100%;
  border: none;
  border-radius: var(--radius-md);
  background: white;
}

.fallback-preview {
  text-align: center;
  padding: 3rem 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.fallback-icon {
  font-size: 3.5rem;
}

.fallback-title {
  font-size: 1.2rem;
  font-weight: 800;
  color: var(--text-main);
}

.fallback-desc {
  font-size: 0.88rem;
  color: var(--text-muted);
  max-width: 420px;
}
</style>
