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
	prevData     int64
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
		prevData:     0,
		targetFile:   targetFile,
		targetKey:    targetKey,
	}
}

func (s *NginxStatusHint) read() (string, error) {
	resp, err := http.Get(s.targetFile)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *NginxStatusHint) parse(dat string) (int64, error) {
	sp := strings.Split(dat, " ")
	out := make(map[string]int64)
	for k, v := range s.dataLocation {
		dat, err := strconv.Atoi(sp[v])
		if err != nil {
			return 0, err
		}
		out[k] = int64(dat)
	}
	return out[s.targetKey], nil
}

func (s *NginxStatusHint) postProcess(dat int64) int64 {
	// initialize
	if s.prevData == 0 {
		s.prevData = dat + 1
	}
	// -1 removes access by this program
	now := dat - s.prevData - 1
	s.prevData = dat
	return now
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
	graph := cuimeter.NewGraph(targets)
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
