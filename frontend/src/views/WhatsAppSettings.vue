<template>
  <div class="wa-container">
    <!-- Banner Header -->
    <div class="glass-card welcome-banner">
      <div class="welcome-text">
        <h2>Integrasi Native WhatsApp Whatsmeow Admin 📱</h2>
        <p>Tautkan akun WhatsApp Admin secara native melalui pemindaian QR Code Whatsmeow.</p>
      </div>
      <div class="banner-actions">
        <button class="btn btn-secondary" @click="fetchStatus" :disabled="loading">
          🔄 Refresh Status
        </button>
        <button v-if="waStatus.connected" class="btn btn-danger" @click="handleDisconnect" :disabled="disconnecting">
          {{ disconnecting ? 'Memproses...' : '🔴 Log Out Sesi WA' }}
        </button>
      </div>
    </div>

    <!-- Status Overview & QR Section -->
    <div class="grid-2col">
      <!-- Connection Status Card -->
      <div class="glass-card card-box">
        <h3 class="card-title">Status Sesi Whatsmeow Admin</h3>
        
        <div class="status-box">
          <div class="status-badge-lg" :class="waStatus.connected ? 'status-connected' : 'status-disconnected'">
            <span class="status-dot-lg"></span>
            <span>{{ waStatus.connected ? 'TERHUBUNG (SESSION ACTIVE)' : 'BELUM TERHUBUNG / PERLU SCAN QR' }}</span>
          </div>
          <p class="status-desc">
            {{ waStatus.connected 
              ? `Sesi WhatsApp Admin (${waStatus.phone || 'Admin'}) terhubung secara aktif. Pesan notifikasi dikirimkan langsung dari akun WA Admin ini.` 
              : 'Sesi WhatsApp belum terhubung. Silakan pindai (scan) QR Code di samping menggunakan aplikasi WhatsApp pada ponsel Admin.' }}
          </p>
        </div>

        <div class="info-list">
          <div class="info-item">
            <span class="info-label">Pengirim Notifikasi:</span>
            <span class="info-value">Admin WhatsApp Gateway (Whatsmeow)</span>
          </div>
          <div class="info-item">
            <span class="info-label">Nomor WhatsApp Terhubung:</span>
            <span class="info-value code-text">{{ waStatus.phone || 'Belum Ada Sesi' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Status Engine:</span>
            <span class="info-value text-emerald">Native Go Whatsmeow Multi-Device</span>
          </div>
        </div>
      </div>

      <!-- QR Code & Pair Box -->
      <div class="glass-card card-box flex-center">
        <h3 class="card-title">Pindai (Scan) QR Code WhatsApp</h3>
        
        <div v-if="waStatus.connected" class="connected-state">
          <div class="check-circle">✓</div>
          <h4>Sesi Admin Terhubung!</h4>
          <p>WhatsApp Admin <strong>{{ waStatus.phone }}</strong> sudah aktif dan siap mengirimkan notifikasi otomatis.</p>
        </div>

        <div v-else-if="waStatus.qr_code" class="qr-box">
          <img :src="waStatus.qr_code" alt="WhatsApp QR Code" class="qr-img" />
          <p class="qr-tip">Buka WhatsApp di HP &rarr; Pengaturan &rarr; Perangkat Tertaut &rarr; Tautkan Perangkat</p>
        </div>

        <div v-else class="empty-qr">
          <div class="qr-placeholder-icon">⌛</div>
          <p>Memuat QR Code Whatsmeow... Silakan klik refresh jika QR Code belum tampil.</p>
          <button class="btn btn-primary btn-sm mt-2" @click="fetchStatus">
            Muat Ulang QR Code
          </button>
        </div>
      </div>
    </div>

    <!-- Test Send WhatsApp Section -->
    <div class="glass-card card-box mt-4">
      <h3 class="card-title">🧪 Uji Coba Kirim Pesan WhatsApp</h3>
      <p class="card-subtitle">Kirim pesan pengujian dari nomor Admin ke nomor tujuan mana pun.</p>

      <form @submit.prevent="handleSendTest" class="form-grid">
        <div class="form-group">
          <label class="form-label">Nomor WhatsApp Tujuan *</label>
          <input 
            v-model="testForm.target_phone" 
            type="text" 
            class="form-input" 
            placeholder="Contoh: 081234567890" 
            required 
          />
        </div>

        <div class="form-group">
          <label class="form-label">Isi Pesan Uji Coba *</label>
          <textarea 
            v-model="testForm.message" 
            class="form-input" 
            rows="3" 
            placeholder="Tuliskan pesan uji coba di sini..." 
            required
          ></textarea>
        </div>

        <button type="submit" class="btn btn-primary" :disabled="sendingTest || !waStatus.connected">
          <span>{{ sendingTest ? 'Mengirim...' : (waStatus.connected ? '🚀 Kirim Pesan Uji Coba' : '⚠️ Sambungkan WA Terlebih Dahulu') }}</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import apiClient from '../services/api'
import { alertSuccess, alertError, confirmDialog } from '../utils/swal'

const loading = ref(false)
const sendingTest = ref(false)
const disconnecting = ref(false)
let pollTimer = null

const waStatus = reactive({
  enabled: true,
  connected: false,
  phone: '',
  qr_code: null
})

const testForm = reactive({
  target_phone: '',
  message: 'Halo! Ini adalah pesan uji coba dari Sistem Pengajuan Dokumen Kelayakan.'
})

const fetchStatus = async () => {
  loading.value = true
  try {
    const res = await apiClient.get('/whatsapp/status')
    if (res.data?.data) {
      Object.assign(waStatus, res.data.data)
    }
  } catch (err) {
    console.error('Failed to fetch WA status:', err)
  } finally {
    loading.value = false
  }
}

const handleDisconnect = async () => {
  const isConfirmed = await confirmDialog(
    'Log Out Sesi WA?',
    'Sesi WhatsApp Admin yang terhubung akan diputuskan.',
    'Ya, Log Out'
  )
  if (!isConfirmed) return

  disconnecting.value = true
  try {
    await apiClient.post('/whatsapp/disconnect')
    alertSuccess('Log Out Berhasil', 'Sesi WhatsApp diputuskan. Silakan pindai QR Code baru.')
    fetchStatus()
  } catch (err) {
    alertError('Gagal Log Out', err.response?.data?.error || err.message)
  } finally {
    disconnecting.value = false
  }
}

const handleSendTest = async () => {
  sendingTest.value = true
  try {
    await apiClient.post('/whatsapp/test', testForm)
    alertSuccess('Pesan Terkirim!', `Pesan uji coba berhasil dikirim ke ${testForm.target_phone}`)
  } catch (err) {
    alertError('Gagal Kirim Pesan', err.response?.data?.error || err.message)
  } finally {
    sendingTest.value = false
  }
}

onMounted(() => {
  fetchStatus()
  pollTimer = setInterval(fetchStatus, 4000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<style scoped>
.wa-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.welcome-banner {
  padding: 1.75rem 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.15) 0%, rgba(30, 41, 59, 0.8) 100%);
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.banner-actions {
  display: flex;
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

.grid-2col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.card-box {
  padding: 1.75rem;
  display: flex;
  flex-direction: column;
}

.card-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 0.25rem;
}

.card-subtitle {
  font-size: 0.82rem;
  color: var(--text-muted);
  margin-bottom: 1.25rem;
}

.status-box {
  margin: 1rem 0 1.5rem;
}

.status-badge-lg {
  display: inline-flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.55rem 1.25rem;
  border-radius: 9999px;
  font-weight: 800;
  font-size: 0.85rem;
  letter-spacing: 0.02em;
}

.status-connected {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.4);
}

.status-disconnected {
  background: rgba(244, 63, 94, 0.15);
  color: #fb7185;
  border: 1px solid rgba(244, 63, 94, 0.4);
}

.status-dot-lg {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 10px currentColor;
}

.status-desc {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-top: 0.75rem;
  line-height: 1.4;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding-top: 1rem;
  border-top: 1px dashed var(--border-color);
}

.info-item {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
}

.info-label {
  color: var(--text-muted);
}

.info-value {
  font-weight: 600;
  color: var(--text-main);
}

.code-text {
  font-family: monospace;
  color: var(--accent-primary);
}

.text-emerald {
  color: #34d399;
}

.flex-center {
  align-items: center;
  text-align: center;
}

.connected-state {
  margin-top: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.check-circle {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
  border: 2px solid #34d399;
  font-size: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
}

.connected-state h4 {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-main);
}

.connected-state p {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-top: 0.3rem;
  max-width: 280px;
}

.qr-box {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.qr-img {
  width: 210px;
  height: 210px;
  border-radius: 12px;
  border: 4px solid var(--border-color);
  background: white;
  padding: 8px;
}

.qr-tip {
  font-size: 0.78rem;
  color: var(--text-muted);
  margin-top: 0.85rem;
  max-width: 280px;
}

.empty-qr {
  margin-top: 2rem;
  color: var(--text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.qr-placeholder-icon {
  font-size: 3rem;
  margin-bottom: 0.5rem;
}

.form-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.mt-2 { margin-top: 0.5rem; }
.mt-4 { margin-top: 1rem; }
.btn-sm { padding: 0.4rem 0.85rem; font-size: 0.82rem; }

.btn-danger {
  background: rgba(244, 63, 94, 0.2);
  color: #fb7185;
  border: 1px solid rgba(244, 63, 94, 0.4);
}
.btn-danger:hover {
  background: rgba(244, 63, 94, 0.3);
}

@media (max-width: 1024px) {
  .grid-2col { grid-template-columns: 1fr; }
}
</style>
