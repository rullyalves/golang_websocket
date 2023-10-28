package relationship_service

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go_ws/services/relationship_service/blocked_user"
	blockDao "go_ws/services/relationship_service/blocked_user/dao/neo4j"
	"go_ws/services/relationship_service/complaint"
	"go_ws/services/relationship_service/complaint_category"
	"go_ws/services/relationship_service/dislike"
	dislikeUsecases "go_ws/services/relationship_service/dislike/usecases"
	"go_ws/services/relationship_service/like"
	likeDao "go_ws/services/relationship_service/like/dao/neo4j"
	likeUsecases "go_ws/services/relationship_service/like/usecases"
	"go_ws/shared/http_router"
)

func Handle(router http_router.Router, neo4jDriver *neo4j.DriverWithContext, sqsClient *sqs.Client) {

	// likes
	findReceivedLikes := likeDao.FindByProfileId(neo4jDriver)
	router.Get("users/:userId/received-likes", like.FindReceivedLikesHandler(findReceivedLikes))

	sendLike := likeUsecases.NewSendLike(neo4jDriver, sqsClient)
	router.Post("users/:userId/send-like", like.SendLikeHandler(sendLike))

	// dislikes
	sendDislike := dislikeUsecases.NewSendDislike(neo4jDriver, sqsClient)
	router.Post("users/:userId/send-dislike", dislike.SendDislikeHandler(sendDislike))

	// blocks
	saveBlock := blockDao.Save(neo4jDriver)
	router.Post("users/:userId/send-block", blocked_user.SendBlockHandler(saveBlock))

	// complaint
	router.Post("/complaints", complaint.SendComplaintHandler())

	// complaint category
	router.Get("/complaint-categories", complaint_category.FindComplaintCategoriesHandler())
	router.Post("/complaint-categories", complaint_category.FindComplaintCategoriesHandler())
	router.Delete("/complaint-categories", complaint_category.FindComplaintCategoriesHandler())
}
