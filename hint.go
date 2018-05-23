package cuimeter

type Hint interface {
	read() (string, error)
	parse(string) (interface{}, error)
	postProcess(interface{}) interface{}
	getTarget() interface{}
	getChan() chan interface{}
}

type BaseHint struct {
	target interface{} // to use string and []string
	Chan   chan interface{}
}

func NewBaseHint(target interface{}) *BaseHint {
	return &BaseHint{
		target: target,
		Chan:   make(chan interface{}),
	}
}

var targetIdx int

func (b *BaseHint) getTarget() interface{} {
	return b.target
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
