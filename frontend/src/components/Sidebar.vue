<template>
  <aside class="sidebar" :class="{ collapsed: isCollapsed }">
    <div class="sidebar-header">
      <div class="brand" v-if="!isCollapsed">
        <div class="brand-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
        </div>
        <div class="brand-text">
          <span class="brand-title">Pengajuan Dokumen</span>
          <span class="brand-subtitle">Sistem Kelayakan</span>
        </div>
      </div>
      <button class="collapse-btn" @click="isCollapsed = !isCollapsed">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline :points="isCollapsed ? '9 18 15 12 9 6' : '15 18 9 12 15 6'"/></svg>
      </button>
    </div>

    <nav class="sidebar-nav">
      <router-link to="/dashboard" class="nav-item" active-class="active" :title="isCollapsed ? 'Dashboard' : ''">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
        <span v-if="!isCollapsed">Dashboard</span>
      </router-link>

      <!-- Admin Section -->
      <template v-if="auth.isAdmin">
        <div class="nav-label" v-if="!isCollapsed">ADMINISTRATOR</div>
        <router-link to="/document-types" class="nav-item" active-class="active" :title="isCollapsed ? 'Master Jenis Dokumen' : ''">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
          <span v-if="!isCollapsed">Master Jenis Dokumen</span>
        </router-link>
        <router-link to="/users" class="nav-item" active-class="active" :title="isCollapsed ? 'Kelola User & Penilai' : ''">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
          <span v-if="!isCollapsed">Kelola User & Penilai</span>
        </router-link>
        <router-link to="/reviews" class="nav-item" active-class="active" :title="isCollapsed ? 'Penilaian Project' : ''">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/></svg>
          <span v-if="!isCollapsed">Penilaian Project</span>
        </router-link>
        <router-link to="/reviews/history" class="nav-item" active-class="active" :title="isCollapsed ? 'Riwayat Audit' : ''">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          <span v-if="!isCollapsed">Riwayat Audit</span>
        </router-link>
      </template>

      <!-- Pemohon Section -->
      <template v-if="auth.isPemohon">
        <div class="nav-label" v-if="!isCollapsed">PEMOHON</div>
        <router-link to="/projects" class="nav-item" active-class="active" :title="isCollapsed ? 'Daftar Project' : ''">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
          <span v-if="!isCollapsed">Daftar Project</span>
        </router-link>
        <router-link to="/projects/create" class="nav-item" active-class="active" :title="isCollapsed ? 'Buat Pengajuan' : ''">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg>
          <span v-if="!isCollapsed">Buat Pengajuan</span>
        </router-link>
      </template>

      <!-- Penilai Section -->
      <template v-if="auth.isPenilai">
        <div class="nav-label" v-if="!isCollapsed">PENILAI</div>
        <router-link to="/reviews" class="nav-item" active-class="active" :title="isCollapsed ? 'Penilaian' : ''">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 11l3 3L22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/></svg>
          <span v-if="!isCollapsed">Penilaian</span>
        </router-link>
        <router-link to="/reviews/history" class="nav-item" active-class="active" :title="isCollapsed ? 'Riwayat Penilaian' : ''">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          <span v-if="!isCollapsed">Riwayat Penilaian</span>
        </router-link>
      </template>
    </nav>

    <div class="sidebar-footer">
      <div class="user-info" v-if="!isCollapsed">
        <div class="user-avatar">{{ auth.userName?.charAt(0)?.toUpperCase() }}</div>
        <div class="user-details">
          <span class="user-name">{{ auth.userName }}</span>
          <span class="user-role">{{ getRoleLabel(auth.userRole) }}</span>
        </div>
      </div>
      <button class="nav-item logout-btn" @click="auth.logout()" :title="isCollapsed ? 'Logout' : ''">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
        <span v-if="!isCollapsed">Logout</span>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const isCollapsed = ref(false)

const getRoleLabel = (role) => {
  if (role === 'admin') return 'Administrator'
  if (role === 'penilai') return 'Penilai / Verifikator'
  return 'Pemohon Dokumen'
}
</script>

<style scoped>
.sidebar {
  width: 260px;
  min-height: 100vh;
  background: rgba(15, 23, 42, 0.95);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  transition: width 0.2s ease;
  position: fixed;
  left: 0; top: 0; bottom: 0;
  z-index: 40;
}
.sidebar.collapsed { width: 68px; }

.sidebar-header {
  padding: 1.25rem 1rem;
  display: flex; align-items: center; justify-content: space-between;
  border-bottom: 1px solid var(--border-color);
  min-height: 70px;
}
.brand { display: flex; align-items: center; gap: 0.7rem; }
.brand-icon {
  width: 38px; height: 38px; border-radius: 10px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.2), rgba(79, 70, 229, 0.4));
  border: 1px solid rgba(99, 102, 241, 0.4);
  display: flex; align-items: center; justify-content: center;
  color: var(--accent-primary); flex-shrink: 0;
}
.brand-title { font-size: 0.95rem; font-weight: 800; color: var(--text-main); display: block; }
.brand-subtitle { font-size: 0.7rem; color: var(--text-muted); }

.collapse-btn {
  background: transparent; border: none; color: var(--text-muted);
  cursor: pointer; padding: 0.4rem; border-radius: 6px;
  display: flex; align-items: center;
}
.collapse-btn:hover { background: rgba(255,255,255,0.05); color: var(--text-main); }

.sidebar-nav { flex: 1; padding: 0.75rem 0.5rem; overflow-y: auto; }

.nav-label {
  font-size: 0.68rem; font-weight: 700; color: var(--text-subtle);
  text-transform: uppercase; letter-spacing: 0.06em;
  padding: 0.85rem 0.75rem 0.35rem; margin-top: 0.25rem;
}

.nav-item {
  display: flex; align-items: center; gap: 0.75rem;
  padding: 0.6rem 0.75rem; border-radius: 8px;
  color: var(--text-muted); text-decoration: none;
  font-size: 0.88rem; font-weight: 500;
  transition: all 0.15s ease; margin-bottom: 2px; cursor: pointer;
  border: none; background: none; width: 100%; text-align: left;
  font-family: var(--font-main);
}
.nav-item:hover { background: rgba(255,255,255,0.05); color: var(--text-main); }
.nav-item.active { background: rgba(99,102,241,0.15); color: #818cf8; }
.nav-item.active svg { color: #818cf8; }

.sidebar-footer {
  padding: 0.75rem; border-top: 1px solid var(--border-color);
}
.user-info {
  display: flex; align-items: center; gap: 0.65rem;
  padding: 0.5rem 0.5rem 0.65rem;
}
.user-avatar {
  width: 34px; height: 34px; border-radius: 8px;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  display: flex; align-items: center; justify-content: center;
  color: white; font-size: 0.85rem; font-weight: 700; flex-shrink: 0;
}
.user-name { font-size: 0.85rem; font-weight: 600; color: var(--text-main); display: block; }
.user-role { font-size: 0.72rem; color: var(--text-muted); }

.logout-btn { color: #fb7185; }
.logout-btn:hover { background: rgba(244,63,94,0.1); }

.collapsed .sidebar-nav .nav-item { justify-content: center; padding: 0.65rem; }
.collapsed .sidebar-footer { padding: 0.5rem; }
.collapsed .sidebar-footer .nav-item { justify-content: center; padding: 0.65rem; }
</style>
