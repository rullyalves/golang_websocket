package models

type RtcSessionDescription struct {
	RoomID string `json:"roomId" bson:"roomId"`
	Sdp    string `json:"sdp" bson:"sdp"`
	Type   string `json:"type" bson:"type"`
}
