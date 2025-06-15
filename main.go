// main.go
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// DiscordWebhookPayload represents the data we expect to RECEIVE from the POST request,
// and also the data we will SEND to Discord.
type DiscordWebhookPayload struct {
	Content  string `json:"content,omitempty"`
	Username string `json:"username,omitempty"`
}

func main() {
	// We still need the Discord Webhook URL to forward the message.
	if os.Getenv("DISCORD_WEBHOOK_URL") == "" {
		log.Fatal("FATAL: DISCORD_WEBHOOK_URL environment variable not set.")
	}

	// Create an endpoint named /api/receive that will listen for POST requests.
	http.HandleFunc("/api/receive", sendWebhookHandler)

	// Start the server.
	port := "8080"
	log.Printf("INFO: Server started on port %s", port)
	log.Printf("INFO: Waiting to receive a message via POST at http://localhost:%s/api/receive", port)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("FATAL: Could not start server: %s", err)
	}
}

// sendWebhookHandler is the core function that handles the incoming POST request.
func sendWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// We only want to handle POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}


	// 1. Create a variable to hold the incoming data.
	//    The structure of `DiscordWebhookPayload` must match the JSON you send.
	var payload DiscordWebhookPayload

	// 2. Decode the JSON from the request's body.
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error decoding JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Now, the message is available in `payload.Content`.
	//    We can log it to the console to prove we received it.
	log.Printf("SUCCESS: Received message from POST request: '%s'", payload.Content)

    // Optional: Check if the message is empty.
    if payload.Content == "" {
        http.Error(w, "The 'content' field cannot be empty.", http.StatusBadRequest)
        return
    }

	// --- Now we forward this received message to Discord ---

	// 4. Get the Discord URL from the environment variable.
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	// 5. Marshal the payload we received back into JSON format to send it to Discord.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error preparing message for Discord", http.StatusInternalServerError)
		return
	}

	// 6. Send the request to Discord.
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Error sending webhook to Discord", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 7. Check Discord's response and inform the original caller.
	if resp.StatusCode >= 300 {
		log.Printf("WARN: Discord returned a non-success status: %s", resp.Status)
		// We can still call our own request a success since we did receive and process it.
	}

	// 8. Send a success response back to the client that sent the POST request.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":          "success",
		"receivedMessage": payload.Content, // Echo back the message we got.
	})
}