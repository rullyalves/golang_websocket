package models

type MessageType string

const (
	MessageTypeUpdate   MessageType = "update"
	MessageTypeDelete   MessageType = "delete"
	MessageTypeSnapshot MessageType = "snapshot"
)
