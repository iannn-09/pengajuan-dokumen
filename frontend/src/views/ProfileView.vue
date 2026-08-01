<template>
  <div class="profile-container">
    <div class="header-nav">
      <h1 class="page-title">Profil Saya & Pengaturan Akun</h1>
      <p class="page-subtitle">Kelola informasi data diri, foto profil, dan kata sandi akun Anda.</p>
    </div>

    <div class="grid-2col">
      <!-- Left Column: User Summary Card -->
      <div class="glass-card summary-card">
        <div class="avatar-box">
          <img 
            v-if="auth.user?.avatar" 
            :src="getAvatarUrl(auth.user.avatar)" 
            :alt="auth.userName" 
            class="avatar-large-img" 
          />
          <div v-else class="avatar-large">
            {{ (form.name || auth.userName || 'U').charAt(0).toUpperCase() }}
          </div>

          <button 
            type="button" 
            class="avatar-upload-btn" 
            @click="$refs.avatarInput.click()" 
            :disabled="uploadingAvatar"
            title="Klik untuk ubah foto profil"
          >
            <span v-if="uploadingAvatar">...</span>
            <span v-else>📷</span>
          </button>

          <input 
            type="file" 
            ref="avatarInput" 
            class="hidden-input" 
            accept="image/png, image/jpeg, image/jpg, image/webp" 
            @change="handleAvatarUpload" 
          />
          
          <span class="role-badge" :class="auth.userRole">{{ getRoleLabel(auth.userRole) }}</span>
        </div>

        <h3 class="user-display-name">{{ form.name || auth.userName }}</h3>
        <p class="user-email-text">{{ auth.user?.email }}</p>

        <div class="summary-details">
          <div class="summary-item">
            <span class="item-label">Status Akun:</span>
            <span class="item-value text-success">Active ✓</span>
          </div>
          <div class="summary-item">
            <span class="item-label">Perusahaan / Instansi:</span>
            <span class="item-value">{{ form.company || auth.userCompany || auth.user?.company || '-' }}</span>
          </div>
          <div class="summary-item">
            <span class="item-label">No. Telepon / HP:</span>
            <span class="item-value">{{ form.phone || auth.user?.phone || '-' }}</span>
          </div>
        </div>
      </div>

      <!-- Right Column: Edit Profile Form -->
      <div class="glass-card edit-card">
        <h3 class="card-title">Edit Data Diri</h3>

        <form @submit.prevent="handleSaveProfile" class="profile-form">
          <div class="form-group">
            <label class="form-label">Nama Lengkap *</label>
            <input 
              v-model="form.name" 
              type="text" 
              class="form-input" 
              placeholder="Masukkan nama lengkap Anda..." 
              required 
            />
          </div>

          <div class="form-group">
            <label class="form-label">Alamat Email (Otomatis)</label>
            <input 
              :value="auth.user?.email" 
              type="email" 
              class="form-input form-disabled" 
              disabled 
            />
            <small class="form-hint">🔒 Email digunakan sebagai ID login utama dan tidak dapat diubah.</small>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">No. Telepon / HP</label>
              <input 
                v-model="form.phone" 
                type="text" 
                class="form-input" 
                placeholder="Contoh: 081234567890" 
              />
            </div>

            <div class="form-group">
              <label class="form-label">Nama Perusahaan / Instansi</label>
              <input 
                v-model="form.company" 
                type="text" 
                class="form-input" 
                placeholder="Contoh: PT Medika Utama" 
              />
              <small class="form-hint">Otomatis terisi saat Anda membuat permohonan baru.</small>
            </div>
          </div>

          <hr class="form-divider" />

          <div class="form-group">
            <label class="form-label">Kata Sandi Baru (Opsional)</label>
            <input 
              v-model="form.password" 
              type="password" 
              class="form-input" 
              placeholder="Kosongkan jika tidak ingin mengubah kata sandi..." 
              minlength="6"
            />
            <small class="form-hint">Minimal 6 karakter. Biarkan kosong untuk tetap menggunakan kata sandi lama.</small>
          </div>

          <div class="form-actions">
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span>{{ submitting ? 'Menyimpan Perubahan...' : 'Simpan Perubahan Data Diri ✓' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import apiClient from '../services/api'
import { alertSuccess, alertError } from '../utils/swal'

const auth = useAuthStore()
const submitting = ref(false)
const uploadingAvatar = ref(false)
const avatarInput = ref(null)

const form = reactive({
  name: '',
  phone: '',
  company: '',
  password: ''
})

const populateForm = () => {
  if (auth.user) {
    form.name = auth.user.name || ''
    form.phone = auth.user.phone || ''
    form.company = auth.user.company || ''
    form.password = ''
  }
}

const handleAvatarUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('avatar', file)

  uploadingAvatar.value = true
  try {
    const res = await apiClient.post('/auth/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    if (res.data?.data) {
      auth.updateUser(res.data.data)
      alertSuccess('Foto Profil Diperbarui', 'Foto profil Anda berhasil diunggah.')
    }
  } catch (err) {
    alertError('Gagal Upload Foto', err.response?.data?.error || err.message)
  } finally {
    uploadingAvatar.value = false
    if (avatarInput.value) avatarInput.value.value = ''
  }
}

const handleSaveProfile = async () => {
  submitting.value = true
  try {
    const payload = {
      name: form.name,
      phone: form.phone,
      company: form.company
    }
    if (form.password.trim()) {
      payload.password = form.password.trim()
    }

    const res = await apiClient.put('/auth/profile', payload)
    if (res.data?.data) {
      auth.updateUser(res.data.data)
      alertSuccess('Profil Diperbarui!', 'Data diri Anda berhasil diperbarui.')
      form.password = ''
    }
  } catch (err) {
    alertError('Gagal Menyimpan', err.response?.data?.error || err.message)
  } finally {
    submitting.value = false
  }
}

const getAvatarUrl = (path) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${apiClient.defaults.baseURL.replace('/api/v1', '')}${path}`
}

const getRoleLabel = (role) => {
  if (role === 'admin') return 'Administrator'
  if (role === 'penilai') return 'Penilai / Verifikator'
  return 'Pemohon Dokumen'
}

onMounted(async () => {
  if (!auth.user) {
    await auth.fetchMe()
  }
  populateForm()
})
</script>

<style scoped>
.profile-container {
  max-width: 960px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.page-title {
  font-size: 1.4rem;
  font-weight: 800;
  color: var(--text-main);
}

.page-subtitle {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-top: 0.15rem;
}

.grid-2col {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 1.5rem;
}

.summary-card {
  padding: 1.75rem 1.25rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.avatar-box {
  position: relative;
  margin-bottom: 1rem;
}

.avatar-large {
  width: 84px;
  height: 84px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 2.2rem;
  font-weight: 800;
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.35);
}

.avatar-large-img {
  width: 84px;
  height: 84px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--accent-primary);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.35);
}

.avatar-upload-btn {
  position: absolute;
  right: -4px;
  bottom: 6px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--accent-primary);
  border: 2px solid #0f172a;
  color: white;
  font-size: 0.85rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.4);
  transition: transform 0.15s ease;
}

.avatar-upload-btn:hover {
  transform: scale(1.1);
  background: #4f46e5;
}

.hidden-input {
  display: none;
}

.role-badge {
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
  font-size: 0.68rem;
  font-weight: 800;
  padding: 0.15rem 0.6rem;
  border-radius: 9999px;
  background: var(--accent-primary);
  color: white;
}

.role-badge.admin { background: #ef4444; }
.role-badge.penilai { background: #f59e0b; }
.role-badge.pemohon { background: #6366f1; }

.user-display-name {
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--text-main);
  margin-top: 0.5rem;
}

.user-email-text {
  font-size: 0.82rem;
  color: var(--text-muted);
  margin-top: 0.15rem;
}

.summary-details {
  width: 100%;
  margin-top: 1.5rem;
  padding-top: 1.25rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  text-align: left;
}

.summary-item {
  display: flex;
  flex-direction: column;
}

.item-label {
  font-size: 0.75rem;
  color: var(--text-subtle);
  font-weight: 600;
}

.item-value {
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--text-main);
  margin-top: 0.1rem;
}

.text-success {
  color: #34d399;
}

.edit-card {
  padding: 1.75rem;
}

.card-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 1.25rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-disabled {
  background: rgba(255, 255, 255, 0.04);
  color: var(--text-muted);
  cursor: not-allowed;
}

.form-hint {
  font-size: 0.75rem;
  color: var(--text-subtle);
  margin-top: 0.25rem;
  display: block;
}

.form-divider {
  border: none;
  border-top: 1px dashed var(--border-color);
  margin: 1.5rem 0;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

@media (max-width: 768px) {
  .grid-2col {
    grid-template-columns: 1fr;
  }
}
</style>
