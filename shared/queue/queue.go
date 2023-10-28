package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"os"
)

func GetSqsClient(ctx context.Context, accessKey string, accessSecret string, regionStr string) (*sqs.Client, error) {

	provider := credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY"),
			SecretAccessKey: os.Getenv("AWS_SECRET_KEY"),
			CanExpire:       false,
		},
	}

	credentialsProvider := config.WithCredentialsProvider(provider)
	region := config.WithRegion(os.Getenv("AWS_SQS_REGION"))
	cfg, err := config.LoadDefaultConfig(ctx, region, credentialsProvider)

	if err != nil {
		return nil, err
	}

	sqsClient := sqs.NewFromConfig(cfg)

	return sqsClient, nil
}
