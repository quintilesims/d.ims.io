package token

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/quintilesims/d.ims.io/aws"
)

type DynamoAuth struct {
	table string
	aws   *aws.Provider
}

func NewDynamoAuth(table string, a *aws.Provider) *DynamoAuth {
	return &DynamoAuth{
		table: table,
		aws:   a,
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

	if _, err := d.aws.DynamoDB.PutItem(input); err != nil {
		return err
	}

	return nil
}

func (d *DynamoAuth) Authenticate(token string) error {
	return nil
}
