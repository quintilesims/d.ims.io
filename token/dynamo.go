package token

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoAuth struct {
	table  string
	dynamo dynamodbiface.DynamoDBAPI
}

func NewDynamoAuth(table string, d dynamodbiface.DynamoDBAPI) *DynamoAuth {
	return &DynamoAuth{
		table:  table,
		dynamo: d,
	}
}

func (d *DynamoAuth) AddToken(user, token string) error {
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

	if _, err := d.dynamo.PutItem(input); err != nil {
		return err
	}

	return nil
}

func (d *DynamoAuth) Authenticate(token string) error {
	return nil
}
