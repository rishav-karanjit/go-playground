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

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"), // Replace with your region
	)
	if err != nil {
		fmt.Printf("Unable to load SDK config, %v\n", err)
		return
	}

	var defaultBucket = middleware.InitializeMiddlewareFunc("DefaultBucket", func(
		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
	) (
		out middleware.InitializeOutput, metadata middleware.Metadata, err error,
	) {
		// Type switch to check if the input is s3.GetObjectInput, if so and the bucket is not set, populate it with
		// our default.
		switch v := in.Parameters.(type) {
		case *dynamodb.PutItemInput:
			if idAttr, ok := v.Item["ID"]; ok {
				// Check if it's of type *types.AttributeValueMemberN
				if idNum, ok := idAttr.(*types.AttributeValueMemberN); ok {
					if idNum.Value == "12345" {
						v.Item["intercepted attribute"] = &types.AttributeValueMemberS{Value: "Intercepted by me"}
						fmt.Println("Intercepted by me")
					}
				}
			}
		}
		// Middleware must call the next middleware to be executed in order to continue execution of the stack.
		// If an error occurs, you can return to prevent further execution.
		return next.HandleInitialize(ctx, in)
	})

	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
		// Attach the custom middleware to the beginning of the Initialize step
		return stack.Initialize.Add(defaultBucket, middleware.Before)
	})

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Create the input for PutItem
	input := &dynamodb.PutItemInput{
		TableName: aws.String("TestEverythingHere"), // Replace with your table name
		Item: map[string]types.AttributeValue{
			"ID":   &types.AttributeValueMemberN{Value: "12345"},
			"Name": &types.AttributeValueMemberS{Value: "John Doe"},
			"Age":  &types.AttributeValueMemberN{Value: "30"},
		},
	}

	// Make the PutItem request
	result, err := client.PutItem(context.TODO(), input)
	if err != nil {
		fmt.Printf("Failed to put item, %v\n", err)
		return
	}

	fmt.Println("Successfully added item to DynamoDB table")
	fmt.Println(result.Attributes)
}
