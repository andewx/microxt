package net

import "fmt"

/* We are using a frame class to facilitate streaming UDP connection data,
 this helps us ensure that we are sending and recieving data in a consistent manner
with an expected bit rate*/

// Frame Structure
const FRAME_HEADER_SIZE = 8
const FRAME_MESSAGE_SIZE = 1024

// Frame Structure
type Frame struct {
	Header  []byte
	Message []byte
	free    bool
}

func NewFrame() *Frame {
	frame := new(Frame)
	frame.Header = make([]byte, FRAME_HEADER_SIZE)
	frame.Message = make([]byte, FRAME_MESSAGE_SIZE)
	return frame
}

func (f *Frame) Clear() {
	f.free = true
	copy(f.Header, make([]byte, FRAME_HEADER_SIZE))
}

func (f *Frame) Fill(message []byte) error {

	if f.free {
		return fmt.Errorf("%sError%s Frame is free a full", CS_RED, CS_WHITE)
	}

	d := int32(len(message))
	if d > FRAME_MESSAGE_SIZE {
		return fmt.Errorf("%sError%s Message size %d too large for frame", CS_RED, CS_WHITE, d)
	}
	header := make([]byte, FRAME_HEADER_SIZE)

	// Convert int32 to bytes
	header[0] = byte(d >> 24)
	header[1] = byte(d >> 16)
	header[2] = byte(d >> 8)
	header[3] = byte(d)

	copy(f.Header, header)
	copy(f.Message, message)
	f.free = false
	return nil
}

// Frame Queue
type FrameQueue struct {
	items []*Frame
	back  int
}

/*Frame Queue Cycles (FQC) is the number of frames that can be stored in the queue and uses these as its memory base for sending message */

func NewFrameQueue(cycle int) *FrameQueue {
	q := new(FrameQueue)
	q.items = make([]*Frame, cycle)
	q.back = 0
	return q
}

func (q *FrameQueue) Full() bool {
	return q.back == len(q.items)
}

func (q *FrameQueue) Enqueue(message []byte) error {
	if q.Full() {
		return fmt.Errorf("%sError%s Frame Queue is full", CS_RED, CS_WHITE)
	}
	q.items[q.back].Fill(message)
	q.back += 1
	return nil
}

func (q *FrameQueue) Deqeue() *Frame {
	if q.back == 0 {
		return nil
	}
	q.back -= 1
	temp := q.items[0]
	temp.free = true

	//Shift all items in the queue forward
	for i := 0; i < len(q.items); i++ {
		q.items[i] = q.items[i+1]
	}

	q.items[len(q.items)-1] = temp

	return temp
}
