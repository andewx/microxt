package net

import "time"

type Stopwatch struct {
	StartTime int64
	EndTime   int64
	Interval  int64
}

func NewStopwatch() *Stopwatch {
	sw := new(Stopwatch)
	return sw
}

func (p *Stopwatch) Start() {
	p.StartTime = time.Now().UnixNano()
}

func (p *Stopwatch) Stop() {
	p.EndTime = time.Now().UnixNano()
}

func (p *Stopwatch) GetInterval() int64 {
	return p.EndTime - p.StartTime
}

func (p *Stopwatch) Reset() {
	p.StartTime = 0
	p.EndTime = 0
}
