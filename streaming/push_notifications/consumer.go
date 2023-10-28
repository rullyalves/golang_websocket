package messaging

import (
	"fmt"
	"go_ws/shared/queue"
)

func StartPushConsumers(processQueue queue.ProcessMessage, processor queue.MessageProcessor) {
	work := queue.StartWorker(processQueue, processor)

	matchQueueUrl := fmt.Sprintf(queue.MatchPushQueueUrl)
	likeQueueUrl := fmt.Sprintf(queue.LikePushQueueUrl)
	messageQueueUrl := fmt.Sprintf(queue.MessagePushQueueUrl)

	for i := 0; i < 10; i++ {
		fmt.Printf("notification worker %v started \n", i)

		go work(matchQueueUrl)
		go work(likeQueueUrl)
		go work(messageQueueUrl)
	}
}
