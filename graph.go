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
	Hints     []Hint
	Unit      string
	Interval  time.Duration

	// log y axis
	IsLog bool
	// used for no scaling
	IsFixed bool
	// used for fixed graph
	Max interface{}
	Min interface{}
	// used for number of axis
	NumSeparate int
}

func NewGraph(hints []Hint, unit string, interval time.Duration) *Graph {
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

	var targets []string
	targetInterface := hints[0].getTarget()
	switch target := targetInterface.(type) {
	case string:
		targets = append(targets, target)
		for i := 1; i < len(hints); i++ {
			ti := hints[i].getTarget()
			targets = append(targets, ti.(string))
		}
	case []string:
		for _, t := range target {
			targets = append(targets, t)
		}
	default:
		panic("target should be string or []string")
	}

	status := make([]*Status, len(targets))
	for i := 0; i < len(targets); i++ {
		status[i] = NewStatus(wr)
	}
	return &Graph{
		Width:       uint16(wr),
		Height:      uint16(hr),
		AllStatus:   status,
		Targets:     targets,
		Buff:        buff,
		Hints:       hints,
		Unit:        unit,
		Interval:    interval,
		NumSeparate: 4, // default
	}
}

func (g *Graph) SetMinMax(min, max interface{}) {
	g.Min = min
	g.Max = max
}

func (g *Graph) Visualize() error {
	g.FillBuff()
	g.FillAxis()
	var lineBuffer bytes.Buffer

	lineBuffer.WriteString(fmt.Sprintf("\x1b[%d;0H", g.Height-3))
	for h := int(g.Height) - 1; h >= 0; h-- {
		for w := 0; w < int(g.Width); w++ {
			lineBuffer.WriteString(fmt.Sprintf(GetFillChar(g.Buff[h][w])))
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
		fmt.Printf("data type [%v] is not supported", reflect.TypeOf(data))
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
	var globalMax interface{}
	if g.IsFixed {
		globalMax = g.Max
	} else {
		globalMax = g.GetGlobalMax()
	}

	height := int(g.Height)
	width := int(g.Width)

	// TODO: will be removed for optimization
	for h := 0; h < height; h++ {
		for w := width - 2; w >= 0; w-- {
			g.Buff[h][w] = ' '
		}
	}

	for i, st := range g.AllStatus {
		colorID := rune(i + 1) // colorMap starts from rune(1)
		st.SeekToHead()
		for w := width - 2; st.Data.Tail != (st.getIdx+1)%len(st.Data.Data); w-- {
			dat := st.GetData()
			localHeight := g.ScaledHeight(dat, globalMax, height)

			for h := 0; h <= localHeight-1; h++ {
				if g.Buff[h][w] < colorID {
					// for the part of overrapping
					g.Buff[h][w] = rune(len(g.AllStatus)) + colorID
				} else {
					g.Buff[h][w] = colorID
				}
			}
		}
	}
}

func (g *Graph) FillAxis() error {
	sep := g.NumSeparate
	putLabel := func(h, i int) {
		var nums string
		if g.IsFixed {
			nums = fmt.Sprintf("%.2f ", (g.Max.(float64)-g.Min.(float64))/float64(sep)*float64(i))
		} else {
			// TODO: dynamically changed
			//nums = fmt.Sprintf("%.2f ", (g.Max.(float64)-g.Min.(float64))/sep*float64(i))
		}
		copy(g.Buff[h][:len(nums)], []rune(nums))
		if h == int(g.Height-1) {
			unit := fmt.Sprintf(" %s", g.Unit)
			//copy(g.Buff[h-1][:len(unit)], []rune(unit))
			copy(g.Buff[h][int(g.Width)-1-len(unit):g.Width-1], []rune(unit))
		}
	}
	fill := func(h, i int) {
		for w := 0; w < int(g.Width)-1; w++ {
			if g.Buff[h][w] == ' ' {
				g.Buff[h][w] = '─'
			} else {
				g.Buff[h][w] += '─'
			}
		}
		putLabel(h, i)
	}

	i := sep
	for h := int(g.Height - 1); h > 0; h -= int(g.Height) / sep {
		fill(h, i)
		i--
	}
	fill(0, i)

	return nil
}

func (g *Graph) ShowLabel(interval time.Duration) {
	var lineBuffer bytes.Buffer
	if interval == 0 {
		interval = g.Interval
	}

	for i, status := range g.AllStatus {
		lineBuffer.WriteString(fmt.Sprintf("%s [%.2v %s/%.2fs]  ",
			g.Targets[i],
			status.Data.TailData(),
			g.Unit,
			float64(interval)/float64(time.Second)))
	}
	fmt.Printf("%s\n", lineBuffer.String())
}

func (g *Graph) Get(hintID int) {
	strData, err := g.Hints[hintID].read()
	if err != nil {
		// error channel?
		panic(err)
	}
	data, err := g.Hints[hintID].parse(strData)
	if err != nil {
		// error channel?
		panic(err)
	}
	data = g.Hints[hintID].postProcess(data)
	g.Hints[hintID].getChan() <- data
}

func (g *Graph) Set(status *Status, Chan chan interface{}, wg *sync.WaitGroup) {
	p := <-Chan
	status.SetData(p)
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

func (g *Graph) Run() {
	g.runWithInterval()
}

func (g *Graph) runWithInterval() {
	wg := &sync.WaitGroup{}
	count := uint64(0)
	sleep := g.Interval
	hintNum := len(g.Hints)
	for {
		now := time.Now()
		wg.Add(hintNum)
		for i, hint := range g.Hints {
			go g.Get(i)
			go g.Set(g.AllStatus[i], hint.getChan(), wg)
		}
		wg.Wait()

		g.Visualize()
		g.ShowLabel(0)
		count++
		time.Sleep(sleep - time.Now().Sub(now))
	}
}

func (g *Graph) RunWithPipe() {
	count := uint64(0)
	before := time.Now()
	for {
		go g.Get(0)
		g.SetForPipe(g.Hints[0].getChan())

		g.Visualize()
		after := time.Now()
		g.ShowLabel(after.Sub(before))
		count++
		before = after
	}
}
