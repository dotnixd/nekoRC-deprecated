package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/logrusorgru/aurora"
)

// SockAddr is socket path
const SockAddr = "/run/nekoRC/ctl.sock"

// ServerListener is nekoCTL listener
var ServerListener net.Listener

func startServer() {
	os.RemoveAll(SockAddr)

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		_check(errors.New("Couldn`t to listen /run/nekoRC/ctl.sock. nekoCTL shouldn`t working"), false)
	}

	ServerListener = l
}

func listenServer(cfg Config) {
	for {
		conn, err := ServerListener.Accept()
		_check(err, false)

		go handleConnection(conn, cfg)
	}
}

func handleConnection(c net.Conn, cfg Config) {
	c.SetDeadline(time.Now().Add(time.Second * time.Duration(cfg.CTL.ConnectionTimeout)))

	body, err := bufio.NewReader(c).ReadString('\n')
	_check(err, false)

	action := strings.Fields(strings.ReplaceAll(body, "\n", ""))

	if len(action) == 0 {
		_, err := c.Write([]byte(fmt.Sprintln(prefixWarning, "Something went wrong")))
		_check(err, false)
		c.Close()
		return
	}

	data := ""
	var buff bytes.Buffer

	switch action[0] {
	case "help":
	case "reboot":
		stopServices()
		_check(syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART), true)
	case "shutdown":
		stopServices()
		syscall.Sync()
		_check(syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF), true)
	case "service":
		if len(action) != 1 {
			switch action[1] {
			case "start":
				if len(action) != 2 {
					for _, v := range action[2:] {
						srv, err := searchService(v)
						if err != nil {
							data += fmt.Sprintln(prefixWarning, err)
						} else {
							srv.Start(&buff)
							data += buff.String()
						}
					}
				}
			case "stop":
				if len(action) != 2 {
					for _, v := range action[2:] {
						srv, err := searchService(v)
						if err != nil {
							data += fmt.Sprintln(prefixWarning, err)
						} else {
							srv.Stop(&buff)
							data += buff.String()
						}
					}
				}
			case "restart":
				if len(action) != 2 {
					for _, v := range action[2:] {
						srv, err := searchService(v)
						if err != nil {
							data += fmt.Sprintln(prefixWarning, err)
						} else {
							srv.Stop(&buff)
							data += buff.String()
							srv.Start(&buff)
							data += buff.String()
						}
					}
				}
			case "status":
				if len(action) != 2 {
					for _, v := range action[2:] {
						srv, err := searchService(v)
						if err != nil {
							data += fmt.Sprintln(prefixWarning, err)
						} else {
							srv.Stop(&buff)

							data += fmt.Sprintln(prefixInfo, aurora.Magenta(srv.Name), "status:", aurora.Yellow(srv.Status.String()))
						}
					}
				}
			case "list":
				if len(action) != 2 {
					switch action[2] {
					// running|stopped|starting|errored
					case "running":
						for _, s := range Services {
							if s.Status == StatusServiceUp {
								data += fmt.Sprintln(prefixInfo, aurora.Magenta(s.Name), "is running")
							}
						}
					case "stopped":
						for _, s := range Services {
							if s.Status == StatusServiceDown {
								data += fmt.Sprintln(prefixInfo, aurora.Magenta(s.Name), "is stopped")
							}
						}
					case "starting":
						for _, s := range Services {
							if s.Status == StatusServiceStarting {
								data += fmt.Sprintln(prefixInfo, aurora.Magenta(s.Name), "is starting")
							}
						}
					case "errored":
						for _, s := range Services {
							if s.Status == StatusServiceError {
								data += fmt.Sprintln(prefixInfo, aurora.Magenta(s.Name), "is errored")
							}
						}
					}
				}
			}
		}
	case "autostart":
		if len(action) != 1 {
			switch action[1] {
			case "add":
				if len(action) != 2 {
					for _, name := range action[2:] {
						var tmp Service
						if tmp.Load("/etc/nekoRC/services/"+name+".neko.yml") != nil {
							data += fmt.Sprintln(prefixWarning, aurora.Magenta(name), aurora.White("- no such service"))
						} else {
							inittab, err := os.OpenFile("/etc/nekoRC/inittab.neko.yml", os.O_APPEND|os.O_WRONLY, 0644)
							_check(err, false)
							inittab.WriteString("\n - " + name)
							inittab.Close()
							data = fmt.Sprintln(prefixWarning, aurora.White("Success!"))
						}
					}
				}
			case "remove":
				if len(action) != 2 {
					for _, name := range action[2:] {
						var tmp Service
						if tmp.Load("/etc/nekoRC/services/"+name+".neko.yml") != nil {
							data += fmt.Sprintln(prefixWarning, aurora.Magenta(name), aurora.White("- no such service"))
						} else {
							inittab, err := ioutil.ReadFile("/etc/nekoRC/inittab.neko.yml")
							_check(err, false)

							ioutil.WriteFile("/etc/nekoRC/inittab.neko.yml", []byte(strings.ReplaceAll(string(inittab), " - "+name, "")), 644)
							data = fmt.Sprintln(prefixWarning, aurora.White("Success!"))
						}
					}
				}
			case "list":
				inittab, err := ioutil.ReadFile("/etc/nekoRC/inittab.neko.yml")
				_check(err, false)
				var tab []string
				_check(yaml.Unmarshal(inittab, tab), false)

				for _, v := range tab {
					data += fmt.Sprintln(prefixWarning, v, " - in autostart")
				}
			}
		}
	default:
		data = fmt.Sprintln(prefixWarning, "Something went wrong")
	}

	_, err = c.Write([]byte(data))
	_check(err, false)
	c.Close()
}
