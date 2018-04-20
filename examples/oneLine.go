package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/ami-GS/cuimeter"
)

type OneLineHint struct {
	Data       int64
	targetFile string
	scanner    *bufio.Scanner
	unit       string
}

func NewOneLineHint(targetFile string) *OneLineHint {
	fp, err := os.Open(targetFile)
	// TODO: need fp.Close(), but when?
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fp)
	return &OneLineHint{
		Data:       0,
		targetFile: targetFile,
		scanner:    scanner,
		unit:       "num",
	}
}

func (s *OneLineHint) Parse(data string) (out map[string]int64, err error) {
	dat, err := strconv.Atoi(data)
	out["data"] = int64(dat)
	return out, nil
}

func (s *OneLineHint) Get(retData *int64, wg *sync.WaitGroup) {
	var tmp int
	var err error
	if s.scanner.Scan() {
		tmp, err = strconv.Atoi(s.scanner.Text())
		if err != nil {
			panic(err)
		}
	}
	*retData = int64(tmp)
	wg.Done() // TODO: channel?
}

func (s *OneLineHint) GetUnit() string {
	return s.unit
}

func ngxbar2(targets []string) {
	hints := make([]cuimeter.Hint, len(targets))
	for i, _ := range hints {
		hints[i] = NewOneLineHint(targets[i])
	}
	graph := cuimeter.NewGraph(len(targets))
	graph.Run(hints)
}

var targets items

type items []string

func (i *items) String() string {
	return fmt.Sprintf("%v", *i)
}
func (i *items) Set(v string) error {
	*i = append(*i, v)
	return nil
}
func main() {
	flag.Var(&targets, "target", "")
	flag.Parse()

	ngxbar2(targets)
}
