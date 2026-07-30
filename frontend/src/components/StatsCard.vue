<template>
  <div class="glass-card stat-card" :class="`type-${type}`" @click="$emit('select')">
    <div class="stat-header">
      <span class="stat-title">{{ title }}</span>
      <div class="stat-icon-wrapper">
        <component :is="icon" :size="20" />
      </div>
    </div>
    <div class="stat-body">
      <span class="stat-value">{{ value }}</span>
      <span class="stat-desc">Dokumen</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: String,
  value: [Number, String],
  icon: Object,
  type: {
    type: String,
    default: 'default'
  }
})

defineEmits(['select'])
</script>

<style scoped>
.stat-card {
  padding: 1.25rem 1.5rem;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: var(--accent-primary);
  border-radius: 4px 0 0 4px;
  opacity: 0.8;
}

.type-pending::before { background: #f59e0b; }
.type-approved::before { background: #10b981; }
.type-rejected::before { background: #f43f5e; }
.type-total::before { background: #6366f1; }

.stat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.stat-title {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-muted);
}

.stat-icon-wrapper {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-main);
}

.type-pending .stat-icon-wrapper { color: #fbbf24; background: rgba(245, 158, 11, 0.12); }
.type-approved .stat-icon-wrapper { color: #34d399; background: rgba(16, 185, 129, 0.12); }
.type-rejected .stat-icon-wrapper { color: #fb7185; background: rgba(244, 63, 94, 0.12); }
.type-total .stat-icon-wrapper { color: #818cf8; background: rgba(99, 102, 241, 0.12); }

.stat-body {
  display: flex;
  align-items: baseline;
  gap: 0.4rem;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 800;
  color: var(--text-main);
  line-height: 1;
}

.stat-desc {
  font-size: 0.8rem;
  color: var(--text-subtle);
}
</style>
