package ws

// SignalMessage defines the struture of messages exchanged via Websocket
// for signalling purposes such as WebRtc negotiation, room joining, etc.
type SignalMessage struct {
	Type    string      `json: "type"`              // offer, answer, cancdidate, join, leave, error
	Payload interface{} `json: "payload"`           // can be anything: SDP object, ICE candidate oject, room name, error
	Sender  string      `json: "sender"`            // Unique ID of the client sending the message
	Target  string      `json: "target, omitempty"` // Optional: Unique ID of the specific client this message is intended for
	Room    string      `json: "room, omitempty"`
}
