package cuimeter

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Hint interface {
	Parse(data string) (map[string]int64, error)
	// http get, or read line and so on, user can define the method
	Get(retData *int64, wg *sync.WaitGroup)
	GetUnit() string
	GetInterval() time.Duration
}

type BaseHint struct {
	unit     string
	interval time.Duration
}

func NewBaseHint(unit string, interval time.Duration) *BaseHint {
	return &BaseHint{
		unit:     unit,
		interval: interval,
	}
}

func (b *BaseHint) GetUnit() string {
	return b.unit
}
func (b *BaseHint) GetInterval() time.Duration {
	return b.interval
}

// One of example
type NginxStubHint struct {
	*BaseHint
	dataLocation map[string]int
	targetKey    string
	targetFile   string
	unit         string
	interval     time.Duration
}

func NewNginxStatusHint(targetFile string, targetKey string, dataLocation map[string]int, unit string) *NginxStubHint {
	keys := make([]string, 0, len(dataLocation))
	for key := range dataLocation {
		keys = append(keys, key)
	}
	return &NginxStubHint{
		BaseHint:     NewBaseHint(unit, 200*time.Millisecond),
		dataLocation: dataLocation,
		targetFile:   targetFile,
		targetKey:    targetKey,
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
