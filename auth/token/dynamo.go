package token

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
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

func (d *DynamoTokenManager) AddToken(user, token string) error {
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
		return err
	}

	if _, err := d.dynamo.PutItem(input); err != nil {
		return err
	}

	return nil
}

func (d *DynamoTokenManager) Authenticate(token string) error {
	return nil
}
