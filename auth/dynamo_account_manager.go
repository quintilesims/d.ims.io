package auth

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoAccountManager struct {
	table    string
	dynamodb dynamodbiface.DynamoDBAPI
}

func NewDynamoAccountManager(table string, dynamodb dynamodbiface.DynamoDBAPI) *DynamoAccountManager {
	return &DynamoAccountManager{
		table:    table,
		dynamodb: dynamodb,
	}
}

func (d *DynamoAccountManager) GrantAccess(accountID string) error {
	item := map[string]*dynamodb.AttributeValue{
		"AccountID": {
			S: &accountID,
		},
	}

	input := &dynamodb.PutItemInput{}
	input.SetTableName(d.table)
	input.SetItem(item)

	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := d.dynamodb.PutItem(input); err != nil {
		return err
	}

	return nil
}

func (d *DynamoAccountManager) RevokeAccess(accountID string) error {
	key := map[string]*dynamodb.AttributeValue{
		"AccountID": {
			S: &accountID,
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

func (d *DynamoAccountManager) Accounts() ([]string, error) {
	input := &dynamodb.ScanInput{}
	input.SetTableName(d.table)

	if err := input.Validate(); err != nil {
		return nil, err
	}

	output, err := d.dynamodb.Scan(input)
	if err != nil {
		return nil, err
	}

	response := make([]string, len(output.Items))
	i := 0
	for _, v := range output.Items {
		response[i] = aws.StringValue(v["AccountID"].S)
		i++
	}

	return response, nil
}
