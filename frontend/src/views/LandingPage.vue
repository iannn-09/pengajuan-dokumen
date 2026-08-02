<template>
  <div class="landing-page" ref="pageRef">
    <!-- Floating Orbs Background -->
    <div class="floating-orbs" aria-hidden="true">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
      <div class="orb orb-4"></div>
    </div>

    <!-- Navbar -->
    <header class="landing-header" :class="{ scrolled: isScrolled }">
      <div class="landing-nav-container">
        <div class="brand">
          <div class="brand-icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
          </div>
          <div class="brand-text">
            <span class="brand-title">Pengajuan Dokumen</span>
            <span class="brand-subtitle">Portal Layanan Kelayakan Digital</span>
          </div>
        </div>
        <nav class="nav-menu">
          <a href="#home" class="nav-link" @click.prevent="scrollTo('home')">Home</a>
          <a href="#features" class="nav-link" @click.prevent="scrollTo('features')">Fitur Utama</a>
          <a href="#workflow" class="nav-link" @click.prevent="scrollTo('workflow')">Alur Proses</a>
        </nav>
        <div class="nav-auth">
          <template v-if="auth.isAuthenticated">
            <router-link to="/dashboard" class="btn btn-primary">
              Buka Dashboard ({{ auth.userName }}) &rarr;
            </router-link>
          </template>
          <template v-else>
            <router-link to="/login" class="btn btn-ghost">Masuk</router-link>
            <router-link to="/register" class="btn btn-primary">Daftar Akun</router-link>
          </template>
        </div>
        <!-- Mobile menu button -->
        <button class="mobile-menu-btn" @click="mobileMenuOpen = !mobileMenuOpen" aria-label="Toggle menu">
          <svg v-if="!mobileMenuOpen" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
          <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <!-- Mobile menu -->
      <div class="mobile-nav" v-if="mobileMenuOpen">
        <a href="#home" class="mobile-nav-link" @click.prevent="scrollTo('home'); mobileMenuOpen = false">Home</a>
        <a href="#features" class="mobile-nav-link" @click.prevent="scrollTo('features'); mobileMenuOpen = false">Fitur Utama</a>
        <a href="#workflow" class="mobile-nav-link" @click.prevent="scrollTo('workflow'); mobileMenuOpen = false">Alur Proses</a>
        <div class="mobile-nav-auth">
          <template v-if="!auth.isAuthenticated">
            <router-link to="/login" class="btn btn-ghost btn-block">Masuk</router-link>
            <router-link to="/register" class="btn btn-primary btn-block">Daftar Akun</router-link>
          </template>
          <template v-else>
            <router-link to="/dashboard" class="btn btn-primary btn-block">Buka Dashboard &rarr;</router-link>
          </template>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <section id="home" class="hero-section">
      <div class="hero-content reveal" ref="heroRef">
        <div class="hero-badge">
          <span class="pulse-dot"></span>
          Sistem Pengajuan Dokumen Kelayakan Modern & Transparan
        </div>
        <h1 class="hero-title">
          Kelola Permohonan<br>Dokumen Kelayakan dengan
          <span class="gradient-text-animated">Cepat, Transparan & Andal</span>
        </h1>
        <p class="hero-subtitle">
          Platform digital terpadu untuk perusahaan dan pemohon dalam mengajukan permohonan dokumen kelayakan, memantau verifikasi administrasi, hingga keputusan resmi secara real-time.
        </p>
        <div class="hero-ctas">
          <template v-if="auth.isAuthenticated">
            <router-link to="/dashboard" class="btn btn-primary btn-lg btn-shimmer">
              Masuk ke Dashboard Saya &rarr;
            </router-link>
          </template>
          <template v-else>
            <router-link to="/register" class="btn btn-primary btn-lg btn-shimmer">
              Mulai Buat Pengajuan &rarr;
            </router-link>
            <router-link to="/login" class="btn btn-ghost btn-lg">
              Masuk ke Akun
            </router-link>
          </template>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="features-section">
      <div class="section-header reveal">
        <span class="section-badge">Fitur</span>
        <h2 class="section-title">Fitur Unggulan Sistem</h2>
        <p class="section-subtitle">Dirancang untuk menjamin efisiensi permohonan dan ketepatan verifikasi oleh instansi.</p>
      </div>
      <div class="features-grid">
        <div class="feature-card reveal" v-for="(feature, idx) in features" :key="idx" :class="'reveal-delay-' + (idx % 4)">
          <div class="feature-icon-wrap" :class="'gradient-' + (idx + 1)">
            <span class="feature-icon">{{ feature.icon }}</span>
          </div>
          <h3>{{ feature.title }}</h3>
          <p>{{ feature.desc }}</p>
          <div class="feature-glow" :class="'glow-' + (idx + 1)"></div>
        </div>
      </div>
    </section>

    <!-- Workflow Timeline Section -->
    <section id="workflow" class="workflow-section">
      <div class="section-header reveal">
        <span class="section-badge">Alur</span>
        <h2 class="section-title">Alur Proses Permohonan</h2>
        <p class="section-subtitle">Empat langkah mudah dari draf hingga keputusan resmi kelayakan.</p>
      </div>
      <div class="timeline-container">
        <div class="timeline-line"></div>
        <div class="timeline-progress" :style="{ width: timelineProgress + '%' }"></div>
        <div class="timeline-steps">
          <div class="timeline-step reveal" v-for="(step, idx) in workflowSteps" :key="idx" :class="'reveal-delay-' + idx">
            <div class="step-marker" :class="{ active: timelineProgress >= (idx + 1) * 25 }">
              <span class="step-num">{{ String(idx + 1).padStart(2, '0') }}</span>
            </div>
            <div class="step-content">
              <h4>{{ step.title }}</h4>
              <p>{{ step.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <div class="cta-card reveal">
        <div class="cta-bg-animated"></div>
        <div class="cta-content">
          <h2>Siap Mengajukan Dokumen Kelayakan Anda?</h2>
          <p>Daftarkan akun perusahaan Anda sekarang dan nikmati proses layanan digital yang cepat, transparan, dan dapat dipantau secara real-time.</p>
          <div class="cta-actions">
            <router-link v-if="!auth.isAuthenticated" to="/register" class="btn btn-white btn-lg btn-shimmer">
              Daftar Akun Pemohon Sekarang &rarr;
            </router-link>
            <router-link v-else to="/dashboard" class="btn btn-white btn-lg btn-shimmer">
              Buka Dashboard Pengajuan &rarr;
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="landing-footer">
      <div class="footer-container">
        <div class="footer-brand">
          <div class="footer-logo">
            <div class="brand-icon small">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            </div>
            <span class="footer-title">Pengajuan Dokumen Kelayakan</span>
          </div>
          <p>Portal Pelayanan Digital Publik Berkinerja Tinggi & Safe Audit Log.</p>
        </div>
        <div class="footer-links">
          <div class="footer-col">
            <h4>Navigasi</h4>
            <a href="#features" @click.prevent="scrollTo('features')">Fitur Utama</a>
            <a href="#workflow" @click.prevent="scrollTo('workflow')">Alur Proses</a>
          </div>
          <div class="footer-col">
            <h4>Akun</h4>
            <router-link to="/login">Masuk</router-link>
            <router-link to="/register">Daftar Baru</router-link>
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        <span>&copy; 2026 Sistem Pengajuan Dokumen Kelayakan. All rights reserved.</span>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

const pageRef = ref(null)
const heroRef = ref(null)
const isScrolled = ref(false)
const mobileMenuOpen = ref(false)
const timelineProgress = ref(0)

const features = [
  { icon: '📁', title: 'Pengajuan & Upload Dokumen', desc: 'Dukungan unggah berkas kelayakan (PDF, JPG, PNG, DOCX) hingga 10MB dengan penamaan UUID aman dan validasi ketat.' },
  { icon: '🔍', title: 'Penilaian Multi-Tahap', desc: 'Verifikator mengevaluasi dokumen, memberikan catatan revisi, serta memutuskan status disetujui, revisi, atau ditolak.' },
  { icon: '📜', title: 'Audit Log & History Timeline', desc: 'Setiap perubahan status dan catatan evaluasi tercatat transparan pada histori audit log secara kronologis.' },
  { icon: '📊', title: 'Visual Analytics Dashboard', desc: 'Grafik Chart.js interaktif untuk statistik pengajuan bulanan, distribusi status, dan performa verifikasi real-time.' },
]

const workflowSteps = [
  { title: 'Buat Draf Project', desc: 'Isi formulir rincian permohonan dan informasi perusahaan.' },
  { title: 'Upload Berkas Lampiran', desc: 'Unggah dokumen pendukung teknis kelayakan.' },
  { title: 'Kirim & Verifikasi Penilai', desc: 'Penilai memeriksa kelengkapan dan kesesuaian dokumen.' },
  { title: 'Keputusan Final', desc: 'Terima hasil: Disetujui, Revisi, atau Ditolak.' },
]

function scrollTo(id) {
  const el = document.getElementById(id)
  if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

// Scroll handler
function handleScroll() {
  isScrolled.value = window.scrollY > 40

  // Timeline progress
  const workflowEl = document.getElementById('workflow')
  if (workflowEl) {
    const rect = workflowEl.getBoundingClientRect()
    const windowH = window.innerHeight
    if (rect.top < windowH && rect.bottom > 0) {
      const progress = Math.min(100, Math.max(0, ((windowH - rect.top) / (rect.height + windowH * 0.3)) * 100))
      timelineProgress.value = progress
    }
  }
}

// Intersection Observer for scroll reveal
let observer = null
function setupScrollReveal() {
  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('revealed')
          observer.unobserve(entry.target)
        }
      })
    },
    { threshold: 0.1, rootMargin: '0px 0px -50px 0px' }
  )
  const revealEls = document.querySelectorAll('.reveal')
  revealEls.forEach((el) => observer.observe(el))
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true })
  handleScroll()

  // Setup scroll reveal with a small delay to ensure DOM is rendered
  setTimeout(setupScrollReveal, 100)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  if (observer) observer.disconnect()
})
</script>

<style scoped>
/* ================================================================
   LANDING PAGE — PREMIUM REDESIGN
   ================================================================ */
.landing-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  color: var(--text-main);
  position: relative;
  overflow-x: hidden;
}

/* ── Floating Orbs Background ─────────────────────────────────── */
.floating-orbs {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}
.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float-orb 20s ease-in-out infinite;
}
.orb-1 {
  width: 500px; height: 500px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.35) 0%, transparent 70%);
  top: -10%; left: -5%;
  animation-duration: 22s;
}
.orb-2 {
  width: 400px; height: 400px;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.25) 0%, transparent 70%);
  top: 50%; right: -8%;
  animation-duration: 26s;
  animation-delay: -5s;
}
.orb-3 {
  width: 350px; height: 350px;
  background: radial-gradient(circle, rgba(168, 85, 247, 0.2) 0%, transparent 70%);
  bottom: 10%; left: 30%;
  animation-duration: 30s;
  animation-delay: -10s;
}
.orb-4 {
  width: 300px; height: 300px;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.2) 0%, transparent 70%);
  top: 25%; left: 60%;
  animation-duration: 24s;
  animation-delay: -15s;
}
@keyframes float-orb {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(30px, -40px) scale(1.05); }
  50% { transform: translate(-20px, 20px) scale(0.95); }
  75% { transform: translate(15px, 35px) scale(1.02); }
}

/* ── Scroll Reveal Animations ─────────────────────────────────── */
.reveal {
  opacity: 0;
  transform: translateY(30px);
  transition: opacity 0.7s ease-out, transform 0.7s ease-out;
}
.reveal.revealed {
  opacity: 1;
  transform: translateY(0);
}
.reveal-delay-0 { transition-delay: 0s; }
.reveal-delay-1 { transition-delay: 0.15s; }
.reveal-delay-2 { transition-delay: 0.3s; }
.reveal-delay-3 { transition-delay: 0.45s; }

/* ── Header Navbar ────────────────────────────────────────────── */
.landing-header {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: transparent;
  transition: all 0.35s ease;
}
.landing-header.scrolled {
  background: rgba(15, 23, 42, 0.92);
  backdrop-filter: blur(20px) saturate(1.4);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.15);
}
.landing-nav-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.brand { display: flex; align-items: center; gap: 0.75rem; z-index: 10; }
.brand-icon {
  width: 40px; height: 40px; border-radius: var(--radius-md);
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.25), rgba(168, 85, 247, 0.35));
  border: 1px solid rgba(99, 102, 241, 0.4);
  display: flex; align-items: center; justify-content: center;
  color: #a78bfa;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}
.brand-icon:hover {
  transform: scale(1.05);
  box-shadow: 0 0 20px rgba(99, 102, 241, 0.3);
}
.brand-icon.small { width: 28px; height: 28px; }
.brand-title { font-size: 1.05rem; font-weight: 800; color: var(--text-main); display: block; }
.brand-subtitle { font-size: 0.72rem; color: var(--text-muted); }
.nav-menu { display: flex; gap: 2rem; }
.nav-link {
  color: var(--text-muted);
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  transition: color 0.25s ease;
  position: relative;
}
.nav-link::after {
  content: '';
  position: absolute;
  bottom: -4px; left: 0;
  width: 0; height: 2px;
  background: linear-gradient(90deg, #6366f1, #a78bfa);
  border-radius: 2px;
  transition: width 0.3s ease;
}
.nav-link:hover { color: var(--text-main); }
.nav-link:hover::after { width: 100%; }
.nav-auth { display: flex; gap: 0.75rem; align-items: center; }

.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  color: var(--text-main);
  cursor: pointer;
  padding: 0.5rem;
  z-index: 10;
}
.mobile-nav {
  display: none;
  padding: 1rem 1.5rem 1.5rem;
  border-top: 1px solid rgba(255,255,255,0.06);
}
.mobile-nav-link {
  display: block;
  padding: 0.75rem 0;
  color: var(--text-muted);
  text-decoration: none;
  font-weight: 500;
  border-bottom: 1px solid rgba(255,255,255,0.04);
}
.mobile-nav-auth { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 1rem; }

/* ── Buttons ──────────────────────────────────────────────────── */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.6rem 1.25rem;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-size: 0.9rem;
  text-decoration: none;
  cursor: pointer;
  border: none;
  transition: all 0.3s ease;
}
.btn-primary {
  background: linear-gradient(135deg, #6366f1, #7c3aed);
  color: #fff;
  box-shadow: 0 4px 15px rgba(99, 102, 241, 0.3);
}
.btn-primary:hover {
  background: linear-gradient(135deg, #4f46e5, #6d28d9);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(99, 102, 241, 0.4);
}
.btn-ghost {
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-main);
  border: 1px solid rgba(255, 255, 255, 0.12);
}
.btn-ghost:hover {
  background: rgba(255, 255, 255, 0.12);
  transform: translateY(-2px);
}
.btn-white {
  background: #fff;
  color: #1e1b4b;
  font-weight: 700;
}
.btn-white:hover {
  background: #f0f0ff;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(255, 255, 255, 0.2);
}
.btn-lg { padding: 0.85rem 1.75rem; font-size: 1rem; }
.btn-block { width: 100%; }

/* Shimmer effect */
.btn-shimmer {
  position: relative;
  overflow: hidden;
}
.btn-shimmer::before {
  content: '';
  position: absolute;
  top: 0; left: -100%;
  width: 100%; height: 100%;
  background: linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.2) 50%, transparent 100%);
  animation: shimmer 3s infinite;
}
@keyframes shimmer {
  0% { left: -100%; }
  100% { left: 200%; }
}

/* ── Hero Section ─────────────────────────────────────────────── */
.hero-section {
  position: relative;
  z-index: 1;
  max-width: 900px;
  margin: 0 auto;
  padding: 8rem 1.5rem 5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}
.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.45rem 1rem;
  border-radius: 9999px;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.25);
  color: #a78bfa;
  font-size: 0.82rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  animation: badge-glow 3s ease-in-out infinite alternate;
}
@keyframes badge-glow {
  0% { box-shadow: 0 0 8px rgba(99, 102, 241, 0.1); }
  100% { box-shadow: 0 0 20px rgba(99, 102, 241, 0.25); }
}
.pulse-dot {
  width: 8px; height: 8px; border-radius: 50%; background: #6366f1;
  box-shadow: 0 0 10px #6366f1;
  animation: pulse-glow 2s ease-in-out infinite;
}
@keyframes pulse-glow {
  0%, 100% { opacity: 1; box-shadow: 0 0 10px #6366f1; }
  50% { opacity: 0.5; box-shadow: 0 0 20px #6366f1, 0 0 40px rgba(99, 102, 241, 0.3); }
}
.hero-title {
  font-size: 3rem;
  font-weight: 800;
  line-height: 1.15;
  letter-spacing: -0.03em;
  margin-bottom: 1.25rem;
}
.gradient-text-animated {
  background: linear-gradient(135deg, #818cf8 0%, #34d399 50%, #818cf8 100%);
  background-size: 200% 200%;
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: gradient-shift 4s ease infinite;
}
@keyframes gradient-shift {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}
.hero-subtitle {
  font-size: 1.05rem;
  color: var(--text-muted);
  line-height: 1.7;
  margin-bottom: 2rem;
  max-width: 550px;
  margin-left: auto;
  margin-right: auto;
}
.hero-ctas { display: flex; gap: 1rem; flex-wrap: wrap; justify-content: center; }

/* ── Section Headers ──────────────────────────────────────────── */
.features-section,
.workflow-section,
.cta-section {
  position: relative;
  z-index: 1;
  max-width: 1280px;
  margin: 0 auto;
  padding: 5rem 1.5rem;
  width: 100%;
}
.section-header {
  text-align: center;
  max-width: 650px;
  margin: 0 auto 3.5rem;
}
.section-badge {
  display: inline-block;
  padding: 0.3rem 0.85rem;
  border-radius: 9999px;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
  color: #818cf8;
  font-size: 0.78rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin-bottom: 1rem;
}
.section-title {
  font-size: 2.2rem;
  font-weight: 800;
  color: var(--text-main);
  letter-spacing: -0.02em;
  margin-bottom: 0.5rem;
}
.section-subtitle {
  font-size: 1rem;
  color: var(--text-muted);
  margin-top: 0.5rem;
  line-height: 1.6;
}

/* ── Feature Cards ────────────────────────────────────────────── */
.features-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}
.feature-card {
  position: relative;
  padding: 2rem 1.5rem;
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  overflow: hidden;
  transition: transform 0.35s ease, border-color 0.35s ease, box-shadow 0.35s ease;
}
.feature-icon { font-size: 2rem; }
.feature-card:hover {
  transform: translateY(-6px);
  border-color: rgba(99, 102, 241, 0.3);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2), 0 0 30px rgba(99, 102, 241, 0.08);
}
.feature-glow {
  position: absolute;
  top: -50%; left: -50%;
  width: 200%; height: 200%;
  opacity: 0;
  transition: opacity 0.5s ease;
  pointer-events: none;
}
.feature-card:hover .feature-glow { opacity: 1; }
.glow-1 { background: radial-gradient(circle at 60% 60%, rgba(99, 102, 241, 0.08) 0%, transparent 60%); }
.glow-2 { background: radial-gradient(circle at 60% 60%, rgba(16, 185, 129, 0.08) 0%, transparent 60%); }
.glow-3 { background: radial-gradient(circle at 60% 60%, rgba(168, 85, 247, 0.08) 0%, transparent 60%); }
.glow-4 { background: radial-gradient(circle at 60% 60%, rgba(59, 130, 246, 0.08) 0%, transparent 60%); }
.feature-icon-wrap {
  width: 50px; height: 50px;
  border-radius: var(--radius-md);
  display: flex; align-items: center; justify-content: center;
  font-size: 1.4rem;
  transition: transform 0.3s ease;
}
.feature-card:hover .feature-icon-wrap { transform: scale(1.1); }
.gradient-1 { background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(99, 102, 241, 0.05)); }
.gradient-2 { background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(16, 185, 129, 0.05)); }
.gradient-3 { background: linear-gradient(135deg, rgba(168, 85, 247, 0.2), rgba(168, 85, 247, 0.05)); }
.gradient-4 { background: linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(59, 130, 246, 0.05)); }
.feature-card h3 { font-size: 1.1rem; font-weight: 700; color: var(--text-main); }
.feature-card p { font-size: 0.88rem; color: var(--text-muted); line-height: 1.55; }

/* ── Workflow Timeline ────────────────────────────────────────── */
.timeline-container {
  position: relative;
  padding: 2rem 0;
}
.timeline-line {
  position: absolute;
  top: 32px; left: 0; right: 0;
  height: 3px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 3px;
}
.timeline-progress {
  position: absolute;
  top: 32px; left: 0;
  height: 3px;
  background: linear-gradient(90deg, #6366f1, #a78bfa, #34d399);
  border-radius: 3px;
  transition: width 0.3s ease;
  box-shadow: 0 0 12px rgba(99, 102, 241, 0.4);
}
.timeline-steps {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
  position: relative;
}
.step-marker {
  width: 56px; height: 56px;
  border-radius: 50%;
  background: var(--bg-secondary);
  border: 2px solid rgba(255, 255, 255, 0.1);
  display: flex; align-items: center; justify-content: center;
  margin: 0 auto 1.25rem;
  transition: all 0.5s ease;
  position: relative;
  z-index: 2;
}
.step-marker.active {
  border-color: #6366f1;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(168, 85, 247, 0.15));
  box-shadow: 0 0 20px rgba(99, 102, 241, 0.3);
}
.step-num {
  font-size: 1rem;
  font-weight: 800;
  color: var(--text-subtle);
  transition: color 0.5s ease;
}
.step-marker.active .step-num { color: #a78bfa; }
.step-content { text-align: center; }
.step-content h4 {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 0.4rem;
}
.step-content p { font-size: 0.85rem; color: var(--text-muted); line-height: 1.5; }

/* ── CTA Section ──────────────────────────────────────────────── */
.cta-card {
  position: relative;
  padding: 4rem 2rem;
  text-align: center;
  border-radius: var(--radius-xl);
  overflow: hidden;
  border: 1px solid rgba(99, 102, 241, 0.3);
}
.cta-bg-animated {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #312e81 0%, #1e1b4b 30%, #0f172a 60%, #064e3b 100%);
  background-size: 300% 300%;
  animation: cta-gradient 8s ease infinite;
  z-index: 0;
}
@keyframes cta-gradient {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}
.cta-content {
  position: relative;
  z-index: 1;
}
.cta-card h2 { font-size: 2rem; font-weight: 800; color: #fff; margin-bottom: 0.75rem; }
.cta-card p { font-size: 1rem; color: rgba(255,255,255,0.7); margin-bottom: 2rem; max-width: 600px; margin-left: auto; margin-right: auto; line-height: 1.6; }
.cta-actions { display: flex; justify-content: center; }

/* ── Footer ───────────────────────────────────────────────────── */
.landing-footer {
  position: relative;
  z-index: 1;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(15, 23, 42, 0.95);
  margin-top: auto;
}
.footer-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 3rem 1.5rem 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 3rem;
}
.footer-logo { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.5rem; }
.footer-title { font-weight: 700; font-size: 0.95rem; color: var(--text-main); }
.footer-brand p { font-size: 0.82rem; color: var(--text-muted); max-width: 300px; line-height: 1.5; }
.footer-links { display: flex; gap: 4rem; }
.footer-col { display: flex; flex-direction: column; gap: 0.5rem; }
.footer-col h4 { font-size: 0.82rem; font-weight: 700; color: var(--text-main); text-transform: uppercase; letter-spacing: 0.06em; margin-bottom: 0.25rem; }
.footer-col a { font-size: 0.85rem; color: var(--text-muted); text-decoration: none; transition: color 0.2s ease; }
.footer-col a:hover { color: var(--text-main); }
.footer-bottom {
  max-width: 1280px;
  margin: 0 auto;
  padding: 1.25rem 1.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
  text-align: center;
  font-size: 0.78rem;
  color: var(--text-subtle);
}

/* ── Responsive ───────────────────────────────────────────────── */
@media (max-width: 1024px) {
  .hero-section {
    grid-template-columns: 1fr;
    text-align: center;
    padding-top: 7rem;
  }
  .nav-menu { display: none; }
  .nav-auth { display: none; }
  .mobile-menu-btn { display: block; }
  .mobile-nav { display: block; }
  .hero-subtitle { margin-left: auto; margin-right: auto; }
  .hero-ctas { justify-content: center; }
  .features-grid { grid-template-columns: repeat(2, 1fr); }
  .timeline-steps { grid-template-columns: repeat(2, 1fr); gap: 2rem; }
  .timeline-line, .timeline-progress { display: none; }
  .footer-container { flex-direction: column; }
}
@media (max-width: 640px) {
  .nav-menu { display: none; }
  .hero-title { font-size: 2rem; }
  .section-title { font-size: 1.6rem; }
  .features-grid { grid-template-columns: 1fr; }
  .footer-container { flex-direction: column; gap: 1rem; text-align: center; }
  .timeline-steps { grid-template-columns: 1fr; }
  .footer-links { flex-direction: column; gap: 1.5rem; }
  .cta-card { padding: 2.5rem 1.5rem; }
  .cta-card h2 { font-size: 1.5rem; }
}
</style>