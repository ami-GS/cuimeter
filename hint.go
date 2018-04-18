package cuimeter

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Hint interface {
	Parse(data string) (map[string]int64, error)
	// http get, or read line and so on, user can define the method
	Get() (int64, error)
	GetUnit() string
}

// One of example
type NginxStubHint struct {
	dataLocation map[string]int
	targetKey    string
	targetFile   string
	unit         string
}

func NewNginxStatusHint(targetFile string, targetKey string, dataLocation map[string]int, unit string) *NginxStubHint {
	keys := make([]string, 0, len(dataLocation))
	for key := range dataLocation {
		keys = append(keys, key)
	}
	return &NginxStubHint{
		dataLocation: dataLocation,
		targetFile:   targetFile,
		targetKey:    targetKey,
		unit:         unit,
	}
}

func (s *NginxStubHint) Parse(data string) (out map[string]int64, err error) {
	sp := strings.Split(data, " ")
	for k, v := range s.dataLocation {
		dat, err := strconv.Atoi(sp[v])
		if err != nil {
			return nil, err
		}
		out[k] = int64(dat)
	}
	return out, nil
}

func (s *NginxStubHint) Get() (int64, error) {
	resp, _ := http.Get(s.targetFile)
	data, _ := ioutil.ReadAll(resp.Body)
	out, err := s.Parse(string(data))
	if err != nil {
		panic(err)
	}
	return out[s.targetKey], nil
}

func (s *NginxStubHint) GetUnit() string {
	return s.unit
}
