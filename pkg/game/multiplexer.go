package game

type Multiplexer struct {
	stream *StreamRequest
}

func NewMultiplexer(_max int) *Multiplexer {
	m := new(Multiplexer)
	m.stream = NewStream(_max)
	return m
}

func (m *Multiplexer) GetStream() *StreamRequest {
	return m.stream
}

func (m *Multiplexer) Read() Request {
	return m.GetStream().Read()
}

func (m *Multiplexer) ReadAll() *[]Request {
	return m.GetStream().ReadAll()
}

func (m *Multiplexer) SoftRead() Request {
	return m.GetStream().SoftRead()
}

func (m *Multiplexer) Write(_p Request) *Multiplexer {
	m.GetStream().Write(_p)
	return m
}
