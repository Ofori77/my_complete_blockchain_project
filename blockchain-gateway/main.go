package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"my_complete_blockchain_project/blockchain-gateway/blockchain" // *** IMPORTANT: UPDATE THIS IMPORT PATH ***
)

// This struct defines the data format we expect from the Django backend.
type RecordHashRequest struct {
	ID       string `json:"id"`
	Hash     string `json:"hash"`
	FarmerID string `json:"farmer_id"`
	BuyerID  string `json:"buyer_id"`
}

// This function will handle the API request to record a transaction hash.
func recordTransactionHash(c *gin.Context) {
	var req RecordHashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Now we call the real Fabric function instead of just logging!
	_, err := blockchain.InvokeChaincode("RecordVerifiedTransactionHash", req.Hash, req.FarmerID, req.BuyerID)
	if err != nil {
		log.Printf("Error invoking chaincode: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to record on blockchain: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hash for ID %s recorded successfully on blockchain.", req.ID)})
}

func main() {
	// --- NEW CODE: Initialize the Fabric Gateway BEFORE starting the server ---
	if err := blockchain.InitializeFabricGateway(); err != nil {
		log.Fatalf("Failed to initialize Fabric Gateway: %v", err)
	}
	// --- END NEW CODE ---

	// Initialize the web server framework
	router := gin.Default()

	// Define the API endpoint that will receive POST requests.
	router.POST("/api/record-verified-transaction-hash", recordTransactionHash)

	// Start the server and listen on port 8080.
	log.Println("Starting Go blockchain gateway on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}