package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type AIChatMessage struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

// GenerateAIResponse generates a response using Google Gemini API or intelligent knowledge engine
func GenerateAIResponse(userMessage string, userRole string) string {
	apiKey := os.Getenv("GEMINI_API_KEY")

	// If Gemini API Key is configured in .env, call Gemini API
	if apiKey != "" {
		geminiResponse, err := callGeminiAPI(apiKey, userMessage, userRole)
		if err == nil && geminiResponse != "" {
			return geminiResponse
		}
		log.Printf("[AI Service] Gemini API call failed, falling back to smart Knowledge Engine. Err: %v\n", err)
	}

	// Smart Built-In AI Knowledge Engine for Pengajuan Dokumen Kelayakan
	return generateKnowledgeEngineResponse(userMessage, userRole)
}

func callGeminiAPI(apiKey, userMessage, userRole string) (string, error) {
	// Models to try in sequence
	models := []string{
		"gemini-1.5-flash",
		"gemini-2.0-flash",
		"gemini-pro",
	}

	systemPrompt := fmt.Sprintf(
		"Anda adalah Asisten AI Sistem Pengajuan Dokumen Kelayakan. Peran pengguna saat ini: %s.\n"+
			"Tugas Anda adalah membantu pengguna menjawab pertanyaan apapun seputar pengajuan dokumen kelayakan, persyarataan berkas (PDF/Gambar max 10MB), "+
			"status permohonan (Draft, Submitted, Under Review, Approved, Revision, Rejected), "+
			"pengeditan profil, pengunduhan rekapan Excel, serta notifikasi WhatsApp.\n"+
			"Berikan jawaban dalam bahasa Indonesia yang ramah, sopan, ringkas, terstruktur dengan format bullet point atau emoji bila perlu.",
		userRole,
	)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": systemPrompt + "\n\nPertanyaan Pengguna: " + userMessage},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, modelName := range models {
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", modelName, apiKey)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("[AI Service] Gemini model %s returned status %d: %s\n", modelName, resp.StatusCode, string(bodyBytes))
			continue
		}

		var geminiResult struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}

		if err := json.Unmarshal(bodyBytes, &geminiResult); err == nil {
			if len(geminiResult.Candidates) > 0 && len(geminiResult.Candidates[0].Content.Parts) > 0 {
				reply := strings.TrimSpace(geminiResult.Candidates[0].Content.Parts[0].Text)
				if reply != "" {
					return reply, nil
				}
			}
		}
	}

	return "", fmt.Errorf("all Gemini API model attempts failed")
}

func generateKnowledgeEngineResponse(msg string, role string) string {
	lowerMsg := strings.ToLower(msg)

	switch {
	case strings.Contains(lowerMsg, "admin") || strings.Contains(lowerMsg, "hubungi") || strings.Contains(lowerMsg, "kontak") || strings.Contains(lowerMsg, "cs") || strings.Contains(lowerMsg, "bantuan"):
		return "📞 *Cara Menghubungi Admin / Layanan Bantuan:*\n\n" +
			"1. **WhatsApp Support**: Anda dapat menghubungi Admin secara langsung melalui notifikasi pesan WhatsApp yang Anda terima di nomor Anda.\n" +
			"2. **Menu Profil**: Masuk ke menu *Profil Saya* untuk memperbarui informasi kontak Anda agar Admin dapat menghubungi Anda kembali.\n" +
			"3. **Tim Verifikator**: Jika ada catatan revisi pada dokumen Anda, nama Verifikator yang menilai berkas Anda juga tercantum pada detail pengajuan."

	case strings.Contains(lowerMsg, "password") || strings.Contains(lowerMsg, "sandi") || strings.Contains(lowerMsg, "lupa"):
		return "🔐 *Cara Mengubah Password / Kata Sandi:*\n\n" +
			"1. Klik foto profil Anda di pojok kiri bawah (Sidebar) atau buka menu *Profil Saya*.\n" +
			"2. Pada tab *Keamanan & Password*, masukkan password lama Anda.\n" +
			"3. Masukkan password baru (minimal 8 karakter) dan konfirmasi.\n" +
			"4. Klik *Simpan Password Baru*."

	case strings.Contains(lowerMsg, "foto") || strings.Contains(lowerMsg, "avatar") || strings.Contains(lowerMsg, "profil"):
		return "🖼️ *Cara Mengubah Foto Profil / Avatar:*\n\n" +
			"1. Buka halaman *Profil Saya* di sidebar navigasi.\n" +
			"2. Klik ikon *Kamera / Upload Foto* pada kartu profil Anda.\n" +
			"3. Pilih file gambar (JPG, PNG max 2MB) dari komputer/HP Anda.\n" +
			"4. Foto profil akan otomatis terpasang!"

	case strings.Contains(lowerMsg, "export") || strings.Contains(lowerMsg, "excel") || strings.Contains(lowerMsg, "laporan") || strings.Contains(lowerMsg, "csv"):
		return "📥 *Cara Mengunduh Rekapan Laporan Excel (CSV):*\n\n" +
			"1. Masuk ke halaman **Dashboard**.\n" +
			"2. Klik tombol **📥 Export Rekapan Laporan (Excel)** pada banner utama.\n" +
			"3. File spreadsheet (.CSV UTF-8) berisi seluruh daftar permohonan & log verifikasi akan otomatis terunduh dan dapat dibuka secara rapi di Microsoft Excel!"

	case strings.Contains(lowerMsg, "syarat") || strings.Contains(lowerMsg, "berkas") || strings.Contains(lowerMsg, "dokumen") || strings.Contains(lowerMsg, "pdf"):
		return "📌 *Persyaratan Dokumen Kelayakan:*\n\n" +
			"1. Berkas Wajib sesuai jenis pengajuan (misal: AMDAL, K3, Kelayakan Lingkungan).\n" +
			"2. Format file yang didukung: *PDF, JPG, PNG*.\n" +
			"3. Ukuran maksimal file: *10 MB* per dokumen.\n" +
			"4. Pastikan dokumen sudah ditandatangani dan jelas terbaca.\n\n" +
			"Anda dapat melihat opsi Jenis Dokumen pada form *Buat Pengajuan Baru*."

	case strings.Contains(lowerMsg, "status") || strings.Contains(lowerMsg, "arti"):
		return "📊 *Penjelasan Status Permohonan:*\n\n" +
			"• ⚪ *Draft*: Dokumen baru dibuat & belum dikirim ke Penilai.\n" +
			"• ⏳ *Submitted*: Dokumen telah dikirim & masuk antrean.\n" +
			"• 🔎 *Under Review*: Dokumen sedang diteliti oleh Verifikator.\n" +
			"• ✅ *Approved*: Permohonan disetujui & memenuhi kelayakan!\n" +
			"• ⚠️ *Revision*: Perlu perbaikan. Cek *Catatan Penilai* pada detail project.\n" +
			"• ❌ *Rejected*: Permohonan ditolak dengan alasan evaluasi tertentu."

	case strings.Contains(lowerMsg, "revisi") || strings.Contains(lowerMsg, "perbaiki"):
		return "📝 *Cara Menangani Status Revisi:*\n\n" +
			"1. Buka menu *Daftar Pengajuan* -> Klik permohonan yang berstatus *Perlu Revisi*.\n" +
			"2. Baca *Catatan Penilai* di bagian evaluasi.\n" +
			"3. Klik tombol *Edit Pengajuan* atau upload dokumen perbaikan baru.\n" +
			"4. Klik *Kirim Ulang Permohonan* agar Verifikator menilai kembali berkas Anda."

	case strings.Contains(lowerMsg, "wa") || strings.Contains(lowerMsg, "whatsapp") || strings.Contains(lowerMsg, "notif"):
		return "📱 *Notifikasi WhatsApp Automated:*\n\n" +
			"Sistem kami secara otomatis mengirimkan notifikasi WhatsApp ke nomor ponsel Anda setiap kali:\n" +
			"• Permohonan berhasil dikirim.\n" +
			"• Verifikator memperbarui status menjadi Disetujui, Perlu Revisi, atau Ditolak.\n\n" +
			"Pastikan nomor telepon di *Profil Saya* sudah benar (menggunakan format 08xxxx)."

	case strings.Contains(lowerMsg, "lama") || strings.Contains(lowerMsg, "waktu") || strings.Contains(lowerMsg, "proses"):
		return "⏱️ *Waktu Proses Verifikasi:*\n\n" +
			"Proses evaluasi dokumen kelayakan umumnya membutuhkan waktu *1 hingga 3 hari kerja* tergantung antrean permohonan yang masuk.\n\n" +
			"Anda akan langsung menerima notifikasi WhatsApp saat status diperbarui oleh Penilai."

	case strings.Contains(lowerMsg, "buat") || strings.Contains(lowerMsg, "tambah") || strings.Contains(lowerMsg, "pengajuan"):
		return "➕ *Cara Membuat Pengajuan Baru:*\n\n" +
			"1. Klik tombol *+ Buat Pengajuan Baru* di Dashboard atau menu *Daftar Pengajuan*.\n" +
			"2. Isi Judul, Perusahaan, Unit Kerja, dan pilih Jenis Dokumen.\n" +
			"3. Upload berkas dokumen kelayakan pada Step 2.\n" +
			"4. Klik *Kirim Permohonan*."

	case strings.Contains(lowerMsg, "halo") || strings.Contains(lowerMsg, "hai") || strings.Contains(lowerMsg, "selamat"):
		return "Halo! 👋 Saya adalah **Asisten AI Kelayakan**. Ada yang bisa saya bantu terkait pengajuan dokumen, cara menghubungi admin, status permohonan, atau aturan berkas hari ini?"

	default:
		return fmt.Sprintf("🤖 Mengenai pertanyaan Anda: *\"%s\"*\n\n"+
			"Untuk bantuan lebih lanjut pada sistem pengajuan dokumen kelayakan ini:\n"+
			"• Jika ada masalah akun/sistem: Silakan cek menu *Profil Saya* atau hubungi Admin.\n"+
			"• Jika terkait permohonan: Cek detail permohonan di menu *Daftar Pengajuan* untuk membaca catatan penilai.\n"+
			"• Jika ingin mengunduh laporan: Klik tombol *Export Rekapan (Excel)* di Dashboard.", msg)
	}
}
