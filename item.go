package cuimeter

type Item struct {
	Data *Queue
	Bar  [][]rune
}

func NewItem(size int) *Item {
	return &Item{
		Data: NewQueue(size, 1),
		Bar:  make([][]rune, size), // [width][height]
	}
}
