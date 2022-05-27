package main

import (
	"bscan/config"
	"bscan/runner/alive"
	"time"
)

func main() {
	for true {
		run()
		time.Sleep(time.Minute)
	}
}

func run() {
	options := config.ParseOptions()
	options.Path = "/"
	options.Ports = []int{443, 80}
	options.TargetMode = 2
	alive.SwitchMode(options)
	//poc.PocScan(options, alive.AliveWebList)
}
