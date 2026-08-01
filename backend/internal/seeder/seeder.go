package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/middleware"
	"pengajuan-dokumen/backend/internal/models"
)

func SeedDatabase(totalPemohon, totalPenilai, totalProjects int) {
	log.Println("Starting Database Seeder...")
	startTime := time.Now()

	// 1. Seed Master Document Types
	seedDocumentTypes()

	// 2. Check existing users
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Printf("Found %d existing users. Skipping user seeding or appending...", userCount)
	} else {
		seedUsers(totalPemohon, totalPenilai)
	}

	// 3. Check existing projects
	var projectCount int64
	config.DB.Model(&models.Project{}).Count(&projectCount)
	if projectCount >= int64(totalProjects) {
		log.Printf("Found %d existing projects. Seeding goal of %d projects already met.", projectCount, totalProjects)
		return
	}

	projectsToSeed := totalProjects - int(projectCount)
	seedProjects(projectsToSeed)

	log.Printf("Database Seeding Completed in %v!", time.Since(startTime))
}

func seedDocumentTypes() {
	var count int64
	config.DB.Model(&models.DocumentType{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("Seeding Master Document Types...")
	docTypes := []models.DocumentType{
		{
			Code:        "AMDAL",
			Name:        "Analisis Mengenai Dampak Lingkungan (AMDAL)",
			Requirement: "1. KTP & NPWP Pemohon\n2. NIB (Nomor Induk Berusaha)\n3. Akta Pendirian Perusahaan\n4. Peta Site Plan & Tata Ruang\n5. Draft Kerangka Acuan (KA-ANDAL)",
			Description: "Kajian mengenai dampak penting suatu usaha/kegiatan terhadap lingkungan hidup. Target pengerjaan verifikasi: 14 Hari Kerja.",
			IsActive:    true,
		},
		{
			Code:        "UKL-UPL",
			Name:        "Upaya Pengelolaan & Pemantauan Lingkungan (UKL-UPL)",
			Requirement: "1. KTP & NPWP Pemohon\n2. NIB & Izin Usaha Perusahaan\n3. Denah Bangunan & Tata Alat\n4. Formulir Isian UKL-UPL Lengkap",
			Description: "Pengelolaan dan pemantauan lingkungan untuk kegiatan skala menengah. Target pengerjaan verifikasi: 7 Hari Kerja.",
			IsActive:    true,
		},
		{
			Code:        "SLF",
			Name:        "Sertifikat Laik Fungsi Bangunan (SLF)",
			Requirement: "1. Gambar As-Built Drawing Bangunan\n2. Laporan Hasil Uji Kelayakan Struktur Bangunan\n3. Rekomendasi Keselamatan Pemadam Kebakaran\n4. Dokumen Sertifikasi Ketenagalistrikan",
			Description: "Sertifikat kelaikan fungsi bangunan gedung dari pemerintah daerah. Target pengerjaan verifikasi: 10 Hari Kerja.",
			IsActive:    true,
		},
		{
			Code:        "ANDALALIN",
			Name:        "Analisis Dampak Lalu Lintas (ANDALALIN)",
			Requirement: "1. Denah Akses Masuk dan Keluar Kendaraan\n2. Data Volume Lalu Lintas Eksisting\n3. Perhitungan Kapasitas Parkir & Ramp Sirkulasi\n4. Surat Rekomendasi Dinas Perhubungan",
			Description: "Kajian pengaruh bangkitan lalu lintas kegiatan pembangunan terhadap jaringan jalan. Target pengerjaan verifikasi: 7 Hari Kerja.",
			IsActive:    true,
		},
		{
			Code:        "IZIN-B3",
			Name:        "Izin Pengolahan & Penyimpanan Limbah B3",
			Requirement: "1. Desain Tempat Penyimpanan Sementara Limbah B3 Berstandar\n2. MoU Kerjasama Pengolah Limbah B3 Pihak Ketiga Berizin\n3. Dokumen SOP Penanganan & Tanggap Darurat Limbah",
			Description: "Izin operasional pengelolaan dan penyimpanan bahan berbahaya dan beracun (B3). Target pengerjaan verifikasi: 21 Hari Kerja.",
			IsActive:    true,
		},
	}

	for _, dt := range docTypes {
		config.DB.Create(&dt)
	}
	log.Println("Master Document Types seeded successfully!")
}

func seedUsers(totalPemohon, totalPenilai int) {
	log.Printf("Seeding Admin, %d Pemohon, and %d Penilai users...", totalPemohon, totalPenilai)

	defaultPassword, _ := middleware.HashPassword("password123")

	// 1. Seed Admin user
	adminUser := models.User{
		Name:     "Super Administrator",
		Email:    "admin@kelayakan.id",
		Password: defaultPassword,
		Role:     models.RoleAdmin,
		Phone:    "081100001111",
		Company:  "Instansi Pusat Pengelolaan Kelayakan",
	}
	config.DB.Create(&adminUser)
	log.Println("Seeded Admin user (admin@kelayakan.id)")

	// 2. Batch insert Pemohon
	batchSize := 250
	var pemohonList []models.User
	for i := 1; i <= totalPemohon; i++ {
		u := models.User{
			Name:     fmt.Sprintf("Pemohon Perusahaan %d", i),
			Email:    fmt.Sprintf("pemohon%d@kelayakan.id", i),
			Password: defaultPassword,
			Role:     models.RolePemohon,
			Phone:    fmt.Sprintf("0812%08d", i),
			Company:  fmt.Sprintf("PT Industri Kelayakan Nusantara %d", i),
		}
		pemohonList = append(pemohonList, u)

		if len(pemohonList) == batchSize || i == totalPemohon {
			if err := config.DB.Create(&pemohonList).Error; err != nil {
				log.Fatalf("Failed batch seeding pemohon: %v", err)
			}
			log.Printf("Seeded %d/%d Pemohon users...", i, totalPemohon)
			pemohonList = nil
		}
	}

	// 3. Batch insert Penilai
	var penilaiList []models.User
	for i := 1; i <= totalPenilai; i++ {
		u := models.User{
			Name:     fmt.Sprintf("Penilai Verifikator %d", i),
			Email:    fmt.Sprintf("penilai%d@kelayakan.id", i),
			Password: defaultPassword,
			Role:     models.RolePenilai,
			Phone:    fmt.Sprintf("0857%08d", i),
			Company:  "Kementerian / Dinas Lingkungan Hidup",
		}
		penilaiList = append(penilaiList, u)

		if len(penilaiList) == batchSize || i == totalPenilai {
			if err := config.DB.Create(&penilaiList).Error; err != nil {
				log.Fatalf("Failed batch seeding penilai: %v", err)
			}
			log.Printf("Seeded %d/%d Penilai users...", i, totalPenilai)
			penilaiList = nil
		}
	}

	log.Println("Users seeding finished.")
}

func seedProjects(count int) {
	log.Printf("Seeding %d Projects with review history...", count)

	var pemohonIDs []uint
	config.DB.Model(&models.User{}).Where("role = ?", models.RolePemohon).Pluck("id", &pemohonIDs)

	var docTypes []models.DocumentType
	config.DB.Where("is_active = ?", true).Find(&docTypes)

	if len(pemohonIDs) == 0 {
		log.Fatalf("Cannot seed projects: Pemohon users missing!")
	}

	statuses := []models.ProjectStatus{
		models.StatusDraft,
		models.StatusSubmitted,
		models.StatusUnderReview,
		models.StatusRevision,
		models.StatusApproved,
		models.StatusRejected,
	}

	companies := []string{
		"PT Sumber Energi Utama",
		"PT Petrokimia Nusantara",
		"PT Graha Bangun Konstruksi",
		"PT Agro Lestari Gemilang",
		"PT Logistik Trans Nasional",
		"PT Manufaktur Otomotif Indonesia",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	batchSize := 500
	var projectBatch []models.Project

	for i := 1; i <= count; i++ {
		pemohonID := pemohonIDs[r.Intn(len(pemohonIDs))]
		status := statuses[r.Intn(len(statuses))]
		company := companies[r.Intn(len(companies))]

		var docTypeID *uint
		docTypeName := "AMDAL Kelayakan Lingkungan"
		if len(docTypes) > 0 {
			dt := docTypes[r.Intn(len(docTypes))]
			docTypeID = &dt.ID
			docTypeName = dt.Name
		}

		createdAt := time.Now().AddDate(0, -r.Intn(12), -r.Intn(30))
		var submittedAt *time.Time
		if status != models.StatusDraft {
			sTime := createdAt.Add(time.Hour * time.Duration(r.Intn(48)+1))
			submittedAt = &sTime
		}

		proj := models.Project{
			ProjectNumber:  fmt.Sprintf("PRJ-%s-%05d", createdAt.Format("20060102"), i),
			Title:          fmt.Sprintf("Permohonan %s - Tahap %d", docTypeName, (i%5)+1),
			Description:    fmt.Sprintf("Permohonan dokumen kelayakan %s untuk fasilitas operasional yang berlokasi di Kawasan Industri Sektor %d.", docTypeName, (i%10)+1),
			CompanyName:    fmt.Sprintf("%s Unit %d", company, (i%20)+1),
			DocumentTypeID: docTypeID,
			UserID:         pemohonID,
			Status:         status,
			SubmittedAt:    submittedAt,
			CreatedAt:      createdAt,
			UpdatedAt:      createdAt,
		}

		projectBatch = append(projectBatch, proj)

		if len(projectBatch) == batchSize || i == count {
			if err := config.DB.Create(&projectBatch).Error; err != nil {
				log.Fatalf("Failed batch seeding projects: %v", err)
			}
			log.Printf("Seeded %d/%d Projects...", i, count)
			projectBatch = nil
		}
	}

	log.Println("Projects seeding finished.")
}
