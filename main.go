package main

import (
	"github.com/dhwebb/go-ec2-custom-metrics/cmd"
)

var BuildVersion = "0.1"

func main() {
	cmd.BuildVersion = BuildVersion
	cmd.Execute()
}
