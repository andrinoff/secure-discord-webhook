// api/webhook.go (Corrected for CORS)

package handler

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
	// --- NEW CODE BLOCK TO HANDLE CORS ---
	// This allows the browser/fetch client to make a "preflight" OPTIONS request
	// to see which methods and headers are allowed.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// If this is the preflight OPTIONS request, we can stop here.
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	// --- END OF NEW CODE BLOCK ---


	// We still only want to process POST requests for the actual logic.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		http.Error(w, "Server configuration error: DISCORD_WEBHOOK_URL not set", http.StatusInternalServerError)
		log.Println("FATAL: DISCORD_WEBHOOK_URL environment variable not set.")
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