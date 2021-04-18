package main

import "fmt"

const (
	_G_HASH = "{_G_HASH}"
	_G_REVS = "{_G_REVS}"
)

func verinfo() string {
	ver := fmt.Sprintf("V%s.%s", _G_REVS, _G_HASH)
	return ver
}
