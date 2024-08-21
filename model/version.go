package model

import "strconv"

var versions = []string{
	"24.04",
	"24.02",
	"23.12",
	"23.09",
	"23.07",
	"23.05",
	"23.02",
	"22.12",
	"22.09",
	"22.07",
	"22.05",
	"2019.0.0",
}

var (
	CurrentVersion string = versions[0]
	BuildNumber    string
)

func BuildNumberInt() int {
	n, _ := strconv.Atoi(BuildNumber)
	return n
}
