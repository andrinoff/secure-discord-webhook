// api/webhook.go
package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// DiscordWebhookPayload is the structure for the JSON data.
type DiscordWebhookPayload struct {
	Content  string `json:"content,omitempty"`
	Username string `json:"username,omitempty"`
}

// Handler is the Vercel serverless function entrypoint.
// Vercel invokes this function for each incoming request.
func Handler(w http.ResponseWriter, r *http.Request) {
	// We only want to handle POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	// Get the Discord Webhook URL from Vercel Environment Variables.
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		http.Error(w, "Server configuration error: DISCORD_WEBHOOK_URL not set", http.StatusInternalServerError)
		log.Println("FATAL: DISCORD_WEBHOOK_URL environment variable not set.")
		return
	}

	// Decode the JSON from the incoming request body.
	var payload DiscordWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Error decoding JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Log the received message for debugging (visible in Vercel logs).
	log.Printf("Received message: '%s'", payload.Content)

	// Marshal the payload to be sent to Discord.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error preparing message for Discord", http.StatusInternalServerError)
		return
	}

	// Create a new request to forward to Discord.
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Error creating Discord request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request using a client with a timeout.
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending webhook to Discord", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check Discord's response status.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("WARN: Discord returned a non-success status: %s", resp.Status)
		// Consider how to report this error back to the client if needed.
	}

	// Send a success response back to the original client.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":          "success",
		"receivedMessage": payload.Content,
	})
}