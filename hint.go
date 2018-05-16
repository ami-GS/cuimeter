package cuimeter

import (
	"time"
)

type Hint interface {
	read() (string, error)
	parse(string) (interface{}, error)
	postProcess(interface{}) interface{}
	getUnit() string
	getInterval() time.Duration
	// currently allows int64 or []int64
	getChan() chan interface{}
}

type BaseHint struct {
	unit     string
	interval time.Duration
	Chan     chan interface{}
}

func NewBaseHint(unit string, interval time.Duration) *BaseHint {
	return &BaseHint{
		unit:     unit,
		interval: interval,
		Chan:     make(chan interface{}),
	}
}

func (b *BaseHint) getUnit() string {
	return b.unit
}
func (b *BaseHint) getInterval() time.Duration {
	return b.interval
}
func (b *BaseHint) getChan() chan interface{} {
	return b.Chan
}

func (b *BaseHint) read() (string, error) {
	return "", nil
}

func (b *BaseHint) parse(string) (interface{}, error) {
	return 0, nil
}
func (b *BaseHint) postProcess(data interface{}) interface{} {
	return data
}
