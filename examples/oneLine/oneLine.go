package main

import (
	"bufio"
	"os"
	"strconv"
	"time"

	"github.com/ami-GS/cuimeter"
)

type OneLineHint struct {
	*cuimeter.BaseHint
	targetFile string
	scanner    *bufio.Scanner
}

func NewOneLineHint(targetFile string, interval time.Duration) *OneLineHint {
	fp, err := os.Open(targetFile)
	// TODO: need fp.Close(), but when?
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fp)
	return &OneLineHint{
		BaseHint:   cuimeter.NewBaseHint("num", interval),
		targetFile: targetFile,
		scanner:    scanner,
	}
}

func (s *OneLineHint) read() (string, error) {
	if s.scanner.Scan() {
		tmp := s.scanner.Text()
		return tmp, nil
	}
	return "", nil
}

func (s *OneLineHint) parse(dat string) (interface{}, error) {
	tmp, err := strconv.Atoi(dat)
	if err != nil {
		return 0, err
	}
	return int64(tmp), nil
}

//func oneline(targets []string) {
func oneline() {
	targets := cuimeter.Targets
	hints := make([]cuimeter.Hint, len(targets))
	for i, _ := range hints {
		hints[i] = NewOneLineHint(targets[i], 200*time.Millisecond)
	}
	graph := cuimeter.NewGraph(targets)
	graph.Run(hints)
}

func main() {
	oneline()
}
