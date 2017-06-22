package auth

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/d.ims.io/mock"
)

func TestDynamoCreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamoDB := mock.NewMockDynamoDBAPI(ctrl)
	target := NewDynamoTokenManager("table", mockDynamoDB)

	validatePutItemInput := func(input *dynamodb.PutItemInput) {
		if v, want := aws.StringValue(input.TableName), "table"; v != want {
			t.Errorf("Table was '%v', expected '%v'", v, want)
		}

		if v, want := aws.StringValue(input.Item["User"].S), "user"; v != want {
			t.Errorf("Column 'User' was '%v', expected '%v'", v, want)
		}

		if input.Item["Token"].S == nil {
			t.Error("Column 'Token' was nil")
		}
	}

	mockDynamoDB.EXPECT().
		PutItem(gomock.Any()).
		Do(validatePutItemInput).
		Return(&dynamodb.PutItemOutput{}, nil)

	if _, err := target.CreateToken("user"); err != nil {
		t.Fatal(err)
	}
}
