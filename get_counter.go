package getcounter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type response struct {
	Counter int64 `json:"counter"`
}

func init() {
	functions.HTTP("getCounter", GetCounter)
	http.ListenAndServe(":8080", nil)
}

func GetCounter(w http.ResponseWriter, r *http.Request) {
	projectID := "sigma-tractor-429314-n0"

	// Initialize Firestore client
	ctx := context.Background()
	// No need to specify credentials if running on GCP
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create Firestore client: %v\n", err)
		return
	}
	defer client.Close()

	counter, err := getFireStore(ctx, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := &response{Counter: counter}
	json.NewEncoder(w).Encode(jsonResponse)
}

func getFireStore(ctx context.Context, client *firestore.Client) (int64, error) {
	// Get the Document snapshot and extract the field data
	docRef := client.Collection("portfolio").Doc("counters")
	dsnap, err := docRef.Get(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to get document: %w", err)
	}
	counterField, err := dsnap.DataAt("simple-counter")
	if err != nil {
		return -1, fmt.Errorf("failed to get field: %w", err)
	}

	return counterField.(int64), nil
}
