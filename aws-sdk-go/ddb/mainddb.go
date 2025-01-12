package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

type Item struct {
	ID   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}

func main() {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create a DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Table name
	tableName := "TestTable"

	// Put an item
	// item := Item{ID: "1", Name: "Example Item"}
	// err = putItem(client, tableName, item)
	// if err != nil {
	// 	log.Fatalf("Failed to put item: %v", err)
	// }

	// Get an item
	retrievedItem, err := getItem(client, tableName, "1")
	if err != nil {
		log.Fatalf("Failed to get item: %v", err)
	}
	fmt.Printf("Retrieved item: %+v\n", retrievedItem)
}

// func putItem(client *dynamodb.Client, tableName string, item Item) error {

// 	// Create an empty map with string keys and int values
// 	beacon := make(map[string]string)

// 	// Add key-value pairs
// 	beacon["branch-key-id"] = "test-mutate-encryption-add-value-da72ac87-5d28-4b84-8a56-96b4ea6e4326"
// 	beacon["type"] = "beacon:ACTIVE"

// 	// Create an empty map with string keys and int values
// 	activeitem := make(map[string]string)

// 	// Add key-value pairs
// 	activeitem["branch-key-id"] = "test-mutate-encryption-add-value-da72ac87-5d28-4b84-8a56-96b4ea6e4326"
// 	activeitem["type"] = "branch:ACTIVE"

// 	av, err := attributevalue.MarshalMap(beacon)
// 	av2, err := attributevalue.MarshalMap(activeitem)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal item: %w", err)
// 	}
// 	a, err := client.TransactWriteItems(context.TODO(), &dynamodb.TransactWriteItemsInput{
// 		TransactItems: []types.TransactWriteItem{
// 			{
// 				Put: &types.Put{
// 					TableName:           aws.String(tableName),
// 					Item:                av,
// 					ConditionExpression: aws.String(fmt.Sprintf("attribute_not_exists(%s)", "d5erf6tg")),
// 				},
// 			},
// 			{
// 				Put: &types.Put{
// 					TableName:           aws.String(tableName),
// 					Item:                av2,
// 					ConditionExpression: aws.String(fmt.Sprintf("attribute_exists(%s)", "d5erf6tg")),
// 				},
// 			},
// 		},
// 	})

// 	if err != nil {

// 		// Check if the error is because the condition was not met
// 		var ccfe *types.ConditionalCheckFailedException
// 		if ok := errors.As(err, &ccfe); ok {
// 			return fmt.Errorf("item already exists")
// 		}
// 		return fmt.Errorf("failed to put item: %w", err)
// 	}

// 	// _, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
// 	// 	TableName: aws.String(tableName),
// 	// 	Item:      av,
// 	// })
// 	fmt.Println(a)
// 	return err
// }

func getItem(client *dynamodb.Client, tableName string, id string) (Item, error) {
	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		// Check if the error is because the condition was not met
		var ccfe smithy.APIError
		var validation *types.validationException
		if ok := errors.As(err, &ccfe); ok {
			fmt.Println(ccfe.ErrorCode())
			fmt.Println(ccfe.ErrorMessage())
		}
	}
	if err != nil {
		return Item{}, err
	}

	var item Item
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return Item{}, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return item, nil
}
