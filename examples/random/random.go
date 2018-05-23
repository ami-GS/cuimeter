package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/ami-GS/cuimeter"
)

type RandomHint struct {
	*cuimeter.BaseHint
}

func NewRandomHint(target string) *RandomHint {
	return &RandomHint{
		BaseHint: cuimeter.NewBaseHint(target),
	}
}

func (r *RandomHint) read() (string, error) {
	return strconv.FormatFloat(rand.Float64(), 'E', -1, 64), nil
}

func (r *RandomHint) parse(dat string) (interface{}, error) {
	data, _ := strconv.ParseFloat(dat, 64)
	return data, nil
}

func random() {
	hints := make([]cuimeter.Hint, 2)
	hints[0] = NewRandomHint("random:1")
	hints[1] = NewRandomHint("random:2")
	graph := cuimeter.NewGraph(hints, "num", 200*time.Millisecond)
	graph.SetMinMax(0.0, 1.0)
	graph.IsFixed = true
	graph.NumSeparate = 8
	graph.Run()
}

func main() {
	random()
}
