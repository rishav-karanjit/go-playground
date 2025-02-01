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

type DbEsdkMiddleware struct {
	originalRequests map[string]interface{}
}

func NewDbEsdkMiddleware() *DbEsdkMiddleware {
	return &DbEsdkMiddleware{
		originalRequests: make(map[string]interface{}),
	}
}

func (m *DbEsdkMiddleware) CreateMiddleware() func(options *dynamodb.Options) {
	return func(options *dynamodb.Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			// Add request interceptor at the beginning of Initialize step
			requestIntercetor := m.createRequestInterceptor()
			if err := stack.Initialize.Add(requestIntercetor, middleware.Before); err != nil {
				return err
			}
			// Add response interceptor at the end of Finalize step
			return stack.Finalize.Add(m.createResponseInterceptor(), middleware.After)
		})
	}
}

func (m *DbEsdkMiddleware) createRequestInterceptor() middleware.InitializeMiddleware {
	return middleware.InitializeMiddlewareFunc("RequestInterceptor", func(
		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
	) (
		out middleware.InitializeOutput, metadata middleware.Metadata, err error,
	) {
		m.handleRequestInterception(in.Parameters)
		return next.HandleInitialize(ctx, in)
	})
}

// handleRequestInterception handles the interception logic before the DynamoDB operation
func (m *DbEsdkMiddleware) handleRequestInterception(params interface{}) {
	if v, ok := params.(*dynamodb.PutItemInput); ok {
		m.originalRequests["PutItemInput"] = *v
		fmt.Println("Original PutItemInput:")
		fmt.Println(m.originalRequests["PutItemInput"].(dynamodb.PutItemInput).Item["intercepted attribute"].(*types.AttributeValueMemberS))

		if idAttr, ok := v.Item["ID"]; ok {
			if idNum, ok := idAttr.(*types.AttributeValueMemberN); ok {
				if idNum.Value == "1" {
					v.Item["intercepted attribute"] = &types.AttributeValueMemberS{Value: "Intercepted by me"}
				}
			}
		}
	}
}

// createResponseInterceptor creates and returns the middleware interceptor for responses
func (m *DbEsdkMiddleware) createResponseInterceptor() middleware.FinalizeMiddleware {
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
		m.handleResponseInterception(result.Result)
		return result, metadata, err
	})
}

// handleResponseInterception handles the interception logic after the DynamoDB operation
func (m *DbEsdkMiddleware) handleResponseInterception(response interface{}) {
	if _, ok := response.(*dynamodb.PutItemOutput); ok {
		fmt.Println("PutItemOutput Response intercepted:")
		if att, ok := m.originalRequests["PutItemInput"].(dynamodb.PutItemInput).Item["intercepted attribute"].(*types.AttributeValueMemberS); ok {
			fmt.Println("intercepted attribute:", att.Value)
		}
		// You can modify the response here if needed
	}
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
	// configureInterceptors(&cfg)

	// Create DynamoDB client
	ddbMiddleware := NewDbEsdkMiddleware().CreateMiddleware()
	client := dynamodb.NewFromConfig(cfg, ddbMiddleware)

	// Create the input for PutItem
	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String("TestEverythingHere"),
		Item: map[string]types.AttributeValue{
			"ID":                    &types.AttributeValueMemberN{Value: "1"},
			"Name":                  &types.AttributeValueMemberS{Value: "John Doe"},
			"Age":                   &types.AttributeValueMemberN{Value: "30"},
			"intercepted attribute": &types.AttributeValueMemberS{Value: "Yo"},
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
