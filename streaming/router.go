package streaming

import (
	"go.mongodb.org/mongo-driver/mongo"
	_ "go_ws/shared/http_router"
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
	"go_ws/streaming/shared/websocket"
)

func Handle(wsRouter *websocket.Router, mongoDriver *mongo.Client) {

	//messageDao.DeleteAll(mongoDriver)(context.Background())

	chatBr := chat.GetChatStream(mongoDriver)
	messagesBr := message.GetMessageStream(mongoDriver)
	deliveryBr := delivery.GetDeliveryStream(mongoDriver)
	profileBr := profile.GetProfileStream(mongoDriver)
	likeBr := like.GetLikeStream(mongoDriver)
	dislikeBr := dislike.GetDislikeStream(mongoDriver)
	blockBr := block.GetBlockStream(mongoDriver)

	wsRouter.Handle("chats", chat.GetChatsAsStream(chatBr, chatDao.FindByUserId(mongoDriver)))
	wsRouter.Handle("messages", message.GetMessagesAsStream(messagesBr, messageDao.FindByUserId(mongoDriver)))
	wsRouter.Handle("deliveries", delivery.GetDeliveriesAsStream(deliveryBr, deliveryDao.FindByUserId(mongoDriver)))
	wsRouter.Handle("profile", profile.GetProfilesAsStream(profileBr, profileDao.FindByUserId(mongoDriver)))
	wsRouter.Handle("likes", like.GetLikesAsStream(likeBr, likeDao.FindByUserId(mongoDriver)))
	wsRouter.Handle("dislikes", dislike.GetDislikesAsStream(dislikeBr, dislikeDao.FindByUserId(mongoDriver)))
	wsRouter.Handle("blocks", block.GetBlocksAsStream(blockBr, blockDao.FindByUserId(mongoDriver)))
}
