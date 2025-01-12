package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	kmsv1 "github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/smithy-go"
)

func main() {
	fmt.Println("\n v2 client \n")
	mainv2()
	// fmt.Println("\n v1 client \n")
	// mainv1()
}
func mainv2() {
	// Load AWS credentials and configuration from environment variables or shared credentials file
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Failed to load AWS configuration:", err)
		return
	}

	// Create KMS service client
	svc := kms.NewFromConfig(cfg)

	// Plaintext data to encrypt
	plaintext := []byte("Hello, World!")
	key := "arn:aws:kms:us-west-2:992382771485:key/ad1d7ff3-79f3-40f4-b31d-7be4d9c8b3cc"
	// Encrypt the data
	input := &kms.EncryptInput{
		KeyId:     &key, // Replace with your KMS key ID
		Plaintext: plaintext,
	}

	result, err := svc.Encrypt(context.TODO(), input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			fmt.Printf("code: %s,\n message: %s,\n fault: %s", ae.ErrorCode(), ae.ErrorMessage(), ae.ErrorFault().String())
		}
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			fmt.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		}

		// var errr *types.Tag
		var err1 *types.CloudHsmClusterInUseException
		var err2 *types.NotFoundException
		switch {
		case errors.As(err, &err1):
			fmt.Println("Error is of type Error1:", err1)
		case errors.As(err, &err2):
			fmt.Println(reflect.TypeOf(err2))
			fmt.Println("Error is of type Error2:", err2)
		}
		// var oe *types.CloudHsmClusterInUseException
		// if errors.As(err, &oe) {
		// 	fmt.Println(reflect.TypeOf(oe))
		// 	// fmt.Println("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		// 	// fmt.Println(reflect.TypeOf(oe.Unwrap()))
		// }
		// fmt.Println(reflect.TypeOf(err))
		// fmt.Println(reflect.TypeOf(err.(*smithy.OperationError).Err))
		// fmt.Println("Failed to encrypt data:", err)
		return
	}

	fmt.Printf("Encrypted data: %x\n", result.CiphertextBlob)
}

func mainv1() {
	// Load AWS credentials and configuration from environment variables or shared credentials file
	sess := session.Must(session.NewSession(&awsv1.Config{
		LogLevel: awsv1.LogLevel(awsv1.LogDebugWithSigning),
	}))

	// Create KMS service client
	svc := kmsv1.New(sess)

	// Plaintext data to encrypt
	plaintext := []byte("Hello, World!")

	// Encrypt the data
	input := &kmsv1.EncryptInput{
		KeyId:     aws.String("arn:aws:kms:us-west-2:658956600833:key/b3537ef1-d8dc-4780-9f5a-55776cbb2f7g"), // Replace with your KMS key ID
		Plaintext: plaintext,
	}

	result, err := svc.Encrypt(input)
	if err != nil {
		// fmt.Println(reflect.TypeOf(err))
		// fmt.Println("Failed to encrypt data:", err)
		return
	}

	fmt.Printf("Encrypted data: %x\n", result.CiphertextBlob)
}
