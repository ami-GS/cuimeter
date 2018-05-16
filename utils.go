package cuimeter

import (
	"golang.org/x/crypto/ssh/terminal"
)

func GetDisplayWH() (w int, h int, err error) {
	// 1 can retrieve size even if pipe is used as stdin
	width, height, err := terminal.GetSize(1)
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}
