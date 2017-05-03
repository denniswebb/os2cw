package main

import (
	"github.com/denniswebb/os2cw/cmd"
)

var BuildVersion = "0.1"

func main() {
	cmd.BuildVersion = BuildVersion
	cmd.Execute()
}
