package cuimeter

type TrackDirection byte

const (
	TrackNone TrackDirection = iota
	TrackMax
	TrackMin
)

type Queue struct {
	Head   int
	Tail   int
	Data   []int64
	TrackQ *Queue
	track  TrackDirection // 0 none, 1 max, 2 min
}

func NewQueue(size int, track TrackDirection) *Queue {
	TrackQ := (*Queue)(nil)
	if track != TrackNone {
		TrackQ = NewQueue(size, TrackNone)
	}
	return &Queue{
		Head:   0,
		Tail:   0,
		Data:   make([]int64, size),
		TrackQ: TrackQ,
		track:  track,
	}

}

func (q *Queue) Len() int {
	if q.Tail >= q.Head {
		return q.Tail - q.Head
	}
	return len(q.Data) - (q.Head - q.Tail)
}

func (q *Queue) IsFull() bool {
	return (q.Tail+1)%len(q.Data) == q.Head
}

func (q *Queue) IsEmpty() bool {
	return q.Tail == q.Head
}

func (q *Queue) Enqueue(s int64) int {
	if q.IsFull() {
		return -1
	}
	q.Data[q.Tail] = s
	q.Tail++
	if q.Tail == len(q.Data) {
		q.Tail = 0
	}

	if q.track != TrackNone {
		tLen := q.TrackQ.Len()
		for i := 0; i < tLen; i++ {
			val := q.TrackQ.TailData()
			if (q.track == TrackMax && val < s) ||
				(q.track == TrackMin && val > s) {
				q.TrackQ.Tail--
				if q.TrackQ.Tail == -1 {
					q.TrackQ.Tail = q.TrackQ.Len() - 1
				}
			}
		}
		q.TrackQ.Enqueue(s)
	}

	return 1
}

func (q *Queue) TailData() int64 {
	if q.Tail != 0 {
		return q.Data[q.Tail-1]
	}
	return q.Data[len(q.Data)-1]
}

func (q *Queue) HeadData() int64 {
	return q.Data[q.Head]
}

func (q *Queue) Dequeue() int64 {
	if q.IsEmpty() {
		return -1
	}
	dat := q.Data[q.Head]
	q.Head++
	if q.Head == len(q.Data) {
		q.Head = 0
	}

	if q.track != TrackNone {
		if q.Data[q.Head] == q.TrackQ.Data[q.TrackQ.Head] {
			q.TrackQ.Dequeue()
		}
	}
	return dat
}
