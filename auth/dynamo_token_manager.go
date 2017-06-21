package auth

import (
	"fmt"
	"log"
	//"github.com/aws/aws-sdk-go/service/dynamodb"
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
	return "", fmt.Errorf("CreateToken not implemented")
}

func (d *DynamoTokenManager) DeleteToken(token string) error {
	return fmt.Errorf("DeleteToken not implemented")
}

func (d *DynamoTokenManager) Authenticate(user, pass string) (bool, error) {
	log.Println("[ERROR] - DynamoTokenManager.Authenticate not implemented")
	return true, nil
}
