package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type response struct {
	Counter int64 `json:"counter"`
}

func init() {
	// First argument is the function name as it will be called by the Functions Framework,
	// and the name of the entry-point flag when deploying.
	// Second argument is the function itself
	// The path is specified when deploying the function
	functions.HTTP("GetCounter", getCounter)
	functions.HTTP("IncrementCounter", incrementCounter)
}

func enableCors(w *http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	// Check if the origin is https://s1mong.github.io or its subpaths
	if strings.HasPrefix(origin, "https://s1mong.github.io") {
		(*w).Header().Set("Access-Control-Allow-Origin", origin)
	}
}

func getCounter(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
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

// TODO: ugly repeated code, fix later
func incrementCounter(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
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

	updatedCounter, err := updateFireStore(ctx, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := &response{Counter: updatedCounter}
	json.NewEncoder(w).Encode(jsonResponse)
}

func getFireStore(ctx context.Context, client *firestore.Client) (int64, error) {
	// Get the Document snapshot and extract the field data
	docRef := client.Collection("portfolio").Doc("counters")
	dsnap, err := docRef.Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get document: %w", err)
	}
	counterField, err := dsnap.DataAt("simple-counter")
	if err != nil {
		return 0, fmt.Errorf("failed to get field: %w", err)
	}

	return counterField.(int64), nil
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
		return 0, fmt.Errorf("failed to get document: %w", err)
	}
	counterField, err := dsnap.DataAt("simple-counter")
	if err != nil {
		return 0, fmt.Errorf("failed to get field: %w", err)
	}

	// Increment the field data
	newCounter := counterField.(int64) + 1

	// Set the field data
	docRef.Set(ctx, map[string]interface{}{
		"simple-counter": newCounter,
	})

	return newCounter, nil
}
