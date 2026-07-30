<template>
  <div class="timeline-container">
    <h3 class="timeline-title">Riwayat Penilaian & Workflow</h3>
    <div v-if="!histories || histories.length === 0" class="timeline-empty">
      Belum ada riwayat penilaian. Project masih dalam tahap pengajuan awal.
    </div>
    <div v-else class="timeline-list">
      <div v-for="item in histories" :key="item.id" class="timeline-item">
        <div class="timeline-icon-box" :class="`type-${item.status_to.toLowerCase()}`">
          <span class="timeline-dot"></span>
        </div>
        <div class="timeline-content">
          <div class="timeline-header">
            <div class="timeline-user">
              <span class="reviewer-name">{{ item.reviewer?.name || 'Penilai' }}</span>
              <span class="action-text">mengubah status dari</span>
              <StatusBadge :status="item.status_from" />
              <span class="action-text">menjadi</span>
              <StatusBadge :status="item.status_to" />
            </div>
            <span class="timeline-date">{{ formatDate(item.created_at) }}</span>
          </div>
          <div v-if="item.notes" class="timeline-notes">
            <strong>Catatan:</strong> {{ item.notes }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import StatusBadge from './StatusBadge.vue'

defineProps({
  histories: {
    type: Array,
    default: () => []
  }
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.timeline-container {
  margin-top: 1.5rem;
}

.timeline-title {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 1rem;
}

.timeline-empty {
  font-size: 0.85rem;
  color: var(--text-muted);
  font-style: italic;
  padding: 1rem;
  background: var(--bg-input);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
}

.timeline-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  position: relative;
  padding-left: 1rem;
}

.timeline-list::before {
  content: '';
  position: absolute;
  left: 19px;
  top: 10px;
  bottom: 10px;
  width: 2px;
  background: var(--border-color);
}

.timeline-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  position: relative;
  z-index: 1;
}

.timeline-icon-box {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--bg-secondary);
  border: 2px solid var(--accent-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 0.2rem;
}

.timeline-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent-primary);
}

.type-approved { border-color: #10b981; }
.type-approved .timeline-dot { background: #10b981; }

.type-revision { border-color: #f59e0b; }
.type-revision .timeline-dot { background: #f59e0b; }

.type-rejected { border-color: #f43f5e; }
.type-rejected .timeline-dot { background: #f43f5e; }

.timeline-content {
  flex: 1;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 0.85rem 1rem;
}

.timeline-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.timeline-user {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.85rem;
  flex-wrap: wrap;
}

.reviewer-name {
  font-weight: 700;
  color: var(--text-main);
}

.action-text {
  color: var(--text-muted);
}

.timeline-date {
  font-size: 0.75rem;
  color: var(--text-subtle);
}

.timeline-notes {
  margin-top: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px dashed var(--border-color);
  font-size: 0.85rem;
  color: var(--text-main);
  white-space: pre-wrap;
}
</style>
