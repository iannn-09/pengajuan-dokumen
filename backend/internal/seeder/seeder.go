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

	// 1. Check existing users
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Printf("Found %d existing users. Skipping user seeding or appending...", userCount)
	} else {
		seedUsers(totalPemohon, totalPenilai)
	}

	// 2. Check existing projects
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

func seedUsers(totalPemohon, totalPenilai int) {
	log.Printf("Seeding %d Pemohon and %d Penilai users...", totalPemohon, totalPenilai)

	defaultPassword, _ := middleware.HashPassword("password123")

	// Batch insert Pemohon
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

	// Batch insert Penilai
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

	var penilaiIDs []uint
	config.DB.Model(&models.User{}).Where("role = ?", models.RolePenilai).Pluck("id", &penilaiIDs)

	if len(pemohonIDs) == 0 || len(penilaiIDs) == 0 {
		log.Fatalf("Cannot seed projects: Pemohon or Penilai users missing!")
	}

	statuses := []models.ProjectStatus{
		models.StatusDraft,
		models.StatusSubmitted,
		models.StatusUnderReview,
		models.StatusRevision,
		models.StatusApproved,
		models.StatusRejected,
	}

	types := []string{
		"AMDAL Kelayakan Lingkungan",
		"UKL-UPL Dokumen Lingkungan",
		"Izin Lokasi & Tata Ruang",
		"Sertifikat Laik Fungsi (SLF)",
		"Analisis Dampak Lalu Lintas (Andalalin)",
		"Izin Pengolahan Limbah B3",
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
		docType := types[r.Intn(len(types))]
		company := companies[r.Intn(len(companies))]

		createdAt := time.Now().AddDate(0, -r.Intn(12), -r.Intn(30))
		var submittedAt *time.Time
		if status != models.StatusDraft {
			sTime := createdAt.Add(time.Hour * time.Duration(r.Intn(48)+1))
			submittedAt = &sTime
		}

		proj := models.Project{
			ProjectNumber: fmt.Sprintf("PRJ-%s-%05d", createdAt.Format("20060102"), i),
			Title:         fmt.Sprintf("Permohonan %s - Tahap %d", docType, (i%5)+1),
			Description:   fmt.Sprintf("Permohonan dokumen kelayakan %s untuk fasilitas operasional yang berlokasi di Kawasan Industri Sektor %d.", docType, (i%10)+1),
			CompanyName:   fmt.Sprintf("%s Unit %d", company, (i%20)+1),
			UserID:        pemohonID,
			Status:        status,
			SubmittedAt:   submittedAt,
			CreatedAt:     createdAt,
			UpdatedAt:     createdAt,
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
