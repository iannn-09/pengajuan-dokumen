<template>
  <div class="dashboard-container">
    <!-- Header Welcome Banner -->
    <div class="glass-card welcome-banner" :class="bannerRoleClass">
      <div class="welcome-text">
        <h2>{{ bannerTitle }}</h2>
        <p>{{ bannerSubtitle }}</p>
      </div>
      <div class="banner-actions">
        <button 
          class="btn btn-secondary" 
          @click="handleExportRekapan" 
          :disabled="exporting" 
          title="Unduh Rekapan Laporan Format Excel (CSV)"
        >
          <span>{{ exporting ? 'Mengunduh...' : '📥 Export Rekapan Laporan' }}</span>
        </button>

        <router-link v-if="auth.isAdmin" to="/users" class="btn btn-primary">
          + Kelola User & Penilai
        </router-link>
        <router-link v-else-if="auth.isPenilai" to="/reviews" class="btn btn-primary">
          Mulai Penilaian &rarr;
        </router-link>
        <router-link v-else to="/projects/create" class="btn btn-primary">
          + Buat Pengajuan Baru
        </router-link>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="glass-card stat-card border-total">
        <span class="stat-label">{{ auth.isPemohon ? 'Total Permohonan Saya' : 'Total Permohonan Sistem' }}</span>
        <span class="stat-number">{{ auth.isPemohon ? stats.total : stats.total_projects }}</span>
      </div>

      <div v-if="auth.isAdmin" class="glass-card stat-card border-users">
        <span class="stat-label">Total User Terdaftar</span>
        <span class="stat-number text-indigo">{{ userCount }}</span>
      </div>
      <div v-else-if="auth.isPemohon" class="glass-card stat-card border-draft">
        <span class="stat-label">Draft (Belum Dikirim)</span>
        <span class="stat-number">{{ stats.draft }}</span>
      </div>
      <div v-else-if="auth.isPenilai" class="glass-card stat-card border-myreviews">
        <span class="stat-label">Penilaian Saya</span>
        <span class="stat-number text-indigo">{{ stats.my_reviews }}</span>
      </div>

      <div class="glass-card stat-card border-pending">
        <span class="stat-label">{{ auth.isPemohon ? 'Telah Dikirim / Dinilai' : 'Menunggu Penilaian' }}</span>
        <span class="stat-number text-amber">
          {{ auth.isPemohon ? (stats.submitted || 0) + (stats.under_review || 0) : stats.pending_review }}
        </span>
      </div>

      <div class="glass-card stat-card border-approved">
        <span class="stat-label">Disetujui</span>
        <span class="stat-number text-emerald">{{ stats.approved }}</span>
      </div>

      <div class="glass-card stat-card border-revision">
        <span class="stat-label">Perlu Revisi</span>
        <span class="stat-number text-amber">{{ stats.revision }}</span>
      </div>

      <div class="glass-card stat-card border-rejected">
        <span class="stat-label">Ditolak</span>
        <span class="stat-number text-rose">{{ stats.rejected }}</span>
      </div>
    </div>

    <!-- Quick Navigation Shortcuts (Admin Only) -->
    <div v-if="auth.isAdmin" class="quick-nav-grid">
      <router-link to="/document-types" class="glass-card nav-card">
        <div class="nav-card-icon">📄</div>
        <div>
          <h4>Master Jenis Dokumen</h4>
          <p>Kelola jenis dokumen & berkas wajib</p>
        </div>
      </router-link>

      <router-link to="/users" class="glass-card nav-card">
        <div class="nav-card-icon">👥</div>
        <div>
          <h4>Kelola User & Penilai</h4>
          <p>Tambah/Hapus akun Pemohon, Penilai, & Admin</p>
        </div>
      </router-link>

      <router-link to="/reviews" class="glass-card nav-card">
        <div class="nav-card-icon">⚖️</div>
        <div>
          <h4>Penilaian Project</h4>
          <p>Monitor & evaluasi antrean permohonan</p>
        </div>
      </router-link>

      <router-link to="/reviews/history" class="glass-card nav-card">
        <div class="nav-card-icon">📜</div>
        <div>
          <h4>Riwayat Audit & Log</h4>
          <p>Audit trail seluruh keputusan verifikator</p>
        </div>
      </router-link>

      <router-link to="/whatsapp-settings" class="glass-card nav-card">
        <div class="nav-card-icon">📱</div>
        <div>
          <h4>Gateway WhatsApp</h4>
          <p>Sesi login WA Admin & kirim pesan uji coba</p>
        </div>
      </router-link>
    </div>

    <!-- Spacious Full-Width Analytics Chart Card -->
    <div class="glass-card chart-card">
      <div class="chart-header">
        <h3 class="card-title">{{ auth.isPemohon ? 'Grafik Trend Pengajuan Bulanan' : 'Distribusi Status & Progres Permohonan' }}</h3>
        <small class="text-muted">Visualisasi data statistik real-time permohonan dokumen kelayakan.</small>
      </div>
      <div class="chart-container">
        <template v-if="chartData.loaded">
          <Bar v-if="auth.isPemohon" :data="barChartData" :options="barChartOptions" />
          <Doughnut v-else :data="doughnutChartData" :options="doughnutChartOptions" />
        </template>
        <div v-else class="chart-loading">Memuat Grafik Analytic...</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import { exportToCSV, formatDateCsv } from '../utils/exportCsv'
import { alertSuccess, alertError } from '../utils/swal'

import { 
  Chart as ChartJS, 
  Title, 
  Tooltip, 
  Legend, 
  ArcElement, 
  BarElement, 
  CategoryScale, 
  LinearScale 
} from 'chart.js'
import { Doughnut, Bar } from 'vue-chartjs'

ChartJS.register(Title, Tooltip, Legend, ArcElement, BarElement, CategoryScale, LinearScale)

const auth = useAuthStore()
const exporting = ref(false)
const userCount = ref(0)

const stats = reactive({
  total: 0,
  total_projects: 0,
  draft: 0,
  submitted: 0,
  under_review: 0,
  pending_review: 0,
  revision: 0,
  approved: 0,
  rejected: 0,
  my_reviews: 0
})

const chartData = reactive({
  loaded: false,
  labels: [],
  values: []
})

// Banner Computed Properties
const bannerTitle = computed(() => {
  if (auth.isAdmin) return 'Dashboard Executive Administrator 👑'
  if (auth.isPenilai) return 'Dashboard Penilai / Verifikator 📋'
  return `Selamat Datang, ${auth.userName} 👋`
})

const bannerSubtitle = computed(() => {
  if (auth.isAdmin) return 'Ringkasan statistik sistem, manajemen pengguna, serta laporan audit pengajuan dokumen.'
  if (auth.isPenilai) return 'Panel peninjauan, pencetakan rekapan, dan persetujuan permohonan dokumen kelayakan.'
  return 'Kelola permohonan dokumen kelayakan Anda di portal pelayanan digital.'
})

const bannerRoleClass = computed(() => {
  if (auth.isAdmin) return 'banner-admin'
  if (auth.isPenilai) return 'banner-penilai'
  return 'banner-pemohon'
})

// Chart Data Configs
const doughnutChartData = computed(() => ({
  labels: chartData.labels,
  datasets: [
    {
      backgroundColor: ['#94a3b8', '#3b82f6', '#c084fc', '#fbbf24', '#34d399', '#fb7185'],
      data: chartData.values
    }
  ]
}))

const doughnutChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'bottom', labels: { color: '#94a3b8' } }
  }
}

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

const barChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    x: { grid: { display: false }, ticks: { color: '#94a3b8' } },
    y: { grid: { color: 'rgba(255,255,255,0.05)' }, ticks: { color: '#94a3b8' } }
  }
}

// Export CSV Functionality
const handleExportRekapan = async () => {
  exporting.value = true
  try {
    if (auth.isPemohon) {
      const res = await apiClient.get('/projects?per_page=1000')
      const projects = res.data?.data || []
      const rows = [
        ['No', 'No. Project', 'Judul Permohonan', 'Perusahaan / Instansi', 'Unit Kerja', 'Status Permohonan', 'Tanggal Dibuat']
      ]
      projects.forEach((item, index) => {
        rows.push([
          index + 1,
          item.project_number || '-',
          item.title || '-',
          item.company_name || '-',
          item.unit || '-',
          item.status || '-',
          formatDateCsv(item.created_at)
        ])
      })
      exportToCSV(`Rekapan_Permohonan_Saya_${new Date().toISOString().slice(0, 10)}.csv`, rows)
    } else {
      const res = await apiClient.get('/reviews/all-history?per_page=1000')
      const histories = res.data?.data || []
      const rows = [
        ['No', 'Tanggal Verifikasi', 'No. Project', 'Judul Permohonan', 'Pemohon', 'Perusahaan / Instansi', 'Penilai / Verifikator', 'Status Keputusan', 'Catatan Evaluasi']
      ]
      histories.forEach((item, index) => {
        rows.push([
          index + 1,
          formatDateCsv(item.created_at),
          item.project?.project_number || '-',
          item.project?.title || '-',
          item.project?.user?.name || '-',
          item.project?.company_name || '-',
          item.reviewer?.name || 'Penilai',
          item.status_to || '-',
          item.notes || '-'
        ])
      })
      exportToCSV(`Rekapan_Sistem_${new Date().toISOString().slice(0, 10)}.csv`, rows)
    }
    alertSuccess('Export Berhasil!', 'Rekapan laporan berhasil diunduh.')
  } catch (err) {
    alertError('Gagal Export', err.response?.data?.error || err.message)
  } finally {
    exporting.value = false
  }
}

// Load Dashboard Data
const fetchData = async () => {
  try {
    const resStats = await apiClient.get('/dashboard/stats')
    if (resStats.data?.data) {
      Object.assign(stats, resStats.data.data)
    }

    if (auth.isAdmin) {
      const resUsers = await apiClient.get('/users?per_page=1')
      if (resUsers.data?.meta) {
        userCount.value = resUsers.data.meta.total
      }
    }

    const resChart = await apiClient.get('/dashboard/chart-data')
    if (auth.isPemohon && resChart.data?.data?.monthly_submissions) {
      const subs = resChart.data.data.monthly_submissions
      chartData.labels = subs.map(s => s.month)
      chartData.values = subs.map(s => s.count)
      chartData.loaded = true
    } else if (!auth.isPemohon && resChart.data?.data?.status_distribution) {
      const dist = resChart.data.data.status_distribution
      chartData.labels = dist.map(s => s.status)
      chartData.values = dist.map(s => s.count)
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
  gap: 1rem;
}

.banner-admin {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(30, 41, 59, 0.8) 100%);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.banner-penilai {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.15) 0%, rgba(30, 41, 59, 0.7) 100%);
  border: 1px solid rgba(168, 85, 247, 0.3);
}

.banner-pemohon {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, rgba(30, 41, 59, 0.7) 100%);
  border: 1px solid rgba(99, 102, 241, 0.3);
}

.banner-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
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
.border-users::before { background: #c084fc; }
.border-draft::before { background: #94a3b8; }
.border-myreviews::before { background: #818cf8; }
.border-pending::before { background: #f59e0b; }
.border-revision::before { background: #fbbf24; }
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
.text-indigo { color: #c084fc; }
.text-emerald { color: #34d399; }
.text-rose { color: #fb7185; }

.quick-nav-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 1rem;
}

.nav-card {
  padding: 1.25rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  text-decoration: none;
  transition: all 0.2s ease;
  border: 1px solid var(--border-color);
}

.nav-card:hover {
  transform: translateY(-3px);
  border-color: var(--accent-primary);
  background: rgba(255, 255, 255, 0.04);
}

.nav-card-icon {
  font-size: 1.8rem;
}

.nav-card h4 {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-main);
}

.nav-card p {
  font-size: 0.78rem;
  color: var(--text-muted);
  margin-top: 0.15rem;
}

.chart-card {
  padding: 1.5rem 2rem;
  min-height: 320px;
}

.chart-header {
  margin-bottom: 1.25rem;
}

.card-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: var(--text-main);
}

.text-muted { color: var(--text-muted); }

.chart-container {
  height: 260px;
  position: relative;
}

.chart-loading {
  padding: 3rem;
  text-align: center;
  color: var(--text-subtle);
  font-size: 0.9rem;
}

@media (max-width: 1400px) {
  .quick-nav-grid { grid-template-columns: repeat(3, 1fr); }
}

@media (max-width: 1280px) {
  .stats-grid { grid-template-columns: repeat(3, 1fr); }
  .quick-nav-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 640px) {
  .welcome-banner { flex-direction: column; align-items: flex-start; }
  .banner-actions { width: 100%; flex-direction: column; }
  .banner-actions button, .banner-actions a { width: 100%; text-align: center; }
  .stats-grid { grid-template-columns: 1fr 1fr; }
  .quick-nav-grid { grid-template-columns: 1fr; }
}
</style>
