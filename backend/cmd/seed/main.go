package main

import (
	"log"
	"time"

	"pengajuan-dokumen/backend/internal/config"
	"pengajuan-dokumen/backend/internal/models"
)

func main() {
	log.Println("Starting Standalone Database Seeder Tool...")

	// Connect Database
	config.ConnectDatabase()

	startTime := time.Now()

	// Seed Master Document Types
	seedDocumentTypes()

	log.Printf("Database Seeding Completed in %v!", time.Since(startTime))
}

func seedDocumentTypes() {
	var count int64
	config.DB.Model(&models.DocumentType{}).Count(&count)
	if count > 0 {
		log.Println("Master Document Types already exist. Skipping...")
		return
	}

	log.Println("Seeding Master Document Types...")

	docTypes := []models.DocumentType{
		{
			Code: "AMDAL",
			Name: "Analisis Mengenai Dampak Lingkungan (AMDAL)",
			Requirement: `<h2>Persyaratan Dokumen</h2>
<ul>
	<li>KTP &amp; NPWP Pemohon</li>
	<li>Nomor Induk Berusaha (NIB)</li>
	<li>Akta Pendirian Perusahaan</li>
	<li>Peta Site Plan &amp; Tata Ruang</li>
	<li>Draft Kerangka Acuan (KA-ANDAL)</li>
</ul>`,
			Description: `<p>Kajian mengenai dampak penting suatu usaha atau kegiatan terhadap lingkungan hidup.</p>
<p><strong>Target pengerjaan verifikasi:</strong> 14 Hari Kerja.</p>`,
			IsActive: true,
		},
		{
			Code: "UKL-UPL",
			Name: "Upaya Pengelolaan & Pemantauan Lingkungan (UKL-UPL)",
			Requirement: `<h2>Persyaratan Dokumen</h2>
<ul>
	<li>KTP &amp; NPWP Pemohon</li>
	<li>NIB &amp; Izin Usaha Perusahaan</li>
	<li>Denah Bangunan &amp; Tata Letak Alat</li>
	<li>Formulir Isian UKL-UPL Lengkap</li>
</ul>`,
			Description: `<p>Dokumen pengelolaan dan pemantauan lingkungan untuk kegiatan skala menengah.</p>
<p><strong>Target pengerjaan verifikasi:</strong> 7 Hari Kerja.</p>`,
			IsActive: true,
		},
		{
			Code: "SLF",
			Name: "Sertifikat Laik Fungsi Bangunan (SLF)",
			Requirement: `<h2>Persyaratan Dokumen</h2>
<ul>
	<li>As-Built Drawing Bangunan</li>
	<li>Laporan Uji Kelayakan Struktur</li>
	<li>Rekomendasi Keselamatan Kebakaran</li>
	<li>Sertifikat Instalasi Ketenagalistrikan</li>
</ul>`,
			Description: `<p>Sertifikat yang menyatakan bangunan telah memenuhi persyaratan kelaikan fungsi.</p>
<p><strong>Target pengerjaan verifikasi:</strong> 10 Hari Kerja.</p>`,
			IsActive: true,
		},
		{
			Code: "ANDALALIN",
			Name: "Analisis Dampak Lalu Lintas (ANDALALIN)",
			Requirement: `<h2>Persyaratan Dokumen</h2>
<ul>
	<li>Denah Akses Kendaraan</li>
	<li>Data Volume Lalu Lintas Eksisting</li>
	<li>Perhitungan Kapasitas Parkir</li>
	<li>Rekomendasi Dinas Perhubungan</li>
</ul>`,
			Description: `<p>Dokumen analisis dampak lalu lintas akibat pembangunan atau kegiatan tertentu.</p>
<p><strong>Target pengerjaan verifikasi:</strong> 7 Hari Kerja.</p>`,
			IsActive: true,
		},
		{
			Code: "IZIN-B3",
			Name: "Izin Pengolahan & Penyimpanan Limbah B3",
			Requirement: `<h2>Persyaratan Dokumen</h2>
<ul>
	<li>Desain TPS Limbah B3</li>
	<li>MoU dengan Pengolah Limbah Berizin</li>
	<li>SOP Penanganan Limbah B3</li>
	<li>Dokumen Tanggap Darurat</li>
</ul>`,
			Description: `<p>Izin operasional pengelolaan dan penyimpanan limbah bahan berbahaya dan beracun (B3).</p>
<p><strong>Target pengerjaan verifikasi:</strong> 21 Hari Kerja.</p>`,
			IsActive: true,
		},
	}

	if err := config.DB.Create(&docTypes).Error; err != nil {
		log.Fatalf("Failed to seed document types: %v", err)
	}

	log.Printf("Successfully seeded %d document types.", len(docTypes))
}