<template>
  <span class="status-badge" :class="badgeClass">
    <span class="status-dot"></span>
    <span>{{ statusText }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
    default: 'DRAFT'
  }
})

const badgeClass = computed(() => {
  switch (props.status) {
    case 'DRAFT': return 'badge-draft'
    case 'SUBMITTED': return 'badge-submitted'
    case 'UNDER_REVIEW': return 'badge-review'
    case 'REVISION': return 'badge-revision'
    case 'APPROVED': return 'badge-approved'
    case 'REJECTED': return 'badge-rejected'
    default: return 'badge-draft'
  }
})

const statusText = computed(() => {
  switch (props.status) {
    case 'DRAFT': return 'Draft'
    case 'SUBMITTED': return 'Telah Dikirim'
    case 'UNDER_REVIEW': return 'Sedang Dinilai'
    case 'REVISION': return 'Perlu Revisi'
    case 'APPROVED': return 'Disetujui'
    case 'REJECTED': return 'Ditolak'
    default: return props.status
  }
})
</script>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.3rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 600;
  line-height: 1;
  border: 1px solid transparent;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.badge-draft {
  background: rgba(148, 163, 184, 0.12);
  color: #94a3b8;
  border-color: rgba(148, 163, 184, 0.3);
}

.badge-submitted {
  background: rgba(59, 130, 246, 0.12);
  color: #60a5fa;
  border-color: rgba(59, 130, 246, 0.3);
}

.badge-review {
  background: rgba(168, 85, 247, 0.12);
  color: #c084fc;
  border-color: rgba(168, 85, 247, 0.3);
}

.badge-revision {
  background: rgba(245, 158, 11, 0.12);
  color: #fbbf24;
  border-color: rgba(245, 158, 11, 0.3);
}

.badge-approved {
  background: rgba(16, 185, 129, 0.12);
  color: #34d399;
  border-color: rgba(16, 185, 129, 0.3);
}

.badge-rejected {
  background: rgba(244, 63, 94, 0.12);
  color: #fb7185;
  border-color: rgba(244, 63, 94, 0.3);
}
</style>
