package apis

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go_ws/shared/database/mongo"
	neo4jdb "go_ws/shared/database/neo4j"
	"go_ws/shared/queue"
	"log"
	"os"
)

func StartConnections(ctx context.Context) (*sqs.Client, *mongo.Client, *neo4j.DriverWithContext, error) {

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsAccessSecret := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_SQS_REGION")

	sqsClient, sqsErr := queue.GetSqsClient(ctx, awsAccessKey, awsAccessSecret, awsRegion)

	if sqsErr != nil {
		log.Printf("error while starting %v\n", sqsErr)
		return nil, nil, nil, sqsErr
	}

	mongoUri := os.Getenv("MONGO_URI")

	mongoDriver, mongoErr := mongodb.NewDriver(mongoUri)

	if mongoErr != nil {
		log.Printf("error while connecting to mongodb %v\n", mongoErr)
		return nil, nil, nil, mongoErr
	}

	neo4jUri := os.Getenv("NEO4J_URI")
	neo4jUsername := os.Getenv("NEO4J_USERNAME")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")

	neo4jDriver, neo4jErr := neo4jdb.NewDriver(neo4jUri, neo4jUsername, neo4jPassword)

	if neo4jErr != nil {
		log.Printf("error while connecting to neo4j %v", neo4jErr)
		return nil, nil, nil, neo4jErr
	}

	return sqsClient, mongoDriver, neo4jDriver, nil
}
