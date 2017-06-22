deps:
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen

mocks:
	mockgen -package mock github.com/aws/aws-sdk-go/service/ecr/ecriface ECRAPI > mock/mock_ecr.go


.PHONY: mocks
