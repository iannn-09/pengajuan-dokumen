# Technical Test Programmer — Sistem Pengajuan Dokumen Kelayakan

Aplikasi Full Stack **Sistem Pengajuan Dokumen Kelayakan** berbasis **Go (Gin Framework, GORM ORM)** pada bagian Backend API, **Vue 3 (Vite, Pinia, Chart.js)** pada bagian Frontend, serta **PostgreSQL** sebagai sistem manajemen basis data utama.

Proyek ini dirancang untuk menangani ribuan hingga ratusan ribu data permohonan dokumen beserta riwayat audit penilaiannya dengan performa tinggi, validasi keandalan data, serta menerapkan standar keamanan dan *clean code*.

---

## 🚀 Fitur Utama

### 1. Multi-Role Authentication & Access Control (JWT & Bcrypt)
- **Role Pemohon Dokumen**:
  - Registrasi & Login dengan enkripsi password `bcrypt`.
  - Membuat, mengedit, dan mengelola draft permohonan dokumen.
  - Mengunggah dokumen lampiran (PDF, PNG, JPG, DOCX) dengan penamaan aman berformat UUID v4 dan validasi ukuran file max 10MB.
  - Mengirim permohonan untuk verifikasi penilaian.
  - Memantau status permohonan real-time (`DRAFT`, `SUBMITTED`, `UNDER_REVIEW`, `REVISION`, `APPROVED`, `REJECTED`) serta membaca catatan revisi dari penilai.
- **Role Penilai**:
  - Peninjauan antrean permohonan masuk (`SUBMITTED`, `UNDER_REVIEW`).
  - Mengunduh & memverifikasi dokumen lampiran pemohon.
  - Mengambil keputusan penilaian: **Setujui (APPROVED)**, **Minta Revisi (REVISION)**, atau **Tolak (REJECTED)** disertai catatan evaluasi.
  - Riwayat audit (Audit Log) dan histori penilaian yang tersimpan secara kronologis.

### 2. High-Performance Dashboard & Visual Analytics
- Interaktif Chart.js (`Vue-ChartJS`) untuk memvisualisasikan tren pengajuan dokumen bulanan dan distribusi status project permohonan.
- Ringkasan statistik real-time (Total Project, Draft, Perlu Revisi, Disetujui, Ditolak).

### 3. High-Speed CLI Data Seeder (10.000 Projects + 2.000 Users)
- Dilengkapi tool CLI seeder berkecepatan tinggi yang mampu mengisi 1.000 akun Pemohon, 1.000 akun Penilai, dan 10.000 data Project Permohonan secara otomatis hanya dalam hitungan detik.

---

## 🛠️ Stack Teknologi & Dependensi

- **Backend**: Go 1.25.3, Gin Web Framework, GORM v1.31.2, Golang-JWT v5, Bcrypt.
- **Frontend**: Vue 3, Vite, Pinia (State Management), Vue Router, Chart.js, Axios.
- **Database**: PostgreSQL 18.
- **DevOps & Tooling**: Docker, Docker Compose, Air (Hot-Reloading Go).

---

## 📂 Struktur Direktori Proyek

```text
pengajuan-dokumen/
├── go.work                     # Workspace Go Multi-Module (backend & wa-bot)
├── docker-compose.yml          # Containerization (PostgreSQL, Backend API, Frontend)
├── README.md                   # Dokumentasi Utama Proyek
├── backend/
│   ├── cmd/
│   │   ├── api/main.go         # Entry Point Server API Backend (Auto Migration)
│   │   └── seed/main.go        # CLI Tool Data Seeder (10.000 Projects + 2.000 Users)
│   ├── internal/
│   │   ├── config/             # Koneksi Database & Environment Loader
│   │   ├── controllers/        # REST Controllers (Auth, Project, Review, Upload, Dashboard)
│   │   ├── middleware/         # JWT Middleware & Role-Based Access Control
│   │   ├── models/             # Schema Model GORM (User, Project, Document, ReviewHistory)
│   │   ├── routes/             # Gin Routes & Config CORS
│   │   └── seeder/             # Logika High-Speed Batch Seeding
│   ├── uploads/                # Direktori Penyimpanan Berkas Lampiran
│   ├── .env.example            # Template Konfigurasi Environment Backend
│   └── Dockerfile              # Dockerfile Multi-Stage Build Backend Go
├── wa-bot/
│   ├── main.go                 # Entry Point WhatsApp Bot API
│   ├── routes.go               # Routing API Session & Message
│   ├── controllers.go          # Handler Session, Send Message, Bulk Message
│   ├── session_manager.go      # Manajemen sesi WhatsApp
│   ├── types.go                # Tipe request JSON untuk API bot
│   ├── Dockerfile              # Dockerfile service WhatsApp bot
│   └── docker-compose.yml      # Compose terpisah untuk menjalankan WA bot
└── frontend/
    ├── src/
    │   ├── assets/             # Global CSS & Design System
    │   ├── components/         # Reusable UI (Sidebar, StatusBadge, Pagination, Timeline, Chart)
    │   ├── router/             # Vue Router & Navigation Guards (Auth/Role Check)
    │   ├── services/           # Axios HTTP Client
    │   ├── stores/             # Pinia Store (Authentication & State Management)
    │   └── views/              # Halaman Aplikasi (Login, Register, Dashboard, Projects, Reviews)
    ├── Dockerfile              # Dockerfile Nginx Static Build Frontend
    └── package.json            # NPM Dependencies & Scripts
```

---

## 📋 Prasyarat Sistem (Prerequisites)

Sebelum menjalankan proyek, pastikan perangkat Anda telah terpasang:
- **Go**: v1.24 atau v1.25 ([Download Go](https://go.dev/dl/))
- **Node.js**: v18.x atau v20.x ([Download Node.js](https://nodejs.org/))
- **PostgreSQL**: v15 atau v16 (jika menjalankan lokal tanpa Docker)
- **Docker & Docker Compose**: ([Download Docker Desktop](https://www.docker.com/products/docker-desktop/)) *(Opsional jika menggunakan Docker)*

---

## 💻 Panduan Menjalankan Proyek

Anda dapat memilih salah satu dari 2 opsi cara menjalankan proyek di bawah ini.

### 🌟 OPSI 1: Menjalankan Menggunakan Docker Compose (Sangat Direkomendasikan)

Metode ini paling praktis karena seluruh *environment* (Database PostgreSQL, Backend API, dan Frontend Vue) akan disiapkan dan dihubungkan secara otomatis.

1. **Jalankan Seluruh Container**:
   Buka terminal di root directory proyek (`pengajuan-dokumen/`):
   ```bash
   docker-compose up --build -d
   ```

2. **Jalankan Data Seeder (10.000 Projects + 2.000 Users)**:
   Setelah seluruh container berjalan, isi database awal dengan menjalankan perintah berikut di terminal:
   ```bash
   docker exec -it pengajuan_backend go run ./cmd/seed/main.go
   ```

3. **Akses Aplikasi**:
   - **Frontend App**: [http://localhost:5173](http://localhost:5173)
   - **Backend API Health Check**: [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)

---

### 🛠️ OPSI 2: Menjalankan Secara Manual (Local Development)

Jika Anda ingin menjalankan proyek secara lokal untuk kebutuhan pengujian atau pengembangan:

#### 1. Persiapan Database PostgreSQL
- Buat database baru bernama `dokumen` di PostgreSQL lokal Anda.
- **Atau** Anda bisa menggunakan container PostgreSQL dari Docker Compose:
  ```bash
  docker-compose up -d postgres
  ```

#### 2. Konfigurasi & Jalankan Backend Go dengan Air
Buka terminal pada folder `backend/`:
```bash
cd backend

# Salin file environment (.env)
# Windows (PowerShell):
Copy-Item .env.example .env
# Linux / macOS:
cp .env.example .env
```

Pastikan isi file `backend/.env` sesuai dengan konfigurasi PostgreSQL lokal Anda(untuk secret key sengaja saya beri agar tinggal memakai saja seharusnya itu tidak di tampilkan):
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=dokumen
DB_SSLMODE=disable
CORS_ALLOW_ORIGIN=http://localhost:5173
JWT_SECRET=rahasia-super-aman-pengajuan-dokumen-2026

WHATSAPP_NOTIFIER_ENABLED=true
WHATSAPP_API_URL=
WHATSAPP_API_TOKEN=
WHATSAPP_DEFAULT_COUNTRY_CODE=62
WHATSAPP_DEFAULT_DELAY=0
WHATSAPP_SENDER_NAME="${APP_NAME}"

WHATSAPP_BOT_URL=http://127.0.0.1:5173
WHATSAPP_BOT_SECRET=8f4d2a9c1b7e6f0a3d8c5e9b2f1a7c4d6e8f9a0b3c5d7e1f2a4b6c8d9e0f1a3

GEMINI_API_KEY=AQ.Ab8RN6IGG118QrDccf54sCengB-NHtGMRSu6zVUqFEpk2KeGuw
```

Selanjutnya, install dependensi dan jalankan server backend dengan Air:
```bash
# Download dependensi Go
go mod tidy

# Jalankan server backend dengan Air (Auto Migration akan berjalan otomatis)
air
```
*(Backend akan berjalan di `http://localhost:8080`)*.

#### 3. Jalankan CLI Data Seeder
Buka terminal baru di folder `backend/` untuk mengisi 10.000 data permohonan dan 2.000 user:
```bash
cd backend
go run ./cmd/seed/main.go
```

#### 4. Menjalankan Frontend Vue 3
Buka terminal pada folder `frontend/`:
```bash
cd frontend

# Install dependensi npm
npm install

# Jalankan server development Vite
npm run dev
```

Akses aplikasi frontend melalui browser pada alamat: **`http://localhost:5173`**

#### 5. Menjalankan WhatsApp Bot dengan Air
Service `wa-bot` berjalan terpisah dari backend utama dan dipakai untuk pengelolaan sesi WhatsApp, kirim pesan, serta bulk message. Service ini juga sudah disiapkan untuk dijalankan dengan Air melalui file [wa-bot/.air.toml](wa-bot/.air.toml).

1. Buka terminal pada folder `wa-bot/`:
   ```bash
   cd wa-bot
   ```

2. Siapkan environment file `.env` di folder `wa-bot/` dengan minimal isi seperti berikut:
   ```env
   API_SECRET=8f4d2a9c1b7e6f0a3d8c5e9b2f1a7c4d6e8f9a0b3c5d7e1f2a4b6c8d9e0f1a3
   WEBHOOK_URL=http://127.0.0.1:8000/api/whatsapp-webhook
   ```

3. Jalankan service bot dengan Air:
   ```bash
   air
   ```

4. Service akan tersedia di port `3000`.

5. Buka endpoint status untuk memastikan bot hidup:
   ```bash
   curl http://localhost:3000/
   ```

6. Gunakan endpoint session dan message berikut sesuai kebutuhan:
   - `POST /sessions` untuk membuat session baru dan menampilkan QR code.
   - `GET /sessions` untuk melihat daftar session.
   - `GET /sessions/{id}` untuk melihat detail session.
   - `PATCH /sessions/{id}` untuk mengganti nama session.
   - `POST /sessions/{id}/logout` untuk logout session.
   - `DELETE /sessions/{id}` untuk menghapus session.
   - `POST /send-message` untuk mengirim pesan satu penerima.
   - `POST /send-bulk-same-message` untuk mengirim pesan yang sama ke banyak nomor.
   - `POST /send-bulk-different-messages` untuk mengirim pesan berbeda ke banyak nomor.
   - `GET /groups` untuk mengambil daftar grup WhatsApp yang tersedia.

Catatan:
- Semua request write pada API bot mengirim field `secret` yang harus sama dengan `API_SECRET`.
- Session bot disimpan otomatis di folder `wa-bot/session/`.
- Bila Anda memakai Docker, jalankan compose di folder `wa-bot/` secara terpisah karena service ini tidak ada di `docker-compose.yml` root.

---

## 🔑 Akun Pengujian Default (Default Test Credentials)

Setelah menjalankan perintah **Data Seeder** di atas, Anda dapat langsung login menggunakan salah satu akun default berikut (atau membuat akun baru via menu **Register**):

| Role / Peran | Email | Password | Hak Akses |
| :--- | :--- | :--- | :--- |
| **Administrator** | `admin@kelayakan.id` | `password123` | Akses Penuh Sistem & Management User |
| **Pemohon Dokumen** | `pemohon1@kelayakan.id` | `password123` | Buat & Kelola Permohonan Dokumen |
| **Pemohon Dokumen** | `pemohon2@kelayakan.id` | `password123` | Buat & Kelola Permohonan Dokumen |
| **Penilai** | `penilai1@kelayakan.id` | `password123` | Review, Disetujui, Revisi, Tolak Dokumen |
| **Penilai** | `penilai2@kelayakan.id` | `password123` | Review, Disetujui, Revisi, Tolak Dokumen |

---

## 🛡️ Aspek Keamanan & Standardisasi Kode

- **Password Hashing**: Menggunakan `golang.org/x/crypto/bcrypt` dengan standar *cost factor*.
- **Keamanan Unggah File (Path Traversal & Shell Upload Protection)**: File lampiran disimpan menggunakan penamaan acak **UUID v4**, serta dilakukan pembatasan ekstensi file (`.pdf`, `.png`, `.jpg`, `.docx`) dan batas ukuran maksimum 10MB.
- **SQL Injection Prevention**: Seluruh akses database menggunakan query terparametrisasi via GORM ORM.
- **JWT & Role-Based Access Control (RBAC)**: Endpoint API Backend dan Navigasi Frontend dilindungi oleh middleware JWT yang memverifikasi kecocokan peran (`pemohon` vs `penilai`).
- **Auto Database Migration**: Skema tabel database diperbarui secara otomatis saat aplikasi diawali tanpa memerlukan eksekusi file SQL manual.

---

## ❓ Solusi Masalah Umum (Troubleshooting)

1. **Gagal Koneksi Database (`connection refused`)**:
   - Pastikan service PostgreSQL sudah berjalan di port `5432`.
   - Pastikan informasi `DB_USER`, `DB_PASSWORD`, dan `DB_NAME` di file `backend/.env` sudah sesuai dengan database Anda.

2. **Gariskait Merah pada Import Go di VS Code**:
   - Jika VS Code menampilkan garis merah pada baris import Go, buka Command Palette (`Ctrl+Shift+P`), lalu pilih **`Go: Restart Language Server`**.

3. **Port 8080 atau 5173 Sudah Terpakai**:
   - Anda dapat mengubah nilai `PORT` pada file `backend/.env` atau menyesuaikan `ports` pada file `docker-compose.yml`.