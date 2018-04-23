package cuimeter

import (
	"time"
)

type Hint interface {
	Parse(data string) (map[string]int64, error)
	// http get, or read line and so on, user can define the method
	Get(Chan chan int64)
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
