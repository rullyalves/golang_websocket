package utils

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func ParseMessagesTo[T any](messages []types.Message) ([]T, error) {
	var items []T

	for _, value := range messages {
		var item T
		err := json.Unmarshal([]byte(*value.Body), &item)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
