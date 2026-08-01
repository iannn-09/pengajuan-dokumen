package main

// Request types
type sendRequest struct {
	Secret  string `json:"secret"`
	Target  string `json:"target"`
	Message string `json:"message"`
}

type bulkMessageRequest struct {
	Secret         string   `json:"secret"`
	Targets        []string `json:"targets"`
	Message        string   `json:"message"`
	SessionIDs     []string `json:"session_ids,omitempty"`
	Strategy       string   `json:"strategy,omitempty"` // "single", "even", "random"
	SendSuratTugas bool     `json:"send_surat_tugas,omitempty"`
}

type bulkDifferentMessageRequest struct {
	Secret   string `json:"secret"`
	Messages []struct {
		Targets   string `json:"targets"`
		Message   string `json:"message"`
		SessionID string `json:"session_id"` // Required: BE assigns which session for this message
	} `json:"messages"`
	SendSuratTugas bool `json:"send_surat_tugas,omitempty"`
}

type createSessionRequest struct {
	Secret string `json:"secret"`
	Name   string `json:"name"`
}

type updateSessionRequest struct {
	Secret string `json:"secret"`
	Name   string `json:"name"`
}

type sessionSecretRequest struct {
	Secret string `json:"secret"`
}
