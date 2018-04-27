package auth

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/d.ims.io/mock"
	"github.com/stretchr/testify/assert"
)

func TestDynamoGrantAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mock.NewMockDynamoDBAPI(ctrl)
	target := NewDynamoAccountManager("table", mockDynamoDB)

	item := map[string]*dynamodb.AttributeValue{
		"AccountID": {
			S: aws.String("account-id"),
		},
	}

	input := &dynamodb.PutItemInput{}
	input.SetTableName("table")
	input.SetItem(item)

	mockDynamoDB.EXPECT().
		PutItem(input).
		Return(&dynamodb.PutItemOutput{}, nil)

	if err := target.GrantAccess("account-id"); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoRevokeAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mock.NewMockDynamoDBAPI(ctrl)
	target := NewDynamoAccountManager("table", mockDynamoDB)

	item := map[string]*dynamodb.AttributeValue{
		"AccountID": {
			S: aws.String("account-id"),
		},
	}

	input := &dynamodb.DeleteItemInput{}
	input.SetTableName("table")
	input.SetKey(item)

	mockDynamoDB.EXPECT().
		DeleteItem(input).
		Return(&dynamodb.DeleteItemOutput{}, nil)

	if err := target.RevokeAccess("account-id"); err != nil {
		t.Fatal(err)
	}
}

func TestDynamoGetAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mock.NewMockDynamoDBAPI(ctrl)
	target := NewDynamoAccountManager("table", mockDynamoDB)

	input := &dynamodb.ScanInput{}
	input.SetTableName("table")

	scanOutput := &dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{"AccountID": &dynamodb.AttributeValue{S: aws.String("1")}},
			{"AccountID": &dynamodb.AttributeValue{S: aws.String("2")}},
			{"AccountID": &dynamodb.AttributeValue{S: aws.String("3")}},
		},
	}

	mockDynamoDB.EXPECT().
		Scan(input).
		Return(scanOutput, nil)

	result, err := target.Accounts()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"1", "2", "3"}
	assert.Equal(t, expected, result)
}
