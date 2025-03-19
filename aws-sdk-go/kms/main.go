package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create KMS client
	kmsClient := kms.NewFromConfig(cfg)

	// Get KMS key ARN from environment variable
	keyID := os.Getenv("key1")
	if keyID == "" {
		log.Fatal("AWS_KMS_KEY_ID environment variable is required")
	}

	// Original data to encrypt
	plaintext := []byte("Hello, World!")

	// Encrypt the data
	ciphertext, err := encryptData(kmsClient, keyID, plaintext)
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}
	fmt.Printf("Encrypted data: %x\n", ciphertext)

	// Decrypt the data
	decryptedText, err := decryptData(kmsClient, ciphertext)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}
	fmt.Printf("Decrypted data: %s\n", string(decryptedText))
}

func encryptData(client *kms.Client, keyID string, plaintext []byte) ([]byte, error) {
	input := &kms.EncryptInput{
		KeyId:     &keyID,
		Plaintext: plaintext,
	}

	result, err := client.Encrypt(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("encryption error: %w", err)
	}

	return result.CiphertextBlob, nil
}

func decryptData(client *kms.Client, ciphertext []byte) ([]byte, error) {
	input := &kms.DecryptInput{
		CiphertextBlob: ciphertext,
	}

	result, err := client.Decrypt(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("decryption error: %w", err)
	}

	return result.Plaintext, nil
}
