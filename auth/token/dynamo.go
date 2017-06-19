package token

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"math/rand"
)

type DynamoTokenManager struct {
	table  string
	dynamo dynamodbiface.DynamoDBAPI
}

func NewDynamoTokenManager(table string, d dynamodbiface.DynamoDBAPI) *DynamoTokenManager {
	return &DynamoTokenManager{
		table:  table,
		dynamo: d,
	}
}

func (d *DynamoTokenManager) GenerateToken(user string) (string, error) {
	raw := fmt.Sprintf("%s:%s", randomString(26), randomString(26))
	token := base64.StdEncoding.EncodeToString([]byte(raw))

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

	if _, err := d.dynamo.PutItem(input); err != nil {
		return "", err
	}

	return token, nil
}

func (d *DynamoTokenManager) Authenticate(user, pass string) (bool, error) {
	data := fmt.Sprintf("%s:%s", user, pass)
	token := base64.StdEncoding.EncodeToString([]byte(data))

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

	output, err := d.dynamo.GetItem(input)
	if err != nil {
		return false, err
	}

	return len(output.Item) > 0, nil
}

func randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	runes := make([]rune, length)
	for i := range runes {
		runes[i] = letters[rand.Intn(len(letters))]
	}

	return string(runes)
}
