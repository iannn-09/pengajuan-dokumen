<template>
  <div class="dashboard-container">
    <!-- Header banner -->
    <div class="glass-card welcome-banner">
      <div class="welcome-text">
        <h2>Selamat Datang, {{ auth.userName }} 👋</h2>
        <p>Kelola permohonan dokumen kelayakan Anda di portal pelayanan digital.</p>
      </div>
      <router-link to="/projects/create" class="btn btn-primary">
        + Buat Pengajuan Baru
      </router-link>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="glass-card stat-card border-total">
        <span class="stat-label">Total Permohonan Saya</span>
        <span class="stat-number">{{ stats.total }}</span>
      </div>
      <div class="glass-card stat-card border-draft">
        <span class="stat-label">Draft (Belum Dikirim)</span>
        <span class="stat-number">{{ stats.draft }}</span>
      </div>
      <div class="glass-card stat-card border-submitted">
        <span class="stat-label">Telah Dikirim / Dinilai</span>
        <span class="stat-number">{{ (stats.submitted || 0) + (stats.under_review || 0) }}</span>
      </div>
      <div class="glass-card stat-card border-revision">
        <span class="stat-label">Perlu Revisi</span>
        <span class="stat-number text-amber">{{ stats.revision }}</span>
      </div>
      <div class="glass-card stat-card border-approved">
        <span class="stat-label">Disetujui</span>
        <span class="stat-number text-emerald">{{ stats.approved }}</span>
      </div>
      <div class="glass-card stat-card border-rejected">
        <span class="stat-label">Ditolak</span>
        <span class="stat-number text-rose">{{ stats.rejected }}</span>
      </div>
    </div>

    <!-- Chart & Recent Activity -->
    <div class="grid-2col">
      <div class="glass-card chart-card">
        <h3 class="card-title">Grafik Pengajuan Bulanan</h3>
        <div class="chart-container">
          <Bar v-if="chartData.loaded" :data="barChartData" :options="chartOptions" />
          <div v-else class="chart-loading">Memuat Grafik...</div>
        </div>
      </div>

      <div class="glass-card recent-card">
        <div class="card-header-flex">
          <h3 class="card-title">Project Terbaru Saya</h3>
          <router-link to="/projects" class="link-more">Lihat Semua &rarr;</router-link>
        </div>

        <div v-if="recentProjects.length === 0" class="empty-text">
          Belum ada permohonan yang dibuat.
        </div>
        <div v-else class="recent-list">
          <div v-for="item in recentProjects" :key="item.id" class="recent-item" @click="$router.push(`/projects/${item.id}`)">
            <div class="recent-main">
              <span class="recent-num">{{ item.project_number }}</span>
              <h4 class="recent-title">{{ item.title }}</h4>
            </div>
            <StatusBadge :status="item.status" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import StatusBadge from '../components/StatusBadge.vue'

import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'
import { Bar } from 'vue-chartjs'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

const auth = useAuthStore()

const stats = reactive({
  total: 0,
  draft: 0,
  submitted: 0,
  under_review: 0,
  revision: 0,
  approved: 0,
  rejected: 0
})

const recentProjects = ref([])
const chartData = reactive({
  loaded: false,
  labels: [],
  values: []
})

const barChartData = computed(() => ({
  labels: chartData.labels,
  datasets: [
    {
      label: 'Jumlah Pengajuan',
      backgroundColor: '#6366f1',
      borderColor: '#4f46e5',
      borderRadius: 6,
      data: chartData.values
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  },
  scales: {
    x: { grid: { display: false }, ticks: { color: '#94a3b8' } },
    y: { grid: { color: 'rgba(255,255,255,0.05)' }, ticks: { color: '#94a3b8' } }
  }
}

const fetchData = async () => {
  try {
    const resStats = await apiClient.get('/dashboard/stats')
    if (resStats.data?.data) {
      Object.assign(stats, resStats.data.data)
    }

    const resProjects = await apiClient.get('/projects?per_page=5')
    if (resProjects.data?.data) {
      recentProjects.value = resProjects.data.data
    }

    const resChart = await apiClient.get('/dashboard/chart-data')
    if (resChart.data?.data?.monthly_submissions) {
      const subs = resChart.data.data.monthly_submissions
      chartData.labels = subs.map(s => s.month)
      chartData.values = subs.map(s => s.count)
      chartData.loaded = true
    }
  } catch (err) {
    console.error('Failed to load dashboard data:', err)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.welcome-banner {
  padding: 1.75rem 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, rgba(30, 41, 59, 0.7) 100%);
  border: 1px solid rgba(99, 102, 241, 0.3);
}

.welcome-text h2 {
  font-size: 1.4rem;
  font-weight: 800;
  color: var(--text-main);
}

.welcome-text p {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin-top: 0.2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 1rem;
}

.stat-card {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  position: relative;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; bottom: 0;
  width: 4px;
  border-radius: 4px 0 0 4px;
}

.border-total::before { background: #6366f1; }
.border-draft::before { background: #94a3b8; }
.border-submitted::before { background: #3b82f6; }
.border-revision::before { background: #f59e0b; }
.border-approved::before { background: #10b981; }
.border-rejected::before { background: #f43f5e; }

.stat-label {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--text-muted);
}

.stat-number {
  font-size: 1.75rem;
  font-weight: 800;
  color: var(--text-main);
}

.text-amber { color: #fbbf24; }
.text-emerald { color: #34d399; }
.text-rose { color: #fb7185; }

.grid-2col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.chart-card, .recent-card {
  padding: 1.5rem;
  min-height: 340px;
}

.card-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 1rem;
}

.card-header-flex {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.link-more {
  font-size: 0.82rem;
  color: var(--accent-primary);
  text-decoration: none;
  font-weight: 600;
}

.chart-container {
  height: 250px;
  position: relative;
}

.chart-loading, .empty-text {
  padding: 3rem;
  text-align: center;
  color: var(--text-subtle);
  font-size: 0.9rem;
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.recent-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.85rem 1rem;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.recent-item:hover {
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateX(3px);
}

.recent-num {
  font-family: monospace;
  font-size: 0.75rem;
  color: var(--accent-primary);
  font-weight: 700;
}

.recent-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--text-main);
  margin-top: 0.1rem;
}

@media (max-width: 1280px) {
  .stats-grid { grid-template-columns: repeat(3, 1fr); }
  .grid-2col { grid-template-columns: 1fr; }
}

@media (max-width: 640px) {
  .stats-grid { grid-template-columns: 1fr 1fr; }
}
</style>
