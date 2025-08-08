package blockchain

import (
	"fmt"
	"log"
	"os"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var contract *gateway.Contract

// InitializeFabricGateway sets up the connection to the Hyperledger Fabric network
func InitializeFabricGateway() error {
	log.Println("Initializing Hyperledger Fabric Gateway...")

	// --- PATHS: YOU MUST CHECK AND UPDATE THESE ---
	// Path to your connection profile file
	credPath := "/mnt/c/Users/terra/Documents/my_complete_blockchain_project/chaincode/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.yaml"
	
	// Path to the Admin user's certificate (CORRECTED FILENAME)
	certPath := "/mnt/c/Users/terra/Documents/my_complete_blockchain_project/chaincode/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/cert.pem"
	
	// Path to the Admin user's private key
	// Note: the name of the key file can vary. Check your folder for the correct name (e.g., in `keystore` folder, it's usually a long string ending in `_sk`).
	keyPath := "/mnt/c/Users/terra/Documents/my_complete_blockchain_project/chaincode/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/6514006af8892095b0eb7ed65941b3fc4f023b84042f9cf46e50019027c4cd9c_sk" 
	// --- END PATHS ---

	walletPath := "./wallet"
	userName := "Admin"
	
	// Create the file system wallet
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	// Check if the Admin identity already exists in the wallet.
	if _, err := wallet.Get(userName); err != nil {
		log.Printf("Identity '%s' not found in wallet. Creating it from user credentials...", userName)
		
		// Read the certificate from the file system
		cert, err := os.ReadFile(certPath)
		if err != nil {
			return fmt.Errorf("failed to read certificate from path %s: %w", certPath, err)
		}
		
		// Read the private key from the file system
		key, err := os.ReadFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to read private key from path %s: %w", keyPath, err)
		}

		// Create a new gateway identity using the cert and key
		identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
		
		// Put the identity into the wallet for future use
		if err := wallet.Put(userName, identity); err != nil {
			return fmt.Errorf("failed to put identity into wallet: %w", err)
		}
		log.Printf("Identity '%s' successfully added to wallet.", userName)
	}

	// Set up the gateway connection
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(credPath)),
		gateway.WithIdentity(wallet, userName),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gateway: %w", err)
	}

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		return fmt.Errorf("failed to get network from gateway: %w", err)
	}

	contract = network.GetContract("farmcred")

	log.Println("Successfully connected to Fabric Gateway and contract.")
	return nil
}

// InvokeChaincode submits a transaction to the blockchain.
func InvokeChaincode(functionName string, args ...string) ([]byte, error) {
	if contract == nil {
		return nil, fmt.Errorf("Fabric contract not initialized")
	}
	log.Printf("Invoking chaincode function: %s with args: %v", functionName, args)
	
	result, err := contract.SubmitTransaction(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %w", err)
	}
	
	log.Printf("Transaction submitted successfully. Result: %s", string(result))
	return result, nil
}