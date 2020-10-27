package engosdl

import "strconv"

var ider int = 0

func resetIder() {
	ider = 0
}

func nextIder() string {
	ider++
	return strconv.Itoa(ider)
}
