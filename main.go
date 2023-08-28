package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var client *secretmanager.Client

func main() {
	// Create a new Secret Manager client
	ctx := context.Background()
	var err error
	client, err = secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	router := mux.NewRouter()
	router.HandleFunc("/secrets/{secretName}", GetSecret).Methods("GET")
	router.HandleFunc("/secrets/{secretName}", UpdateSecret).Methods("PUT")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is now listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func GetSecret(w http.ResponseWriter, r *http.Request) {
	secretName := mux.Vars(r)["secretName"]
	secretPath := fmt.Sprintf("projects/your-project-id/secrets/%s/versions/latest", secretName)

	// Access the secret version
	result, err := client.AccessSecretVersion(context.Background(), &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretPath,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to access secret", http.StatusInternalServerError)
		return
	}

	response := struct {
		SecretValue string `json:"secret_value"`
	}{
		SecretValue: string(result.Payload.Data),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateSecret(w http.ResponseWriter, r *http.Request) {
	secretName := mux.Vars(r)["secretName"]
	secretPath := fmt.Sprintf("projects/your-project-id/secrets/%s", secretName)

	var requestData struct {
		SecretValue string `json:"secret_value"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create a new secret version
	version, err := client.AddSecretVersion(context.Background(), &secretmanagerpb.AddSecretVersionRequest{
		Parent: secretPath,
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(requestData.SecretValue),
		},
	})
	if err != nil {
		http.Error(w, "Failed to update secret", http.StatusInternalServerError)
		return
	}

	response := struct {
		Version string `json:"version"`
	}{
		Version: version.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
