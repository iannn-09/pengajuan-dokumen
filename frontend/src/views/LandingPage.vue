<template>
  <div class="landing-page">
    <!-- Navbar -->
    <header class="landing-header">
      <div class="landing-nav-container">
        <div class="brand">
          <div class="brand-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
          </div>
          <div class="brand-text">
            <span class="brand-title">Pengajuan Dokumen</span>
            <span class="brand-subtitle">Portal Layanan Kelayakan Digital</span>
          </div>
        </div>

        <div class="nav-menu">
          <a href="#features" class="nav-link">Fitur Utama</a>
          <a href="#workflow" class="nav-link">Alur Proses</a>
          <a href="#stats" class="nav-link">Keunggulan</a>
        </div>

        <div class="nav-auth">
          <template v-if="auth.isAuthenticated">
            <router-link to="/dashboard" class="btn btn-primary">
              Buka Dashboard ({{ auth.userName }}) &rarr;
            </router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="btn btn-secondary">Masuk</router-link>
            <router-link to="/register" class="btn btn-primary">Daftar Akun</router-link>
          </template>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-badge">
          <span class="pulse-dot"></span>
          Sistem Pengajuan Dokumen Kelayakan Modern & Transparan
        </div>
        <h1 class="hero-title">
          Kelola Permohonan Dokumen Kelayakan dengan <span class="gradient-text">Cepat, Transparan & Andal</span>
        </h1>
        <p class="hero-subtitle">
          Platform digital terpadu untuk perusahaan dan pemohon dalam mengajukan permohonan dokumen kelayakan, memantau verifikasi administrasi, hingga keputusan resmi secara real-time.
        </p>

        <div class="hero-ctas">
          <template v-if="auth.isAuthenticated">
            <router-link to="/dashboard" class="btn btn-primary btn-lg">
              Masuk ke Dashboard Saya &rarr;
            </router-link>
          </template>
          <template v-else>
            <router-link to="/register" class="btn btn-primary btn-lg">
              Mulai Buat Pengajuan &rarr;
            </router-link>
            <router-link to="/login" class="btn btn-secondary btn-lg">
              Masuk ke Akun
            </router-link>
          </template>
        </div>
      </div>

      <!-- Hero Visual Card / Mockup -->
      <div class="hero-visual">
        <div class="glass-card visual-card">
          <div class="mock-header">
            <span class="dot red"></span>
            <span class="dot yellow"></span>
            <span class="dot green"></span>
            <span class="mock-title">Dashboard Status Real-time</span>
          </div>
          <div class="mock-body">
            <div class="mock-stat">
              <span class="label">Project Permohonan</span>
              <span class="val font-mono">PRJ-20260730-001</span>
            </div>
            <div class="mock-status">
              <span class="badge-approved">✓ Disetujui (Approved)</span>
            </div>
            <div class="mock-timeline">
              <div class="t-line active">1. Pengajuan Draft</div>
              <div class="t-line active">2. Verifikasi Penilai</div>
              <div class="t-line active">3. Keputusan Final</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="features-section">
      <div class="section-header">
        <h2 class="section-title">Fitur Unggulan Sistem</h2>
        <p class="section-subtitle">Dirancang untuk menjamin efisiensi permohonan dan ketepatan verifikasi oleh instansi.</p>
      </div>

      <div class="features-grid">
        <div class="glass-card feature-card">
          <div class="feature-icon">📁</div>
          <h3>Pengajuan & Upload Dokumen</h3>
          <p>Dukungan unggah berkas kelayakan (PDF, JPG, PNG, DOCX) hingga 10MB dengan enkripsi dan penyimpanan UUID aman.</p>
        </div>

        <div class="glass-card feature-card">
          <div class="feature-icon">🔍</div>
          <h3>Penilaian Multi-Tahap</h3>
          <p>Verifikator/Penilai dapat mengevaluasi dokumen, memberikan catatan perbaikan, serta memutuskan status disetujui, revisi, atau ditolak.</p>
        </div>

        <div class="glass-card feature-card">
          <div class="feature-icon">📜</div>
          <h3>Audit Log & History Timeline</h3>
          <p>Setiap perubahan status dan catatan evaluasi tercatat transparan pada histori audit log yang dapat dipantau setiap saat.</p>
        </div>

        <div class="glass-card feature-card">
          <div class="feature-icon">📊</div>
          <h3>Visual Analytics Dashboard</h3>
          <p>Integrasi grafik Chart.js interaktif untuk statistik pengajuan bulanan, distribusi status, dan performa verifikasi.</p>
        </div>
      </div>
    </section>

    <!-- Workflow Section -->
    <section id="workflow" class="workflow-section">
      <div class="section-header">
        <h2 class="section-title">Alur Proses Permohonan</h2>
        <p class="section-subtitle">Empat langkah mudah dari draf hingga keputusan resmi kelayakan.</p>
      </div>

      <div class="workflow-steps">
        <div class="step-card">
          <div class="step-num">01</div>
          <h4>Buat Draf Project</h4>
          <p>Isi formulir rincian permohonan dan informasi perusahaan.</p>
        </div>
        <div class="step-arrow">&rarr;</div>

        <div class="step-card">
          <div class="step-num">02</div>
          <h4>Upload Berkas Lampiran</h4>
          <p>Unggah dokumen pendukung teknis kelayakan.</p>
        </div>
        <div class="step-arrow">&rarr;</div>

        <div class="step-card">
          <div class="step-num">03</div>
          <h4>Kirim & Verifikasi Penilai</h4>
          <p>Penilai memeriksa kelengkapan & kesesuaian dokumen.</p>
        </div>
        <div class="step-arrow">&rarr;</div>

        <div class="step-card">
          <div class="step-num">04</div>
          <h4>Keputusan Final</h4>
          <p>Terima hasil Disetujui (Approved), Revisi, atau Ditolak.</p>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <div class="glass-card cta-card">
        <h2>Siap Mengajukan Dokumen Kelayakan Anda?</h2>
        <p>Daftarkan akun perusahaan Anda sekarang dan nikmati proses layanan digital yang cepat dan transparan.</p>
        <div class="cta-actions">
          <router-link v-if="!auth.isAuthenticated" to="/register" class="btn btn-primary btn-lg">
            Daftar Akun Pemohon Sekarang
          </router-link>
          <router-link v-else to="/dashboard" class="btn btn-primary btn-lg">
            Buka Dashboard Pengajuan &rarr;
          </router-link>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="landing-footer">
      <div class="footer-container">
        <div class="footer-brand">
          <span class="footer-title">Sistem Pengajuan Dokumen Kelayakan</span>
          <p>Portal Pelayanan Digital Publik Berkinerja Tinggi & Safe Audit Log.</p>
        </div>
        <div class="footer-copy">
          &copy; 2026 Sistem Pengajuan Dokumen Kelayakan. All rights reserved.
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
</script>

<style scoped>
.landing-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  color: var(--text-main);
}

/* Header Navbar */
.landing-header {
  position: sticky;
  top: 0;
  z-index: 50;
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border-color);
}

.landing-nav-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand { display: flex; align-items: center; gap: 0.75rem; }
.brand-icon {
  width: 40px; height: 40px; border-radius: var(--radius-md);
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(79, 70, 229, 0.4));
  border: 1px solid rgba(99, 102, 241, 0.4);
  display: flex; align-items: center; justify-content: center;
  color: var(--accent-primary);
}
.brand-title { font-size: 1.05rem; font-weight: 800; color: var(--text-main); display: block; }
.brand-subtitle { font-size: 0.72rem; color: var(--text-muted); }

.nav-menu { display: flex; gap: 1.75rem; }
.nav-link { color: var(--text-muted); text-decoration: none; font-size: 0.9rem; font-weight: 500; transition: color 0.2s ease; }
.nav-link:hover { color: var(--text-main); }

.nav-auth { display: flex; gap: 0.75rem; align-items: center; }

/* Hero Section */
.hero-section {
  max-width: 1280px;
  margin: 0 auto;
  padding: 5rem 1.5rem 4rem;
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 3rem;
  align-items: center;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.4rem 0.9rem;
  border-radius: 9999px;
  background: rgba(99, 102, 241, 0.12);
  border: 1px solid rgba(99, 102, 241, 0.3);
  color: #818cf8;
  font-size: 0.82rem;
  font-weight: 600;
  margin-bottom: 1.25rem;
}

.pulse-dot {
  width: 8px; height: 8px; border-radius: 50%; background: #6366f1;
  box-shadow: 0 0 10px #6366f1;
}

.hero-title {
  font-size: 2.8rem;
  font-weight: 800;
  line-height: 1.2;
  letter-spacing: -0.03em;
  margin-bottom: 1.25rem;
}

.gradient-text {
  background: linear-gradient(135deg, #818cf8 0%, #34d399 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-subtitle {
  font-size: 1.05rem;
  color: var(--text-muted);
  line-height: 1.6;
  margin-bottom: 2rem;
  max-width: 600px;
}

.hero-ctas {
  display: flex;
  gap: 1rem;
}

.btn-lg {
  padding: 0.85rem 1.75rem;
  font-size: 1rem;
}

/* Visual Card Mockup */
.visual-card {
  padding: 1.5rem;
  border-radius: var(--radius-xl);
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.12);
  box-shadow: var(--shadow-lg), var(--shadow-glow);
}

.mock-header {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  margin-bottom: 1.25rem;
  padding-bottom: 0.85rem;
  border-bottom: 1px solid var(--border-color);
}

.dot { width: 10px; height: 10px; border-radius: 50%; }
.dot.red { background: #ef4444; }
.dot.yellow { background: #f59e0b; }
.dot.green { background: #10b981; }

.mock-title {
  margin-left: 0.5rem;
  font-size: 0.82rem;
  color: var(--text-muted);
  font-weight: 600;
}

.mock-stat { display: flex; justify-content: space-between; margin-bottom: 1rem; }
.mock-stat .label { font-size: 0.8rem; color: var(--text-muted); }
.mock-stat .val { font-size: 0.85rem; color: var(--accent-primary); font-weight: 700; }

.mock-status { margin-bottom: 1.5rem; }
.badge-approved {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.3);
  padding: 0.4rem 0.85rem;
  border-radius: 9999px;
  font-size: 0.82rem;
  font-weight: 700;
}

.mock-timeline { display: flex; flex-direction: column; gap: 0.6rem; }
.t-line {
  padding: 0.6rem 0.85rem;
  border-radius: var(--radius-md);
  background: var(--bg-input);
  font-size: 0.82rem;
  color: var(--text-muted);
  border-left: 3px solid var(--accent-primary);
}

/* Features Section */
.features-section, .workflow-section, .cta-section {
  max-width: 1280px;
  margin: 0 auto;
  padding: 4rem 1.5rem;
  width: 100%;
}

.section-header {
  text-align: center;
  max-width: 650px;
  margin: 0 auto 3rem;
}

.section-title {
  font-size: 2rem;
  font-weight: 800;
  color: var(--text-main);
  letter-spacing: -0.02em;
}

.section-subtitle {
  font-size: 0.95rem;
  color: var(--text-muted);
  margin-top: 0.5rem;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}

.feature-card {
  padding: 1.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.feature-icon { font-size: 2rem; }
.feature-card h3 { font-size: 1.1rem; font-weight: 700; color: var(--text-main); }
.feature-card p { font-size: 0.88rem; color: var(--text-muted); line-height: 1.5; }

/* Workflow Section */
.workflow-steps {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.step-card {
  flex: 1;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  text-align: center;
}

.step-num {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--accent-primary);
  margin-bottom: 0.5rem;
}

.step-card h4 { font-size: 1rem; font-weight: 700; color: var(--text-main); margin-bottom: 0.4rem; }
.step-card p { font-size: 0.82rem; color: var(--text-muted); }
.step-arrow { font-size: 1.5rem; color: var(--text-subtle); }

/* CTA Section */
.cta-card {
  padding: 3.5rem 2rem;
  text-align: center;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2) 0%, rgba(30, 41, 59, 0.9) 100%);
  border: 1px solid rgba(99, 102, 241, 0.4);
}

.cta-card h2 { font-size: 2rem; font-weight: 800; color: var(--text-main); }
.cta-card p { font-size: 1rem; color: var(--text-muted); margin: 0.75rem auto 2rem; max-width: 600px; }
.cta-actions { display: flex; justify-content: center; }

/* Footer */
.landing-footer {
  border-top: 1px solid var(--border-color);
  background: rgba(15, 23, 42, 0.95);
  margin-top: auto;
}

.footer-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.footer-title { font-weight: 700; font-size: 0.95rem; color: var(--text-main); }
.footer-brand p { font-size: 0.8rem; color: var(--text-muted); margin-top: 0.2rem; }
.footer-copy { font-size: 0.8rem; color: var(--text-subtle); }

@media (max-width: 1024px) {
  .hero-section { grid-template-columns: 1fr; text-align: center; }
  .hero-ctas { justify-content: center; }
  .features-grid { grid-template-columns: repeat(2, 1fr); }
  .workflow-steps { flex-direction: column; }
  .step-arrow { transform: rotate(90deg); }
}

@media (max-width: 640px) {
  .nav-menu { display: none; }
  .features-grid { grid-template-columns: 1fr; }
  .footer-container { flex-direction: column; gap: 1rem; text-align: center; }
}
</style>
