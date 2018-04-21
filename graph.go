package cuimeter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var colorMap = map[rune]string{
	0:    "\x1b[38;2;255;82;197;48;2;255;82;197m█\x1b[0m",
	1:    "\x1b[38;2;128;200;197;48;2;128;200;197m█\x1b[0m",
	2:    "\x1b[38;2;128;200;197;48;2;255;82;197m▒\x1b[0m",
	3:    "\x1b[38;2;255;82;197;48;2;128;200;197m▒\x1b[0m",
	' ':  " ",
	'\n': "\n",
}

type Graph struct {
	Width     uint16
	Height    uint16
	AllStatus []*Item
	Buff      [][]rune
	Target    string // tracking key for queue
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

	status := make([]*Item, numq)
	for i := 0; i < numq; i++ {
		status[i] = NewItem(wr)
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

func (g *Graph) Run(hints []Hint) {
	wg := &sync.WaitGroup{}
	count := uint64(0)
	sleep := hints[0].GetInterval()
	DataBuff := make([]int64, len(hints))
	for {
		now := time.Now()
		for i, v := range hints {
			wg.Add(1)
			go v.Get(&DataBuff[i], wg)
		}
		wg.Wait()
		for i, status := range g.AllStatus {
			wg.Add(1)
			go func(status *Item, data int64) {
				if status.Data.IsFull() {
					_ = status.Data.Dequeue()
				}
				status.Data.Enqueue(data)
				wg.Done()
			}(status, DataBuff[i])
		}
		wg.Wait()
		g.Visualize()
		label := ""
		for _, status := range g.AllStatus {
			label += fmt.Sprintf("%d %s/%.2fs  ",
				status.Data.TailData(),
				hints[0].GetUnit(),
				float64(sleep)/float64(time.Second))
		}
		fmt.Printf("%s\n", label)
		count++
		time.Sleep(time.Duration(sleep) - time.Now().Sub(now))
	}
}
