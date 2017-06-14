package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Provider struct {
	DynamoDB dynamodbiface.DynamoDBAPI
	EC2      ec2iface.EC2API
	S3       s3iface.S3API
}

func NewProvider(config *aws.Config) *Provider {
	session := session.New(config)
	return &Provider{
		DynamoDB: dynamodb.New(session),
		EC2:      ec2.New(session),
		S3:       s3.New(session),
	}
}
