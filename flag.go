package cuimeter

import (
	"flag"
	"fmt"
	"os"
)

type item []string

func (i *item) String() string {
	return fmt.Sprintf("%v", *i)
}
func (i *item) Set(v string) error {
	*i = append(*i, v)
	return nil
}

var Targets item

func init() {
	flag.Var(&Targets, "target", "need set at least one")
	flag.Parse()
	if len(Targets) == 0 {
		fmt.Println("At least one --target must be set")
		os.Exit(1)
	}
}
