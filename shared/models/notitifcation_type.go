package models

type NotificationType string

const (
	NotificationTypeMatch   = "match"
	NotificationTypeMessage = "message"
	NotificationTypeLike    = "like"
)
