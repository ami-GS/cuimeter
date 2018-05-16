package cuimeter

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
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
	if terminal.IsTerminal(0) {
		flag.Var(&Targets, "target", "need set at least one")
		flag.Parse()
		if len(Targets) == 0 {
			fmt.Println("At least one --target must be set")
			os.Exit(1)
		}
	} else {
		// for pipe based target
	}
}
