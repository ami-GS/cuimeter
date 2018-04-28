package cuimeter

import (
	"time"
)

type Hint interface {
	read() (string, error)
	parse(string) (int64, error)
	postProcess(int64) int64
	getUnit() string
	getInterval() time.Duration
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

func (b *BaseHint) getUnit() string {
	return b.unit
}
func (b *BaseHint) getInterval() time.Duration {
	return b.interval
}
func (b *BaseHint) read() (string, error) {
	return "", nil
}
func (b *BaseHint) parse(string) (int64, error) {
	return 0, nil
}
func (b *BaseHint) postProcess(data int64) int64 {
	return data
}
