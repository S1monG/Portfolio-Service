package incrementcounter

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
	// Replace with your project ID
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

	functions.HTTP("increment-simple-counter", func(w http.ResponseWriter, r *http.Request) {
		// Increment the counter
		updatedCounter, err := updateFireStore(ctx, client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the response as JSON
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := &response{Counter: updatedCounter}
		json.NewEncoder(w).Encode(jsonResponse)
	})

	functions.HTTP("simple-counter", func(w http.ResponseWriter, r *http.Request) {
		counter, err := getFireStore(ctx, client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the response as JSON
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := &response{Counter: counter}
		json.NewEncoder(w).Encode(jsonResponse)
	})

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

/*
In the future I might want to think about updating the counter with a transaction insead of directly like im doing here

It depends on the requierments of the application which are not fully set yet. Transactions adds more complexity and
overhead but they also garantee correctness by atomicity and eliminating race conditions.
*/
func updateFireStore(ctx context.Context, client *firestore.Client) (int64, error) {
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

	// Increment the field data
	newCounter := counterField.(int64) + 1

	// Set the field data
	docRef.Set(ctx, map[string]interface{}{
		"simple-counter": newCounter,
	})

	return newCounter, nil
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
