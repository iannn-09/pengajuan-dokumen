<template>
  <div v-if="totalPages > 1" class="pagination-wrapper">
    <div class="pagination-info">
      Menampilkan Halaman {{ page }} dari {{ totalPages }} (Total {{ total }} Data)
    </div>
    <div class="pagination-controls">
      <button 
        class="page-btn" 
        :disabled="page <= 1" 
        @click="$emit('change-page', page - 1)"
      >
        &laquo; Prev
      </button>
      
      <button 
        v-for="p in visiblePages" 
        :key="p" 
        class="page-btn" 
        :class="{ active: p === page }" 
        @click="$emit('change-page', p)"
      >
        {{ p }}
      </button>

      <button 
        class="page-btn" 
        :disabled="page >= totalPages" 
        @click="$emit('change-page', page + 1)"
      >
        Next &raquo;
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  page: { type: Number, default: 1 },
  totalPages: { type: Number, default: 1 },
  total: { type: Number, default: 0 }
})

defineEmits(['change-page'])

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, props.page - 2)
  const end = Math.min(props.totalPages, props.page + 2)
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})
</script>

<style scoped>
.pagination-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 0;
  border-top: 1px solid var(--border-color);
  margin-top: 1.5rem;
}

.pagination-info {
  font-size: 0.82rem;
  color: var(--text-muted);
}

.pagination-controls {
  display: flex;
  gap: 0.35rem;
}

.page-btn {
  padding: 0.35rem 0.75rem;
  font-size: 0.82rem;
  font-weight: 600;
  background: var(--bg-input);
  color: var(--text-main);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
}

.page-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
}

.page-btn.active {
  background: var(--accent-primary);
  color: #fff;
  border-color: var(--accent-primary);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
