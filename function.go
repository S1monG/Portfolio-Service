package function

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
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

	r := gin.Default()

	r.POST("/increment-simple-counter", func(c *gin.Context) {
		// Increment the counter
		updatedCounter, err := updateFireStore(ctx, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Send the response as JSON
		c.JSON(http.StatusOK, gin.H{"counter": updatedCounter})
	})

	r.GET("/simple-counter", func(c *gin.Context) {
		counter, err := getFireStore(ctx, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Send the response as JSON
		c.JSON(http.StatusOK, gin.H{"counter": counter})
	})

	r.Use(cors.Default())
	fmt.Println("Server listening on port 8080...")
	r.Run(":8080")
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
