package main

import (
	"bscan/config"
	"bscan/runner/alive"
	"bscan/runner/poc"
)

func main() {
	options := config.ParseOptions()
	alive.SwitchMode(options)
	// 扫描存活网址
	poc.PocScan(options, alive.AliveWebList)

}
