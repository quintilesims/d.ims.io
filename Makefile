SHELL:=/bin/bash
VERSION?=$(shell git describe --tags --always)
CURRENT_DOCKER_IMAGE=quintilesims/d.ims.io:$(VERSION)
LATEST_DOCKER_IMAGE=quintilesims/d.ims.io:latest

deps:
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen

mocks:
	mockgen -package mock github.com/aws/aws-sdk-go/service/ecr/ecriface ECRAPI > mock/mock_ecr.go
	mockgen -package mock github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface DynamoDBAPI > mock/mock_dynamodb.go
	mockgen -package mock github.com/quintilesims/d.ims.io/auth TokenManager > mock/mock_token_manager.go
	mockgen -package mock github.com/quintilesims/d.ims.io/auth AccountManager > mock/mock_account_manager.go


build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags "-X main.Version=$(VERSION)" -o d.ims.io . 
	docker build -t $(CURRENT_DOCKER_IMAGE) .

release: build
	docker push $(CURRENT_DOCKER_IMAGE)
	docker tag  $(CURRENT_DOCKER_IMAGE) $(LATEST_DOCKER_IMAGE)
	docker push $(LATEST_DOCKER_IMAGE)

.PHONY: deps mocks build release
