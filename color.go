package cuimeter

// TODO: more flexible
var ColorMap map[rune]string

func init() {
	ColorMap = map[rune]string{
		0:    "\x1b[38;2;255;82;197;48;2;255;82;197m█\x1b[0m",
		1:    "\x1b[38;2;128;200;197;48;2;128;200;197m█\x1b[0m",
		2:    "\x1b[38;2;128;200;197;48;2;255;82;197m▒\x1b[0m",
		3:    "\x1b[38;2;255;82;197;48;2;128;200;197m▒\x1b[0m",
		' ':  " ",
		'\n': "\n",
	}
}
