package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go/middleware"
)

// handleRequestInterception handles the interception logic before the DynamoDB operation
func handleRequestInterception(params interface{}) {
	if v, ok := params.(*dynamodb.PutItemInput); ok {
		if idAttr, ok := v.Item["ID"]; ok {
			if idNum, ok := idAttr.(*types.AttributeValueMemberN); ok {
				if idNum.Value == "1" {
					v.Item["intercepted attribute"] = &types.AttributeValueMemberS{Value: "Intercepted by me"}
					fmt.Println("Request intercepted")
				}
			}
		}
	}
}

// handleResponseInterception handles the interception logic after the DynamoDB operation
func handleResponseInterception(response interface{}) {
	if _, ok := response.(*dynamodb.PutItemOutput); ok {
		fmt.Println("PutItemOutput Response intercepted:")
		// You can modify the response here if needed
	}
	if getItemOutput, ok := response.(*dynamodb.GetItemOutput); ok {
		fmt.Println("GetItemOutput Response intercepted:")
		if age, ok := getItemOutput.Item["Age"].(*types.AttributeValueMemberN); ok {
			fmt.Println("Age:", age.Value)
		}
		if id, ok := getItemOutput.Item["ID"].(*types.AttributeValueMemberN); ok {
			fmt.Println("ID:", id.Value)
		}
		if name, ok := getItemOutput.Item["Name"].(*types.AttributeValueMemberS); ok {
			fmt.Println("Name:", name.Value)
		}
		if intercepted, ok := getItemOutput.Item["intercepted attribute"].(*types.AttributeValueMemberS); ok {
			fmt.Println("intercepted attribute:", intercepted.Value)
		}
		getItemOutput.Item["intercepted attribute"] = &types.AttributeValueMemberS{Value: "I read your data "}
		// You can modify the response here if needed
	}
}

// createRequestInterceptor creates and returns the middleware interceptor for requests
func createRequestInterceptor() middleware.InitializeMiddleware {
	return middleware.InitializeMiddlewareFunc("RequestInterceptor", func(
		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
	) (
		out middleware.InitializeOutput, metadata middleware.Metadata, err error,
	) {
		handleRequestInterception(in.Parameters)
		return next.HandleInitialize(ctx, in)
	})
}

// createResponseInterceptor creates and returns the middleware interceptor for responses
func createResponseInterceptor() middleware.FinalizeMiddleware {
	return middleware.FinalizeMiddlewareFunc("ResponseInterceptor", func(
		ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
	) (
		out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
	) {
		// First let the request complete
		result, metadata, err := next.HandleFinalize(ctx, in)
		if err != nil {
			return result, metadata, err
		}

		// Then intercept the response
		handleResponseInterception(result.Result)
		return result, metadata, err
	})
}

// configureInterceptors adds both request and response interceptors to the AWS configuration
func configureInterceptors(cfg *aws.Config) {
	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
		// Add request interceptor at the beginning of Initialize step
		if err := stack.Initialize.Add(createRequestInterceptor(), middleware.Before); err != nil {
			return err
		}
		// Add response interceptor at the end of Finalize step
		return stack.Finalize.Add(createResponseInterceptor(), middleware.After)
	})
}

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		fmt.Printf("Unable to load SDK config, %v\n", err)
		return
	}

	// Configure both interceptors
	configureInterceptors(&cfg)

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Create the input for PutItem
	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String("TestEverythingHere"),
		Item: map[string]types.AttributeValue{
			"ID":   &types.AttributeValueMemberN{Value: "1"},
			"Name": &types.AttributeValueMemberS{Value: "John Doe"},
			"Age":  &types.AttributeValueMemberN{Value: "30"},
		},
	}

	// Make the PutItem request
	_, err = client.PutItem(context.TODO(), putItemInput)
	if err != nil {
		fmt.Printf("Failed to put item, %v\n", err)
		return
	}

	fmt.Println("\tSuccessfully added item to DynamoDB table")

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("TestEverythingHere"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberN{Value: "1"},
		},
	}

	// Make the GetItem request
	resultGetItem, err := client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		fmt.Printf("Failed to put item, %v\n", err)
		return
	}

	fmt.Println("\tSuccessfully got item to DynamoDB table")
	// Print specific attributes
	if age, ok := resultGetItem.Item["Age"].(*types.AttributeValueMemberN); ok {
		fmt.Println("Age:", age.Value)
	}
	if id, ok := resultGetItem.Item["ID"].(*types.AttributeValueMemberN); ok {
		fmt.Println("ID:", id.Value)
	}
	if name, ok := resultGetItem.Item["Name"].(*types.AttributeValueMemberS); ok {
		fmt.Println("Name:", name.Value)
	}
	if intercepted, ok := resultGetItem.Item["intercepted attribute"].(*types.AttributeValueMemberS); ok {
		fmt.Println("intercepted attribute:", intercepted.Value)
	}
}
