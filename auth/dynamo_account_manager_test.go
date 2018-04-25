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

	validatePutItemInput := func(input *dynamodb.PutItemInput) {
		if v, want := aws.StringValue(input.TableName), "table"; v != want {
			t.Errorf("Table was '%v', expected '%v'", v, want)
		}

		if v, want := aws.StringValue(input.Item["AccountID"].S), "account-id"; v != want {
			t.Errorf("Column 'AccountID' was '%v', expected '%v'", v, want)
		}

		if input.Item["AccountID"].S == nil {
			t.Error("Column 'Token' was nil")
		}
	}

	mockDynamoDB.EXPECT().
		PutItem(gomock.Any()).
		Do(validatePutItemInput).
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

	validateDeleteItemInput := func(input *dynamodb.DeleteItemInput) {
		if v, want := aws.StringValue(input.TableName), "table"; v != want {
			t.Errorf("Table was '%v', expected '%v'", v, want)
		}

		if v, want := aws.StringValue(input.Key["AccountID"].S), "account-id"; v != want {
			t.Errorf("Key 'AccountID' was '%v', expected '%v'", v, want)
		}
	}

	mockDynamoDB.EXPECT().
		DeleteItem(gomock.Any()).
		Do(validateDeleteItemInput).
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

	validateScanInput := func(input *dynamodb.ScanInput) {
		if v, want := aws.StringValue(input.TableName), "table"; v != want {
			t.Errorf("Table was %v, expected '%v'", v, want)
		}
	}

	scanOutput := &dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{"AccountID": &dynamodb.AttributeValue{S: aws.String("1")}},
			{"AccountID": &dynamodb.AttributeValue{S: aws.String("2")}},
			{"AccountID": &dynamodb.AttributeValue{S: aws.String("3")}},
		},
	}

	mockDynamoDB.EXPECT().
		Scan(gomock.Any()).
		Do(validateScanInput).
		Return(scanOutput, nil)

	result, err := target.Accounts()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"1", "2", "3"}
	assert.Equal(t, expected, result)
}
