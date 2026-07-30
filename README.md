# Technical Test Programmer — Sistem Pengajuan Dokumen Kelayakan

Aplikasi Full Stack **Sistem Pengajuan Dokumen Kelayakan** berbasis **Go (Gin, GORM, Air)** di sisi Backend, **Vue 3 (Vite, Pinia, Chart.js)** di sisi Frontend, dan **PostgreSQL** sebagai database utama.

Proyek ini dibangun untuk menangani ratusan ribu hingga jutaan data permohonan dokumen beserta riwayat audit penilaiannya dengan performa tinggi, keandalan query, serta keamanan sesuai prinsip clean code.

---

## 🚀 Fitur Utama Sesuai Studi Kasus

### 1. Multi-Role Authentication & Access Control (JWT)
- **Role Pemohon Dokumen**:
  - Pendaftaran & Login dengan enkripsi bcrypt.
  - Membuat & mengelola draft permohonan dokumen.
  - Unggah dokumen lampiran (PDF, JPG, PNG, DOCX) dengan nama tersimpan UUID & validasi file 10MB.
  - Mengirim permohonan untuk proses verifikasi penilaian.
  - Memantau status real-time & riwayat revisi dari penilai.
- **Role Penilai / Verifikator**:
  - Peninjauan antrean permohonan masuk (`SUBMITTED`, `UNDER_REVIEW`).
  - Mengunduh & memverifikasi dokumen lampiran.
  - Mengambil keputusan penilaian: **Setujui (APPROVED)**, **Minta Revisi (REVISION)**, atau **Tolak (REJECTED)** disertai catatan evaluasi.
  - Log audit & histori penilaian komprehensif.

### 2. High-Performance Dashboard & Analytics
- **Chart.js & Vue-ChartJS**: Visualisasi grafis pengajuan bulanan & distribusi status project.
- Data ringkasan statistik (Total, Draft, Perlu Revisi, Disetujui, Ditolak) secara real-time.

### 3. Big Data Seeding CLI (10.000 Projects + 2.000 Users)
- Dilengkapi tool seeder CLI berkecepatan tinggi yang mampu mengisi 1.000 Pemohon, 1.000 Penilai, dan 10.000 data Project Permohonan secara batch dalam hitungan detik.

---

## 🛠️ Stack Teknologi

- **Backend**: Go (Golang 1.21+), Gin Web Framework, GORM ORM, Golang JWT v5, Bcrypt.
- **Frontend**: Vue 3, Vite, Pinia State Management, Vue Router, Chart.js & Vue-ChartJS, Axios.
- **Database**: PostgreSQL 15.
- **DevOps & Testing**: Docker, Docker Compose, Air (Hot Reloading), GitLab CI/CD Pipeline.

---

## 📂 Struktur Direktori Proyek

```text
pengajuan-dokumen/
├── .gitlab-ci.yml              # GitLab CI/CD Pipeline
├── docker-compose.yml          # Docker Compose (Postgres, Backend, Frontend)
├── README.md                   # Dokumentasi Resmi Proyek
├── backend/
│   ├── cmd/
│   │   ├── api/main.go         # Entry Point Backend API Server
│   │   └── seed/main.go        # CLI Tool Seeder (10.000 Project + 2.000 User)
│   ├── internal/
│   │   ├── config/             # DB & Environment Loader
│   │   ├── controllers/        # REST Controllers (Auth, Project, Review, Upload, Dashboard)
│   │   ├── middleware/         # JWT Auth & Role-Based Access Middleware
│   │   ├── models/             # GORM Models (User, Project, Document, ReviewHistory)
│   │   ├── routes/             # Gin Router Setup & CORS
│   │   └── seeder/             # Batch Seeder Logic
│   ├── uploads/                # Directory Penyimpanan File Terisolasi
│   ├── .env.example
│   ├── .air.toml
│   └── Dockerfile
└── frontend/
    ├── src/
    │   ├── assets/             # Global CSS & Design System
    │   ├── components/         # Reusable Components (Sidebar, StatusBadge, Pagination, Timeline)
    │   ├── router/             # Vue Router dengan Route Guards
    │   ├── services/           # Axios Client
    │   ├── stores/             # Pinia Auth Store
    │   ├── views/              # Pages (Login, Register, Dashboard, Projects, Reviews)
    │   ├── App.vue
    │   └── main.js
    ├── Dockerfile
    └── package.json
```

---

## 💻 Cara Menjalankan Proyek

### Opsi 1: Menggunakan Docker Compose (Direkomendasikan)

Pastikan Docker Desktop aktif di komputer Anda:

```bash
# 1. Build dan jalankan seluruh container (Postgres, Backend, Frontend)
docker-compose up --build -d

# 2. Buka aplikasi di browser
# Frontend : http://localhost:5173
# Backend  : http://localhost:8080/api/v1/health
```

---

### Opsi 2: Menjalankan Secara Manual (Local Development)

#### 1. Persiapan Database (PostgreSQL)
Jalankan PostgreSQL lokal atau gunakan container postgres dari Docker:
```bash
docker-compose up -d postgres
```
Pastikan file `backend/.env` telah terkonfigurasi dengan credential database lokal Anda:
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=dokumen
DB_SSLMODE=disable
CORS_ALLOW_ORIGIN=http://localhost:5173
JWT_SECRET=super-secret-key-change-in-production
```

#### 2. Jalankan Backend Go
Buka terminal pada folder `backend/`:
```bash
cd backend

# Download dependensi
go mod tidy

# Jalankan via Air (Hot-Reloading)
air

# Atau jalankan langsung tanpa Air:
go run ./cmd/api/main.go
```

#### 3. Jalankan Seeder Data (10.000 Project + 2.000 User)
Untuk mengisi database dengan 10.000 data project dan 2.000 user penguji:
```bash
cd backend
go run ./cmd/seed/main.go
```
*Kredensial Akun Pengujian Default setelah Seeding:*
- **Pemohon**: `pemohon1@kelayakan.id` s/d `pemohon1000@kelayakan.id` (Password: `password123`)
- **Penilai**: `penilai1@kelayakan.id` s/d `penilai1000@kelayakan.id` (Password: `password123`)

#### 4. Jalankan Frontend Vue 3
Buka terminal pada folder `frontend/`:
```bash
cd frontend

# Install dependensi npm
npm install

# Jalankan server pengembangan Vite
npm run dev
```
Akses aplikasi melalui browser: **`http://localhost:5173`**

---

## 🔒 Keamanan & Prinsip Clean Code (Security & Best Practices)

- **Password Hashing**: Menggunakan `golang.org/x/crypto/bcrypt` dengan default cost.
- **Path Traversal Protection**: Nama file yang diunggah dikonversi ke UUID v4 unik di atas disk, dan direktori penyimpanan diisolasi dari root web server.
- **SQL Injection Prevention**: Seluruh query menggunakan parametrisasi ORM GORM.
- **Role-Based Access Control (RBAC)**: Endpoint backend dan rute frontend dilindungi middleware JWT sesuai perannya masing-masing.

---

## 📄 Lisensi & Hak Cipta
Dikembangkan untuk keperluan **Technical Test Programmer - Sistem Pengajuan Dokumen Kelayakan**.