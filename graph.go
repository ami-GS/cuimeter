package cuimeter

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type Graph struct {
	Width     uint16
	Height    uint16
	AllStatus []*Status
	Buff      [][]rune
	Targets   []string // tracking key for queue
}

func NewGraph(targets []string) *Graph {
	wr, hr, err := GetDisplayWH()
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
	numq := len(targets)
	status := make([]*Status, numq)
	for i := 0; i < numq; i++ {
		status[i] = NewStatus(wr)
	}
	return &Graph{
		Width:     uint16(wr),
		Height:    uint16(hr),
		AllStatus: status,
		Buff:      buff,
		Targets:   targets,
	}
}

func (g *Graph) Visualize() error {
	g.FillBuff()
	var lineBuffer bytes.Buffer

	lineBuffer.WriteString(fmt.Sprintf("\x1b[%d;0H", g.Height-3))
	for h := int(g.Height) - 1; h >= 0; h-- {
		for w := 0; w < int(g.Width); w++ {
			lineBuffer.WriteString(fmt.Sprintf(ColorMap[g.Buff[h][w]]))
		}
	}
	fmt.Println(lineBuffer.String())

	return nil
}

func (g *Graph) GetGlobalMax() (globalMax int64) {
	// TODO: need to Max/Min
	for _, st := range g.AllStatus {
		max := st.Data.TrackQ.HeadData()
		if max > globalMax {
			globalMax = max
		}
	}
	return globalMax
}

func (g *Graph) FillBuff() {
	// TODO: needs optimization
	globalMax := g.GetGlobalMax()

	height := int(g.Height)
	width := int(g.Width)

	for h := 0; h < height; h++ {
		for w := width - 2; w >= 0; w-- {
			g.Buff[h][w] = ' '
		}
	}

	for i, st := range g.AllStatus {
		data := st.Data.Data
		qIdx := st.Data.Head

		for w := width - 2; w >= 0; w-- {
			localHeight := int(float64(data[qIdx]) / float64(globalMax) * float64(height))

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

func (g *Graph) ShowLabel(unit string, interval time.Duration) {
	var lineBuffer bytes.Buffer
	for _, status := range g.AllStatus {
		lineBuffer.WriteString(fmt.Sprintf("%d %s/%.2fs  ",
			status.Data.TailData(),
			unit,
			float64(interval)/float64(time.Second)))
	}
	fmt.Printf("%s\n", lineBuffer.String())
}

func (g *Graph) Get(hint Hint) {
	strData, err := hint.read()
	if err != nil {
		// error channel?
		panic(err)
	}
	data, err := hint.parse(strData)
	if err != nil {
		// error channel?
		panic(err)
	}
	data = hint.postProcess(data)
	hint.getChan() <- data
}

func (g *Graph) Set(status *Status, Chan chan int64, wg *sync.WaitGroup) {
	dat := <-Chan
	status.SetData(dat)
	wg.Done()
}

func (g *Graph) Run(hints []Hint) {
	wg := &sync.WaitGroup{}
	count := uint64(0)
	sleep := hints[0].getInterval()
	for {
		now := time.Now()
		wg.Add(len(hints))
		for i, v := range hints {
			go g.Get(v)
			go g.Set(g.AllStatus[i], v.getChan(), wg)
		}
		wg.Wait()

		g.Visualize()
		g.ShowLabel(hints[0].getUnit(), sleep)
		count++
		time.Sleep(sleep - time.Now().Sub(now))
	}
}
