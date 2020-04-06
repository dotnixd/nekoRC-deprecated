package main

import (
	"fmt"
	"net"
	"os"

	"github.com/logrusorgru/aurora"
)

// SockAddr is socket path
const SockAddr = "/run/nekoRC/ctl.sock"

var prefixError string = fmt.Sprint(aurora.Red("ERROR")) +
	fmt.Sprint(aurora.Gray(12, " ::"))

func _check(err error) {
	if err != nil {
		fmt.Println(prefixError, aurora.White(err))
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		help(pageAll)
		os.Exit(1)
	}

	c, err := net.Dial("unix", SockAddr)
	_check(err)

	defer c.Close()

	var s Settings
	parseArgs(&s)

	worker(&s, &c)
}

const pageAll = 0
const pageService = 1
const pageAutostart = 2

func help(page uint8) {
	switch page {
	case pageAll:
		fmt.Println("nekoCTL for nekoRC")
		fmt.Println("\tshutdown - shutdown system")
		fmt.Println("\treboot - reboot system")
		help(pageService)
		help(pageAutostart)
		fmt.Println("\thelp - print this message")
	case pageService:
		fmt.Println("\tservice list (running|stopped|starting|errored) - list services")
		fmt.Println("\tservice (start|stop|restart|status) <service(s)> - manipulate with services")
	case pageAutostart:
		fmt.Println("\tautostart list - list services in inittab.neko.yml")
		fmt.Println("\tautostart (add|remove) <service(s)> - add|remove service to|from inittab")
	}
}
