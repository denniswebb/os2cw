package main

import (
	"github.com/denniswebb/os2cw/cmd"
)

//BuildVersion should be set during the build with -ldflags "-X main.BuildVersion=${VERSION}"
var BuildVersion = "0.1"

func main() {
	cmd.BuildVersion = BuildVersion
	cmd.Execute()
}
