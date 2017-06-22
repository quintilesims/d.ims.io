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
	return fmt.Errorf("DeleteToken not implemented")
}

func (d *DynamoTokenManager) Authenticate(user, pass string) (bool, error) {
	log.Println("[ERROR] - DynamoTokenManager.Authenticate not implemented")
	return true, nil
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
