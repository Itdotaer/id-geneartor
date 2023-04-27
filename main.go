package main

import (
	"flag"
	"fmt"
	"github.com/itdotaer/id-generator/generator"
	"os"
	"runtime"
)

var (
	configFile string
)

func initCmd() {
	flag.StringVar(&configFile, "config", "./alloc.json", "where alloc.json is.")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	initEnv()
	initCmd()

	var err error = nil

	if err = generator.LoadConf(configFile); err != nil {
		goto ERROR
	}
	if err = generator.InitMysql(); err != nil {
		goto ERROR
	}

	if err = StartServer(); err != nil {
		goto ERROR
	}

	os.Exit(0)
ERROR:
	fmt.Println(err)
	os.Exit(-1)
}
