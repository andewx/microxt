package net

import (
	"fmt"

	"github.com/andewx/microxt/proto"
	pr "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/*
Talker structure manages system details for handling UDP connection status and streaming, for
example if we are in a UDP_PROCEDURAL_MODE with a Talker system we can expect that requests are
congenial and wait on the readiness status of the reciever. That is a transmission has a sequence
of {NotifyConsumer, VerifyConsumerStatusReady, Send Message, ExpectResponse or EnterIdleStatus.
On the other hand if we are expecting stream data, we can coordinate the stream parameters and
have a producer and consumer stream or sync stream as needed. (Sync Stream occassionally verifies
that our stream is valid an in sync, if the producer has over produced we can reduce the stream
parameters
*/

type Messages struct {
	Ack          proto.Ack
	Adc          proto.AdcData
	Direct       proto.Directive
	Ddat         proto.DdatData
	Fft          proto.FftData
	FrameDone    proto.FrameDone
	Pdat         proto.PdatData
	RadarConfig  proto.RadarConfigData
	StreamConfig proto.StreamConfig
	Sync         proto.Sync
	Payload      proto.PayloadData
	Handshake    proto.Handshake
}

type Talker struct {
	MyID           string // My ID
	DeviceID       string // Device ID
	DeviceStatus   int    // Device Status
	LocalStatus    int    // Conversation Status
	Message        []byte // Message
	Command        string // Message Context
	ConversationID string // Conversation ID
	Mode           int    // Send/Recieve Mode
	Inbox          Messages
}

func NewTalker() *Talker {
	talker := new(Talker)
	talker.DeviceStatus = NOT_READY
	talker.LocalStatus = NOT_READY
	return talker
}

// Heavy weight protocol messaging to the device/client passes all relevant message information
func (t *Talker) Receive(message []byte) error {

	if pr.Unmarshal(message, &t.Inbox.Ack) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.Adc) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.Direct) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.Ddat) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.Fft) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.FrameDone) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.Pdat) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.RadarConfig) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.StreamConfig) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.Sync) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.Payload) == nil {
	} else if pr.Unmarshal(message, &t.Inbox.Handshake) == nil {
	} else {
		return fmt.Errorf("%serror%s message type not recognized over serial port interface", CS_RED, CS_WHITE)
	}
	return nil
}

func (t *Talker) Send(msg protoreflect.ProtoMessage) ([]byte, error) {
	return pr.Marshal(msg)
}
