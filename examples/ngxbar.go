package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// nginx stub format
/*
Active connections: 1
server accepts handled requests
 280006 280006 280006
Reading: 0 Writing: 1 Waiting: 0
*/

type Status struct {
	ActiveConnections int64
	Accepts           uint64
	Handled           uint64
	Requests          uint64
	Reading           uint64
	Writing           int64
	Waiting           int64
}

type TrackStatus struct {
	StatusNow *Status
	StatusPre *Status
}

func NewTrackStatus(target string) *TrackStatus {
	s := &TrackStatus{
		StatusNow: NewStatus(0, 0, 0, 0, 0, 0, 0),
		StatusPre: NewStatus(0, 0, 0, 0, 0, 0, 0),
	}
	err := GetStatus(target, s.StatusNow)
	if err != nil {
		panic(err)
	}
	err = GetStatus(target, s.StatusPre)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *TrackStatus) Sub() *Status {
	return s.StatusNow.Sub(s.StatusPre)
}

func (s *TrackStatus) StoreNow() {
	*s.StatusPre = *s.StatusNow
}

func NewStatus(a int64, b, c, d, e uint64, f, g int64) *Status {
	return &Status{
		ActiveConnections: a,
		Accepts:           b,
		Handled:           c,
		Requests:          d,
		Reading:           e,
		Writing:           f,
		Waiting:           g,
	}
}

func (s *Status) Sub(right *Status) *Status {
	return NewStatus(
		s.ActiveConnections-right.ActiveConnections,
		s.Accepts-right.Accepts,
		s.Handled-right.Handled,
		s.Requests-right.Requests,
		s.Reading-right.Reading,
		s.Writing-right.Writing,
		s.Waiting-right.Waiting,
	)
}

func GetStatus(target string, status *Status) error {
	resp, _ := http.Get(target)
	data, _ := ioutil.ReadAll(resp.Body)

	sp := strings.Split(string(data), " ")
	dat, err := strconv.Atoi(sp[2])
	status.ActiveConnections = int64(dat)
	dat, err = strconv.Atoi(sp[7])
	status.Accepts = uint64(dat)
	dat, err = strconv.Atoi(sp[8])
	status.Handled = uint64(dat)
	dat, err = strconv.Atoi(sp[9])
	status.Requests = uint64(dat)
	dat, err = strconv.Atoi(sp[11])
	status.Reading = uint64(dat)
	dat, err = strconv.Atoi(sp[13])
	status.Writing = int64(dat)
	dat, err = strconv.Atoi(sp[15])
	status.Waiting = int64(dat)

	return err
}

type Queue struct {
	Head   int
	Tail   int
	Data   []*Status
	TrackQ *Queue
	track  byte // 0 none, 1 max, 2 min
}

func NewQueue(size int, track byte) *Queue {
	TrackQ := (*Queue)(nil)
	if track != 0 {
		TrackQ = NewQueue(size, 0)
	}
	return &Queue{
		Head:   0,
		Tail:   0,
		Data:   make([]*Status, size),
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

func (q *Queue) Enqueue(s *Status) int {
	if q.IsFull() {
		return -1
	}
	q.Data[q.Tail] = s
	q.Tail++
	if q.Tail == len(q.Data) {
		q.Tail = 0
	}

	// Track Max
	if q.track != 0 {
		tLen := q.TrackQ.Len()
		for i := 0; i < tLen; i++ {
			val := q.TrackQ.TailData()
			if val.Requests < s.Requests {
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

func (q *Queue) TailData() *Status {
	if q.Tail != 0 {
		return q.Data[q.Tail-1]
	}
	return q.Data[len(q.Data)-1]
}

func (q *Queue) HeadData() *Status {
	return q.Data[q.Head]
}

func (q *Queue) Dequeue() *Status {
	if q.IsEmpty() {
		return nil
	}
	dat := q.Data[q.Head]
	q.Head++
	if q.Head == len(q.Data) {
		q.Head = 0
	}

	// Dequeue Max if needed
	if q.track != 0 {
		if q.Data[q.Head] == q.TrackQ.Data[q.TrackQ.Head] {
			q.TrackQ.Dequeue()
		}
	}
	return dat
}

type StatusQueue struct {
	Data *Queue
	Bar  [][]rune
}

func NewStatusQueue(size int) *StatusQueue {
	return &StatusQueue{
		Data: NewQueue(size, 1),
		Bar:  make([][]rune, size), // [width][height]
	}
}

type Graph struct {
	Width     uint16
	Height    uint16
	AllStatus []*StatusQueue
	Buff      [][]rune
}

func NewGraph(numq int) *Graph {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	outs := strings.Split(string(out[:len(out)-1]), " ")
	hr, err := strconv.Atoi(outs[0])
	wr, err := strconv.Atoi(outs[1])
	if err != nil {
		panic(err)
	}
	hr /= 2
	buff := make([][]rune, hr)
	for h := 0; h < hr; h++ {
		buff[h] = make([]rune, wr)
		for w := 0; w < wr-1; w++ {
			buff[h][w] = ' '
		}
		buff[h][wr-1] = '\n'
	}

	status := make([]*StatusQueue, numq)
	for i := 0; i < numq; i++ {
		status[i] = NewStatusQueue(wr)
	}
	return &Graph{
		Width:     uint16(wr),
		Height:    uint16(hr),
		AllStatus: status,
		Buff:      buff,
	}
}

func (g *Graph) Visualize() error {
	g.FillBuff()
	var lineBuffer bytes.Buffer

	lineBuffer.WriteString(fmt.Sprintf("\x1b[%d;0H", g.Height-3))
	for h := int(g.Height) - 1; h >= 0; h-- {
		for w := 0; w < int(g.Width); w++ {
			lineBuffer.WriteString(fmt.Sprintf(colorMap[g.Buff[h][w]]))
		}
	}
	fmt.Println(lineBuffer.String())

	return nil
}

func (g *Graph) FillBuff() {
	// TODO : should be optimized
	globalMax := uint64(0)
	for _, st := range g.AllStatus {
		max := st.Data.TrackQ.HeadData().Requests
		if max > globalMax {
			globalMax = max
		}
	}
	height := int(g.Height)
	width := int(g.Width)

	for i, st := range g.AllStatus {
		data := st.Data.Data
		qIdx := st.Data.Head

		if i == 0 {
			for h := 0; h < height; h++ {
				for w := width - 2; w >= 0; w-- {
					g.Buff[h][w] = ' '
				}
			}
		}
		for w := width - 2; w >= 0; w-- {
			localHeight := int(float64(data[qIdx].Requests) / float64(globalMax) * float64(height))

			for h := 0; h < height; h++ {
				if h < localHeight {
					if g.Buff[h][w] < rune(i) {
						// for the part of overrapping
						g.Buff[h][w] = rune(len(g.AllStatus) + i)
					} else {
						g.Buff[h][w] = rune(i)
					}
				} /* else if h > int(globalMax) {
					g.Buff[h][w] = ' '
				}*/
			}

			if qIdx == len(data)-1 {
				qIdx = 0
			} else if (qIdx+1)%len(data) == st.Data.Tail {
				break
			} else {
				qIdx++
			}
		}
	}
}

func ngxbar(targets []string) {
	wg := &sync.WaitGroup{}

	ts := make([]*TrackStatus, len(targets))
	for i, _ := range ts {
		ts[i] = NewTrackStatus(targets[i])
	}

	graph := NewGraph(len(targets))

	sleep := 200 * time.Millisecond
	count := uint64(0)
	var err error
	for {
		now := time.Now()
		for i, v := range ts {
			wg.Add(1)
			go func(ii int, vv *TrackStatus) {
				err = GetStatus(targets[ii], vv.StatusNow)
				if err != nil {
					panic(err)
				}
				sub := vv.Sub()
				sub.Requests-- //temporally doing for remove this program's access
				if graph.AllStatus[ii].Data.IsFull() {
					_ = graph.AllStatus[ii].Data.Dequeue()
				}
				graph.AllStatus[ii].Data.Enqueue(sub)
				vv.StoreNow()
				wg.Done()
			}(i, v)
		}
		wg.Wait()
		graph.Visualize()
		label := ""
		for _, status := range graph.AllStatus {
			label += fmt.Sprintf("%d req/%.2fs  ", status.Data.TailData().Requests, float64(sleep)/float64(time.Second))
		}
		fmt.Printf("%s\n", label)
		count++
		time.Sleep(sleep - time.Now().Sub(now))
	}
}

type items []string

func (i *items) String() string {
	return fmt.Sprintf("%v", *i)
}
func (i *items) Set(v string) error {
	*i = append(*i, v)
	return nil
}

var colorMap map[rune]string
var targets items

func init() {
	flag.Var(&targets, "target", "")
	flag.Parse()

	colorMap = map[rune]string{
		0:    "\x1b[38;2;255;82;197;48;2;255;82;197m█\x1b[0m",
		1:    "\x1b[38;2;128;200;197;48;2;128;200;197m█\x1b[0m",
		2:    "\x1b[38;2;128;200;197;48;2;255;82;197m▒\x1b[0m",
		3:    "\x1b[38;2;255;82;197;48;2;128;200;197m▒\x1b[0m",
		' ':  " ",
		'\n': "\n",
	}

}

// will be deprecated?
func main() {
	ngxbar(targets)
}
