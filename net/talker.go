package net

import (
	"fmt"
	"sync"

	"github.com/andewx/microxt/proto"
	pr "google.golang.org/protobuf/proto"
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

const (
	GPR_RECIEVE = 1
	GPR_ERROR   = 2
)

type Protobuf struct {
	Ack                proto.Ack
	Cmd                proto.Command
	CmdBufferQueue     proto.CommandBufferQueue
	SyncCmdBufferQueue proto.SyncCommandBufferQueue
	DataBuffer         proto.DataBuffer
	MultiDataBuffer    proto.MultiDataBuffer
	Start              proto.Start
	Stop               proto.Stop
	GetProfile         proto.GetProfile
	DeviceProfile      proto.DeviceProfile
	Error              proto.Error
	mutex              sync.Mutex
}

func NewProtobuf() *Protobuf {
	protobuf := new(Protobuf)
	return protobuf
}

type Talker struct {
	Mode       int // Send/Recieve Mode
	Local      Protobuf
	Remote     Protobuf
	Connection *TCPConnection
}

func NewTalker(conn *TCPConnection) *Talker {
	talker := new(Talker)
	talker.Connection = conn
	return talker
}

func (t *Talker) Receive(message []byte, message_id uint32) error {

	err := pr.Unmarshal(message, &t.Remote.Ack)

	t.Remote.mutex.Lock()
	defer t.Remote.mutex.Unlock()
	if err != nil {
		err = pr.Unmarshal(message, &t.Remote.Cmd)
		if err != nil {
			err = pr.Unmarshal(message, &t.Remote.CmdBufferQueue)
			if err != nil {
				err = pr.Unmarshal(message, &t.Remote.SyncCmdBufferQueue)
				if err != nil {
					err = pr.Unmarshal(message, &t.Remote.DataBuffer)
					if err != nil {
						err = pr.Unmarshal(message, &t.Remote.MultiDataBuffer)
						if err != nil {
							err = pr.Unmarshal(message, &t.Remote.Start)
							if err != nil {
								err = pr.Unmarshal(message, &t.Remote.Stop)
								if err != nil {
									err = pr.Unmarshal(message, &t.Remote.GetProfile)
									if err != nil {
										err = pr.Unmarshal(message, &t.Remote.DeviceProfile)
										if err != nil {
											err = pr.Unmarshal(message, &t.Remote.Error)
											if err != nil {
												return fmt.Errorf("Unmarshal error: %v", err)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return err
}

func (t *Talker) RecieveAck(msg []byte, message_id uint32) (pr.Message, error) {
	var ack proto.Ack
	t.Remote.mutex.Lock()
	defer t.Remote.mutex.Unlock()
	err := pr.Unmarshal(msg, &ack)
	return &ack, err
}

func (t *Talker) RecieveDataBuffer(msg []byte, message_id uint32) (pr.Message, error) {
	var dataBuffer proto.DataBuffer
	t.Remote.mutex.Lock()
	defer t.Remote.mutex.Unlock()
	err := pr.Unmarshal(msg, &dataBuffer)
	return &dataBuffer, err
}

func (t *Talker) RecieveMultiDataBuffer(msg []byte, message_id uint32) (pr.Message, error) {
	var multiDataBuffer proto.MultiDataBuffer
	t.Remote.mutex.Lock()
	defer t.Remote.mutex.Unlock()
	err := pr.Unmarshal(msg, &multiDataBuffer)
	return &multiDataBuffer, err
}

func (t *Talker) RecieveError(msg []byte, message_id uint32) (pr.Message, error) {
	var error proto.Error
	t.Remote.mutex.Lock()
	defer t.Remote.mutex.Unlock()
	err := pr.Unmarshal(msg, &error)
	return &error, err
}

func (t *Talker) RecieveDeviceProfile(msg []byte, message_id uint32) (pr.Message, error) {
	var deviceProfile proto.DeviceProfile
	t.Remote.mutex.Lock()
	defer t.Remote.mutex.Unlock()
	err := pr.Unmarshal(msg, &deviceProfile)
	return &deviceProfile, err
}

func (t *Talker) SendCommand(message_id uint32) error {
	t.Local.Cmd.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.Cmd)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err

}

func (t *Talker) SendCommandBufferQueue(message_id uint32) error {
	t.Local.CmdBufferQueue.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.CmdBufferQueue)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendSyncCommandBufferQueue(message_id uint32) error {
	t.Local.SyncCmdBufferQueue.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.SyncCmdBufferQueue)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendStart(message_id uint32) error {
	t.Local.Start.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.Ack)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendStop(message_id uint32) error {
	t.Local.Stop.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.Stop)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendGetProfile(message_id uint32) error {
	t.Local.GetProfile.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.GetProfile)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendAck(message_id uint32) error {
	t.Local.Ack.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.Ack)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendDataBuffer(message_id uint32) error {
	t.Local.DataBuffer.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.DataBuffer)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendMultiDataBuffer(message_id uint32) error {
	t.Local.MultiDataBuffer.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.MultiDataBuffer)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendError(message_id uint32) error {
	t.Local.Error.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.Error)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) SendGetDeviceProfile(message_id uint32) error {
	t.Local.GetProfile.MessageId = message_id
	bytes, err := pr.Marshal(&t.Local.DeviceProfile)
	if err != nil {
		_, err = t.Connection.Send(bytes)
	}

	return err
}

func (t *Talker) Listen(m chan int, message_id uint32, handler func([]byte, uint32) (pr.Message, error)) {
	go t.Connection.Listen(m, message_id, 1024, handler)
}
