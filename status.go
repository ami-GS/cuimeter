package cuimeter

type Status struct {
	Data   *Queue
	getIdx int
}

func NewStatus(size int) *Status {
	return &Status{
		Data:   NewQueue(size, 1),
		getIdx: 0,
	}
}

func (s *Status) SetData(data interface{}) {
	if s.Data.IsFull() {
		_ = s.Data.Dequeue()
	}
	s.Data.Enqueue(data)
}

func (s *Status) SeekToHead() {
	s.getIdx = s.Data.Head
}

func (s *Status) GetData() interface{} {
	out := s.Data.Data[s.getIdx]
	if s.getIdx == len(s.Data.Data)-1 {
		s.getIdx = 0
	} else {
		s.getIdx++
	}
	return out
}
