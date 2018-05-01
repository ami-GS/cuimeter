package cuimeter

type Status struct {
	Data *Queue
}

func NewStatus(size int) *Status {
	return &Status{
		Data: NewQueue(size, 1),
	}
}

func (s *Status) SetData(data int64) {
	if s.Data.IsFull() {
		_ = s.Data.Dequeue()
	}
	s.Data.Enqueue(data)
}
