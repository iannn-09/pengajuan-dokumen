<template>
  <div>
    <DashboardPemohon v-if="auth.isPemohon" />
    <DashboardPenilai v-else-if="auth.isPenilai" />
    <div v-else class="loading-box">Loading Dashboard...</div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
import DashboardPemohon from './DashboardPemohon.vue'
import DashboardPenilai from './DashboardPenilai.vue'

const auth = useAuthStore()
</script>

<style scoped>
.loading-box {
  padding: 3rem;
  text-align: center;
  color: var(--text-muted);
}
</style>
