package cuimeter

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: more flexible
var ColorMap map[rune]string

func getColoredString(val rune, fg, bg string) string {
	// 48 for background
	// 38 for foreground
	return fmt.Sprintf("\x1b[38;%s;48;%sm%c\x1b[0m", fg, bg, val)
}

// first color becomes bigger
func overrupColor(first, second string) (string, error) {
	fStrColor := strings.Split(first, ";")[1:]
	sStrColor := strings.Split(second, ";")[1:]
	getRGBbyte := func(strColor []string) (byteColor [3]byte, err error) {
		for i, strColor := range strColor {
			color, err := strconv.Atoi(strColor)
			if err != nil {
				return byteColor, err
			}
			byteColor[i] = byte(color)
		}
		return byteColor, nil
	}

	fColor, err := getRGBbyte(fStrColor)
	if err != nil {
		return "", err
	}
	sColor, err := getRGBbyte(sStrColor)
	if err != nil {
		return "", err
	}

	mColor := [3]byte{}
	for i := 0; i < 3; i++ {
		//mColor[i] = (fColor[i]*2 + sColor[i]) / 3
		mColor[i] = (fColor[i] + sColor[i]) / 2
	}
	return fmt.Sprintf("2;%d;%d;%d", mColor[0], mColor[1], mColor[2]), nil
}

const baseBGColor = "2;100;100;100;197"

var bgColor = []string{
	"2;255;82;197",  // purple
	"2;128;200;197", // right blue
}

const axisColor = "2;200;200;200"

func init() {
	// TODO: change color depends on color order? first and second
	middle01, err := overrupColor(bgColor[0], bgColor[1])
	if err != nil {
		panic(err)
	}
	middle10, err := overrupColor(bgColor[0], bgColor[1])
	if err != nil {
		panic(err)
	}

	ColorMap = map[rune]string{
		1: getColoredString(' ', bgColor[0], bgColor[0]),
		2: getColoredString(' ', bgColor[1], bgColor[1]),
		//3:       getColoredString("▒", bgColor[1], bgColor[0]),
		//4:       getColoredString("▒", bgColor[0], bgColor[1]),
		3:       getColoredString(' ', middle01, middle01),
		4:       getColoredString(' ', middle10, middle10),
		'─':     getColoredString('─', axisColor, baseBGColor),
		'─' + 1: getColoredString('─', axisColor, bgColor[0]),
		'─' + 2: getColoredString('─', axisColor, bgColor[1]),
		'─' + 3: getColoredString('─', axisColor, middle01),
		'─' + 4: getColoredString('─', axisColor, middle10),
		' ':     getColoredString(' ', baseBGColor, baseBGColor),
		'\n':    "\n",
	}
}

func GetFillChar(one rune) string {
	if '.' <= one && one <= '9' || 'a' <= one && one <= 'z' || 'A' <= one && one <= 'Z' {
		return getColoredString(one, axisColor, baseBGColor)
	}
	return ColorMap[one]
}
