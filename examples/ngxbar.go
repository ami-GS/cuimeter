package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ami-GS/cuimeter"
)

type NginxStatusHint struct {
	*cuimeter.BaseHint
	dataLocation map[string]int
	Data         int64
	targetFile   string
	targetKey    string
}

func NewNginxStatusHint(targetFile string, targetKey string, dataLocation map[string]int, unit string) *NginxStatusHint {
	keys := make([]string, 0, len(dataLocation))
	for key := range dataLocation {
		keys = append(keys, key)
	}
	return &NginxStatusHint{
		BaseHint:     cuimeter.NewBaseHint(unit, 200*time.Millisecond),
		dataLocation: dataLocation,
		Data:         0,
		targetFile:   targetFile,
		targetKey:    targetKey,
	}
}

func (s *NginxStatusHint) Parse(data string) (out map[string]int64, err error) {
	sp := strings.Split(data, " ")
	out = make(map[string]int64)
	for k, v := range s.dataLocation {
		dat, err := strconv.Atoi(sp[v])
		if err != nil {
			return nil, err
		}
		out[k] = int64(dat)
	}
	return out, nil
}

func (s *NginxStatusHint) Get(Chan chan int64) {
	resp, _ := http.Get(s.targetFile)
	data, _ := ioutil.ReadAll(resp.Body)

	out, err := s.Parse(string(data))
	if err != nil {
		panic(err)
	}

	// initialize
	if s.Data == 0 {
		s.Data = out[s.targetKey] + 1
	}
	// -1 removes access by this program
	now := out[s.targetKey] - s.Data - 1
	s.Data = out[s.targetKey]
	Chan <- now
}
func ngxstatus(targets []string) {
	hints := make([]cuimeter.Hint, len(targets))
	for i, _ := range hints {
		hints[i] = NewNginxStatusHint(targets[i],
			"Requests",
			map[string]int{
				"Active Connections": 2,
				"Accepts":            7,
				"Handled":            8,
				"Requests":           9,
				"Reading":            11,
				"Writing":            13,
				"Waiting":            15,
			},
			"req",
		)
	}
	graph := cuimeter.NewGraph(len(targets))
	graph.Run(hints)
}

type items []string

func (i *items) String() string {
	return fmt.Sprintf("%v", *i)
}
func (i *items) Set(v string) error {
	*i = append(*i, v)
	return nil
}
func main() {
	var targets items
	flag.Var(&targets, "target", "")
	flag.Parse()

	ngxstatus(targets)
}
