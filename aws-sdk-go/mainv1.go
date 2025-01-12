package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func mainv1() {
	// Load AWS credentials and configuration from environment variables or shared credentials file
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Failed to create AWS session:", err)
		return
	}

	// Create KMS service client
	svc := kms.New(sess)

	// Plaintext data to encrypt
	plaintext := []byte("Hello, World!")

	// Encrypt the data
	input := &kms.EncryptInput{
		KeyId:     aws.String("arn:aws:kms:us-west-2:658956600833:key/b3537ef1-d8dc-4780-9f5a-55776cbb2f7g"), // Replace with your KMS key ID
		Plaintext: plaintext,
	}

	result, err := svc.Encrypt(input)
	if err != nil {
		fmt.Println("Failed to encrypt data:", err)
		return
	}

	fmt.Printf("Encrypted data: %x\n", result.CiphertextBlob)
}
