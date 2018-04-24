package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ami-GS/cuimeter"
)

type OneLineHint struct {
	*cuimeter.BaseHint
	Data       int64
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
		Data:       0,
		targetFile: targetFile,
		scanner:    scanner,
	}
}

func (s *OneLineHint) Parse(data string) (out map[string]int64, err error) {
	dat, err := strconv.Atoi(data)
	out["data"] = int64(dat)
	return out, nil
}

func (s *OneLineHint) Get(Chan chan int64) {
	var tmp int
	var err error
	if s.scanner.Scan() {
		tmp, err = strconv.Atoi(s.scanner.Text())
		if err != nil {
			panic(err)
		}
	}
	Chan <- int64(tmp)
}
func oneline(targets []string) {
	hints := make([]cuimeter.Hint, len(targets))
	for i, _ := range hints {
		hints[i] = NewOneLineHint(targets[i], 200*time.Millisecond)
	}
	graph := cuimeter.NewGraph(len(targets))
	graph.Run(hints)
}

type items []string

func (i *items) String() string {
	return fmt.Sprintf("%v", *i)
}
func (i *items) Set(v string) error {
	*i = append(*i, v)
	return nil
}
func main() {
	var targets items
	flag.Var(&targets, "target", "")
	flag.Parse()

	oneline(targets)
}
