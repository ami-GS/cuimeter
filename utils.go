package cuimeter

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetDisplayWH() (w int, h int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	outs := strings.Split(string(out[:len(out)-1]), " ")
	h, err = strconv.Atoi(outs[0])
	w, err = strconv.Atoi(outs[1])
	if err != nil {
		return 0, 0, err
	}
	return w, h, nil
}
