package models

type RtcIceCandidate struct {
	ID            string `json:"id" bson:"Id"`
	RoomID        string `json:"roomId" bson:"roomId"`
	Candidate     string `json:"candidate" bson:"candidate"`
	SdpMid        string `json:"sdpMid" bson:"sdpMid"`
	SdpMLineIndex *int   `json:"sdpMLineIndex" bson:"sdpMLineIndex"`
}
