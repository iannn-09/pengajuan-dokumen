<template>
  <div class="auth-card glass-card">
    <div class="auth-header">
      <h2 class="auth-title">Daftar Akun Pemohon</h2>
      <p class="auth-subtitle">Lengkapi formulir di bawah ini untuk membuat akun pemohon dokumen baru</p>
    </div>

    <form @submit.prevent="handleRegister" class="auth-form">
      <div v-if="errorMessage" class="alert-error">
        {{ errorMessage }}
      </div>

      <div class="form-group">
        <label class="form-label">Nama Lengkap / Penanggung Jawab *</label>
        <input 
          v-model="form.name" 
          type="text" 
          class="form-input" 
          placeholder="Nama Anda" 
          required 
        />
      </div>

      <div class="form-group">
        <label class="form-label">Email Pemohon *</label>
        <input 
          v-model="form.email" 
          type="email" 
          class="form-input" 
          placeholder="nama@perusahaan.com" 
          required 
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">Nomor Telepon</label>
          <input 
            v-model="form.phone" 
            type="text" 
            class="form-input" 
            placeholder="0812xxxx" 
          />
        </div>
        <div class="form-group">
          <label class="form-label">Nama Perusahaan / Instansi</label>
          <input 
            v-model="form.company" 
            type="text" 
            class="form-input" 
            placeholder="PT. Contoh Kelayakan" 
          />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">Password * (Min 8 Karakter)</label>
        <input 
          v-model="form.password" 
          type="password" 
          class="form-input" 
          placeholder="••••••••" 
          required 
          minlength="8"
        />
      </div>

      <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
        <span v-if="loading">Memproses Pendaftaran...</span>
        <span v-else>Daftar Akun Pemohon</span>
      </button>
    </form>

    <div class="auth-footer">
      <span>Sudah memiliki akun?</span>
      <router-link to="/login" class="auth-link">Masuk Ke Sistem</router-link>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const loading = ref(false)
const errorMessage = ref('')

const form = reactive({
  name: '',
  email: '',
  password: '',
  phone: '',
  company: ''
})

const handleRegister = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    await auth.register(form)
    router.push('/dashboard')
  } catch (err) {
    errorMessage.value = err.response?.data?.error || 'Gagal mendaftar. Pastikan email belum terdaftar.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-card {
  width: 100%;
  max-width: 520px;
  padding: 2.5rem;
}

.auth-header {
  text-align: center;
  margin-bottom: 1.5rem;
}

.auth-title {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--text-main);
}

.auth-subtitle {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-top: 0.3rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.alert-error {
  background: rgba(244, 63, 94, 0.12);
  color: #fb7185;
  border: 1px solid rgba(244, 63, 94, 0.3);
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  font-size: 0.85rem;
  margin-bottom: 1.25rem;
}

.btn-block {
  width: 100%;
  padding: 0.8rem;
  font-size: 0.95rem;
  margin-top: 0.5rem;
}

.auth-footer {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.85rem;
  color: var(--text-muted);
  display: flex;
  justify-content: center;
  gap: 0.4rem;
}

.auth-link {
  color: var(--accent-primary);
  text-decoration: none;
  font-weight: 600;
}

@media (max-width: 640px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
