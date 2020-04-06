package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func worker(s *Settings, c *net.Conn) {
	data := ""
	if s.IsHelp {
		if s.Help.IsAutostart {
			help(pageAutostart)
			_, err := (*c).Write([]byte(data + "help\n"))
			_check(err)
			_, err = bufio.NewReader(*c).ReadString('\n')
			_check(err)
			return
		}
		if s.Help.IsService {
			help(pageService)
		} else {
			help(pageAll)
		}
		_, err := (*c).Write([]byte("help\n"))
		_check(err)
		_, err = bufio.NewReader(*c).ReadString('\n')
		_check(err)
		return
	}
	if s.IsShutdown {
		data = "shutdown"
	}
	if s.IsReboot {
		data = "reboot"
	}
	if s.IsService {
		if s.Service.IsList {
			if s.Service.List.IsRunning {
				data = "service list running"
			}
			if s.Service.List.IsStarting {
				data = "service list starting"
			}
			if s.Service.List.IsStopped {
				data = "service list stopped"
			}
			if s.Service.List.IsErrored {
				data = "service list errored"
			} else {
				data = "service list all"
			}
		}

		if s.Service.IsStart {
			data = "service start"
		}
		if s.Service.IsStop {
			data = "service stop"
		}
		if s.Service.IsStatus {
			data = "service status"
		}
		if s.Service.IsRestart {
			data = "service restart"
		}
		if !s.Service.IsList {
			data += " " + strings.Join(s.Service.Services, " ")
		}
	}
	if s.IsAutostart {
		if s.Autostart.IsAdd {
			data = "autostart add"
		}
		if s.Autostart.IsRemove {
			data = "autostart remove"
		}
		if s.Autostart.IsList {
			data = "autostart list"
		}

		if !s.Autostart.IsList {
			data += " " + strings.Join(s.Autostart.Services, " ")
		}
	}

	_, err := (*c).Write([]byte(data + "\n"))
	_check(err)
	res, err := bufio.NewReader(*c).ReadString('\n')
	_check(err)

	fmt.Println(res)
}
