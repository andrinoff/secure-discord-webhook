package handler

import (
	"bytes"
	"encoding/json"

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
	allowedOrigins := map[string]bool{
		"https://tbilisi.hackclub.com":    true,
		"https://tbilisihc.andrinoff.com": true,
	}
	// 2. Get the origin of the request.
	origin := r.Header.Get("Origin")

	// 3. Check if the origin is in your whitelist and set the header.
	if allowedOrigins[origin] {
		w.Header().Set("Access-control-Allow-Origin", origin)
	}

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
