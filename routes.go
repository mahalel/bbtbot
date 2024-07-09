package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		return
	}
}

func interactionHandler(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")
	publicKey, err := hex.DecodeString(os.Getenv("APPLICATION_PUBLIC_KEY"))
	if err != nil {
		log.Fatalf("Failed to decode public key: %v", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	message := []byte(timestamp + string(body))
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	if !ed25519.Verify(publicKey, message, signatureBytes) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	var decodedBody map[string]interface{}
	err = json.Unmarshal(body, &decodedBody)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	switch decodedBody["type"].(float64) {
	case 1:
		err = json.NewEncoder(w).Encode(map[string]int{"type": 1})
		if err != nil {
			return
		}
	case 2:
		username := decodedBody["member"].(map[string]interface{})["user"].(map[string]interface{})["username"].(string)
		question := decodedBody["data"].(map[string]interface{})["options"].([]interface{})[0].(map[string]interface{})["value"].(string)
		responseContent := fmt.Sprintf("Question from %s: %s", username, question)

		response := map[string]interface{}{
			"type": 4,
			"data": map[string]string{"content": responseContent},
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}

		//go checkOpenAI(decodedBody, signature, timestamp, question)
		fmt.Printf("Checking Claude")
	default:
		http.Error(w, "Unknown interaction type", http.StatusBadRequest)
	}
}
