package auth

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoTokenManager struct {
	table    string
	dynamodb dynamodbiface.DynamoDBAPI
}

func NewDynamoTokenManager(table string, dynamodb dynamodbiface.DynamoDBAPI) *DynamoTokenManager {
	return &DynamoTokenManager{
		table:    table,
		dynamodb: dynamodb,
	}
}

func (d *DynamoTokenManager) CreateToken(user string) (string, error) {
	token := convertToToken(randomString(20), randomString(20))

	item := map[string]*dynamodb.AttributeValue{
		"User": {
			S: &user,
		},
		"Token": {
			S: &token,
		},
	}

	input := &dynamodb.PutItemInput{}
	input.SetTableName(d.table)
	input.SetItem(item)

	if err := input.Validate(); err != nil {
		return "", err
	}

	if _, err := d.dynamodb.PutItem(input); err != nil {
		return "", err
	}

	return token, nil
}

func (d *DynamoTokenManager) DeleteToken(token string) error {
	key := map[string]*dynamodb.AttributeValue{
		"Token": {
			S: &token,
		},
	}

	input := &dynamodb.DeleteItemInput{}
	input.SetTableName(d.table)
	input.SetKey(key)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := d.dynamodb.DeleteItem(input); err != nil {
		return err
	}

	return nil
}

func (d *DynamoTokenManager) Authenticate(user, pass string) (bool, error) {
	log.Printf("[DEBUG] Attempting to authenticate user '%s' through DynamoDB", user)

	token := convertToToken(user, pass)
	key := map[string]*dynamodb.AttributeValue{
		"Token": {
			S: &token,
		},
	}

	input := &dynamodb.GetItemInput{}
	input.SetTableName(d.table)
	input.SetConsistentRead(true)
	input.SetKey(key)

	if err := input.Validate(); err != nil {
		return false, err
	}

	output, err := d.dynamodb.GetItem(input)
	if err != nil {
		return false, err
	}

	if len(output.Item) > 0 {
		log.Printf("[DEBUG] User '%s' sent valid DynamoDB credentials", user)
		return true, nil
	}

	log.Printf("[DEBUG] User '%s' sent invalid DynamoDB credentials", user)
	return false, nil
}

func convertToToken(user, pass string) string {
	s := fmt.Sprintf("%s:%s", user, pass)
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	runes := make([]rune, length)
	for i := range runes {
		runes[i] = letters[rand.Intn(len(letters))]
	}

	return string(runes)
}
