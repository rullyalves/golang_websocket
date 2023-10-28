package queue

import "os"

// delivery queues
var LikeStreamQueueUrl = os.Getenv("LIKE_STREAM_QUEUE_URL")
var DislikeStreamQueueUrl = os.Getenv("DISLIKE_STREAM_QUEUE_URL")
var ChatStreamQueueUrl = os.Getenv("CHAT_STREAM_QUEUE_URL")
var BlockStreamQueueUrl = os.Getenv("BLOCK_STREAM_QUEUE_URL")
var ProfileStreamQueueUrl = os.Getenv("PROFILE_STREAM_QUEUE_URL")
var DeliveryStreamQueueUrl = os.Getenv("DELIVERY_STREAM_QUEUE_URL")
var MessageStreamQueueUrl = os.Getenv("MESSAGE_STREAM_QUEUE_URL")

// push queues
var MatchPushQueueUrl = os.Getenv("MATCH_PUSH_QUEUE_URL")
var LikePushQueueUrl = os.Getenv("LIKE_PUSH_QUEUE_URL")
var MessagePushQueueUrl = os.Getenv("MESSAGE_PUSH_QUEUE_URL")
