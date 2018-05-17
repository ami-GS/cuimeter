package cuimeter

import (
	"bytes"
	"fmt"
	"reflect"
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

func (g *Graph) GetGlobalMax() interface{} {
	// TODO: need to Max/Min
	dat := g.AllStatus[0].Data.TrackQ.HeadData()
	switch data := dat.(type) {
	case int64:
		globalMax := data
		for i := 1; i < len(g.AllStatus); i++ {
			dat := g.AllStatus[i].Data.TrackQ.HeadData()
			if dat.(int64) > globalMax {
				globalMax = dat.(int64)
			}
		}
		return globalMax
	case float64:
		globalMax := data
		for i := 1; i < len(g.AllStatus); i++ {
			dat := g.AllStatus[i].Data.TrackQ.HeadData()
			if dat.(float64) > globalMax {
				globalMax = dat.(float64)
			}
		}
		return globalMax
	default:
		fmt.Println("not supported")
	}
	panic("")
}

func (g *Graph) ScaledHeight(dat, globalMax interface{}, height int) int {
	switch data := dat.(type) {
	case int64:
		assertMax, ok := globalMax.(int64)
		if !ok {
			fmt.Println("type doesn't match")
		}
		return int(float64(data) / float64(assertMax) * float64(height))
	case float64:
		assertMax, ok := globalMax.(float64)
		if !ok {
			fmt.Println("type doesn't match")
		}
		return int(data / assertMax * float64(height))
	default:
		fmt.Printf("data type [%v] is not supported", reflect.TypeOf(data))
	}
	panic("")
}

func (g *Graph) FillBuff() {
	// TODO: needs optimization
	globalMax := g.GetGlobalMax()

	height := int(g.Height)
	width := int(g.Width)

	// TODO: will be removed for optimization
	for h := 0; h < height; h++ {
		for w := width - 2; w >= 0; w-- {
			g.Buff[h][w] = ' '
		}
	}

	for i, st := range g.AllStatus {
		st.SeekToHead()
		for w := width - 2; st.Data.Tail != (st.getIdx+1)%len(st.Data.Data); w-- {
			dat := st.GetData()
			localHeight := g.ScaledHeight(dat, globalMax, height)

			for h := 0; h <= localHeight; h++ {
				if g.Buff[h][w] < rune(i) {
					// for the part of overrapping
					g.Buff[h][w] = rune(len(g.AllStatus) + i)
				} else {
					g.Buff[h][w] = rune(i)
				}
			}
		}
	}
}

func (g *Graph) ShowLabel(unit string, interval time.Duration) {
	var lineBuffer bytes.Buffer
	for i, status := range g.AllStatus {
		lineBuffer.WriteString(fmt.Sprintf("%s [%v %s/%.2fs]  ",
			g.Targets[i],
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

func (g *Graph) Set(status *Status, Chan chan interface{}, wg *sync.WaitGroup) {
	p := <-Chan
	switch dat := p.(type) {
	case int64:
		status.SetData(dat)
	default:
		fmt.Printf("the %v type is not supported yet\n", reflect.TypeOf(dat))
	}
	wg.Done()
}

func (g *Graph) SetForPipe(Chan chan interface{}) {
	p := <-Chan
	switch dat := p.(type) {
	case int64, float64:
		g.AllStatus[0].SetData(dat)
	case []int64:
		for i, st := range g.AllStatus {
			st.SetData(dat[i])
		}
	case []float64:
		for i, st := range g.AllStatus {
			st.SetData(dat[i])
		}
	default:
		fmt.Printf("the %v type is not supported yet\n", reflect.TypeOf(dat))
	}
}

func (g *Graph) Run(hints []Hint) {
	g.runWithInterval(hints)
}

func (g *Graph) runWithInterval(hints []Hint) {
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

func (g *Graph) RunWithPipe(hint Hint) {
	count := uint64(0)
	before := time.Now()
	for {
		go g.Get(hint)
		g.SetForPipe(hint.getChan())

		g.Visualize()
		after := time.Now()
		g.ShowLabel(hint.getUnit(), after.Sub(before))
		count++
		before = after
	}
}
