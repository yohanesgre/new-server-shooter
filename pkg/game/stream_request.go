package game

import (
	"sync"
)

// StreamRequest for Incoming StreamRequest
type StreamRequest struct {
	mu                  sync.Mutex
	queue               []Request
	max_packet_in_queue int
}

func NewStream(_max int) *StreamRequest {
	s := new(StreamRequest)
	s.max_packet_in_queue = _max
	return s
}

func (s *StreamRequest) Read() Request {
	var p Request
	s.mu.Lock()
	defer s.mu.Unlock()
	if 0 < len(s.queue) {
		p = Peek(s.queue)
		s.queue = Dequeue(s.queue)
	}
	return p
}

func (s *StreamRequest) ReadAll() *[]Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	q_packet := s.queue

	s.queue = nil
	return &q_packet
}

func (s *StreamRequest) SoftRead() Request {
	var p Request
	s.mu.Lock()
	defer s.mu.Unlock()
	if 0 < len(s.queue) {
		p = Peek(s.queue)
	}
	return p
}

func (s *StreamRequest) Write(_p Request) *StreamRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = Enqueue(s.queue, _p)
	if s.max_packet_in_queue < len(s.queue) {
		s.Pop()
	}
	return s
}

func (s *StreamRequest) WriteAll(packets []Request) *StreamRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, packet := range packets {
		s.queue = Enqueue(s.queue, packet)
	}
	for s.max_packet_in_queue < len(s.queue) {
		s.Pop()
	}
	return s
}

func (s *StreamRequest) Pop() *StreamRequest {
	s.Read()
	return s
}

func (s *StreamRequest) GetLen() int {
	return len(s.queue)
}
