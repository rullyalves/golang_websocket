package streaming

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go_ws/shared/queue"
	"go_ws/streaming/block"
	blockDao "go_ws/streaming/block/dao"
	"go_ws/streaming/chat"
	chatDao "go_ws/streaming/chat/dao"
	"go_ws/streaming/delivery"
	deliveryDao "go_ws/streaming/delivery/dao"
	"go_ws/streaming/dislike"
	dislikeDao "go_ws/streaming/dislike/dao"
	"go_ws/streaming/like"
	likeDao "go_ws/streaming/like/dao"
	"go_ws/streaming/message"
	messageDao "go_ws/streaming/message/dao"
	"go_ws/streaming/profile"
	profileDao "go_ws/streaming/profile/dao"
)

func StartStreamingConsumers(processQueue queue.ProcessMessage, client *mongo.Client) func() {

	chatWork := chat.StartWorker(processQueue, chatDao.SaveAll(client))
	messageWork := message.StartWorker(processQueue, messageDao.SaveAll(client))
	deliveryWork := delivery.StartWorker(processQueue, deliveryDao.SaveAll(client))
	blockWork := block.StartWorker(processQueue, blockDao.SaveAll(client))
	dislikeWork := dislike.StartWorker(processQueue, dislikeDao.SaveAll(client))
	likeWork := like.StartWorker(processQueue, likeDao.SaveAll(client))
	profileWork := profile.StartWorker(processQueue, profileDao.SaveAll(client))

	chatQueueUrl := fmt.Sprintf(queue.ChatStreamQueueUrl)
	messageQueueUrl := fmt.Sprintf(queue.MessageStreamQueueUrl)
	deliveryQueueUrl := fmt.Sprintf(queue.DeliveryStreamQueueUrl)
	blockQueueUrl := fmt.Sprintf(queue.BlockStreamQueueUrl)
	dislikeQueueUrl := fmt.Sprintf(queue.DislikeStreamQueueUrl)
	likeQueueUrl := fmt.Sprintf(queue.LikeStreamQueueUrl)
	profileQueueUrl := fmt.Sprintf(queue.ProfileStreamQueueUrl)

	return func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("straming worker %v started \n", i)
			go chatWork(chatQueueUrl)
			go messageWork(messageQueueUrl)
			go deliveryWork(deliveryQueueUrl)
			go blockWork(blockQueueUrl)
			go dislikeWork(dislikeQueueUrl)
			go likeWork(likeQueueUrl)
			go profileWork(profileQueueUrl)
		}
	}
}
