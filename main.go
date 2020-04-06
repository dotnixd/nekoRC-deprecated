package main

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

var prefixNeko string = fmt.Sprint(aurora.Bold(aurora.Magenta("nekoRC"))) + fmt.Sprint(aurora.Gray(12, " ::"))

var prefixInfo string = fmt.Sprint(aurora.Bold(aurora.White("INFO"))) + fmt.Sprint(aurora.Gray(12, " ::"))
var prefixFatal string = fmt.Sprint(aurora.Bold(aurora.Red("ERROR"))) + fmt.Sprint(aurora.Gray(12, " ::"))
var prefixWarning string = fmt.Sprint(aurora.Bold(aurora.Yellow("WARNING"))) + fmt.Sprint(aurora.Gray(12, " ::"))

func main() {
	if os.Getpid() != 1 {
		fmt.Println(prefixFatal, aurora.White("Must be run as PID 1"))
		os.Exit(1)
	}

	doImportantThings()

	var cfg Config
	cfg.Load()

	startServer()
	go listenServer(cfg)

	fmt.Println()
	fmt.Println(prefixNeko, aurora.Cyan(cfg.Distribution))
	fmt.Println()

	fmt.Println(prefixInfo, aurora.Green("Entering stage \"DEFAULT\""))
	loadServices()

	for {
		select {}
	}
}

func _check(err error, die bool) {
	if err != nil {
		if die {
			fmt.Println(prefixFatal, aurora.White(err))
			for {
			}
		} else {
			fmt.Println(prefixWarning, aurora.White(err))
		}
	}
}
