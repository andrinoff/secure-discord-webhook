// api/webhook.go (Updated to allow only one origin)

package contact

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type DiscordWebhookPayload struct {
	Content  string `json:"content,omitempty"`
	Username string `json:"username,omitempty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// --- NEW: Get allowed origin from environment variables ---
	allowedOrigin := "https://tbilisi.hackclub.com"
	if allowedOrigin == "" {
		// This is a server configuration error, so we block the request.
		log.Println("FATAL: ALLOWED_ORIGIN environment variable not set.")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	// --- NEW: Check if the request's origin matches the allowed origin ---
	requestOrigin := r.Header.Get("Origin")
	if requestOrigin != allowedOrigin {
		// If the origin does not match, block the request.
		http.Error(w, "Forbidden: Invalid origin", http.StatusForbidden)
		return
	}

	// --- UPDATED: Set CORS headers to be specific, not a wildcard ---
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin) // Only allow YOUR site
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL_2")
	if webhookURL == "" {
		http.Error(w, "Server configuration error: DISCORD_WEBHOOK_URL not set", http.StatusInternalServerError)
		return
	}

	var payload DiscordWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Error decoding JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error preparing message for Discord", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Error creating Discord request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending webhook to Discord", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}
