<template>
  <div class="auth-card glass-card">
    <div class="auth-header">
      <div class="auth-brand-icon">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
      </div>
      <h2 class="auth-title">Masuk ke Sistem</h2>
      <p class="auth-subtitle">Sistem Pengajuan Permohonan Dokumen Kelayakan</p>
    </div>

    <form @submit.prevent="handleLogin" class="auth-form">
      <div v-if="errorMessage" class="alert-error">
        {{ errorMessage }}
      </div>

      <div class="form-group">
        <label class="form-label">Email</label>
        <input 
          v-model="email" 
          type="email" 
          class="form-input" 
          placeholder="nama@email.com" 
          required 
        />
      </div>

      <div class="form-group">
        <label class="form-label">Password</label>
        <input 
          v-model="password" 
          type="password" 
          class="form-input" 
          placeholder="••••••••" 
          required 
        />
      </div>

      <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
        <span v-if="loading">Memproses...</span>
        <span v-else>Masuk Sekarang</span>
      </button>
    </form>

    <div class="auth-footer">
      <span>Belum memiliki akun?</span>
      <router-link to="/register" class="auth-link">Daftar Akun Baru</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    await auth.login(email.value, password.value)
    router.push('/dashboard')
  } catch (err) {
    errorMessage.value = err.response?.data?.error || 'Gagal masuk. Periksa kembali email dan password Anda.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-card {
  width: 100%;
  max-width: 440px;
  padding: 2.5rem;
}

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}

.auth-brand-icon {
  width: 54px;
  height: 54px;
  margin: 0 auto 1rem;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(79, 70, 229, 0.4));
  border: 1px solid rgba(99, 102, 241, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-primary);
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

.auth-link:hover {
  text-decoration: underline;
}
</style>
