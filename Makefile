deps:
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen

mocks:
	mockgen -package mock github.com/aws/aws-sdk-go/service/ecr/ecriface ECRAPI > mock/mock_ecr.go
	mockgen -package mock github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface DynamoDBAPI > mock/mock_dynamodb.go
	mockgen -package mock github.com/quintilesims/d.ims.io/auth TokenManager > mock/mock_token_manager.go

.PHONY: mocks
