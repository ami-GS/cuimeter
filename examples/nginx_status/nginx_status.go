package main

import (
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

func NewNginxStatusHint(targetFile string, targetKey string, dataLocation map[string]int) *NginxStatusHint {
	keys := make([]string, 0, len(dataLocation))
	for key := range dataLocation {
		keys = append(keys, key)
	}
	return &NginxStatusHint{
		BaseHint:     cuimeter.NewBaseHint(targetFile + targetKey),
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

func (s *NginxStatusHint) parse(dat string) (interface{}, error) {
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

func (s *NginxStatusHint) postProcess(dat interface{}) interface{} {
	data, ok := dat.(int64)
	if !ok {
		return 0
	}
	// initialize
	if s.prevData == 0 {
		s.prevData = data + 1
	}
	// -1 removes access by this program
	now := data - s.prevData - 1
	s.prevData = data
	return now
}

func ngxstatus() {
	targets := cuimeter.Targets
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
		)
	}
	graph := cuimeter.NewGraph(hints, "req", 200*time.Millisecond)
	graph.Run()
}

func main() {
	ngxstatus()
}
