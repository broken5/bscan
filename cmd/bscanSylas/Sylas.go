package main

import (
	"bscan/config"
	"bscan/runner/alive"
	"time"
)

func main() {
	options := config.ParseOptions()
	options.Path = "/"
	options.Ports = []int{443, 80}
	options.TargetMode = 2
	for true {
		alive.SwitchMode(options)
		time.Sleep(time.Minute)
	}
}
