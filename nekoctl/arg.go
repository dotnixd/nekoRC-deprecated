package main

import "os"

// Settings is arguments pool
type Settings struct {
	IsShutdown bool
	IsReboot   bool
	IsService  bool
	IsHelp     bool
	Help       struct {
		IsService   bool
		IsAutostart bool
	}
	Service struct {
		IsList bool
		List   struct {
			IsRunning  bool
			IsStopped  bool
			IsStarting bool
			IsErrored  bool
		}
		IsStart   bool
		IsStop    bool
		IsStatus  bool
		IsRestart bool
		Services  []string
	}
	IsAutostart bool
	Autostart   struct {
		IsList   bool
		IsAdd    bool
		IsRemove bool
		Services []string
	}
}

func parseArgs(s *Settings) {
	switch os.Args[1] {
	case "help":
		s.IsHelp = true
		return
	case "shutdown":
		s.IsShutdown = true
		return
	case "reboot":
		s.IsReboot = true
		return
	case "service":
		if len(os.Args) == 2 {
			s.IsHelp = true
			s.Help.IsService = true
			return
		}
		s.IsService = true
		switch os.Args[2] {
		case "list":
			s.Service.IsList = true
			if len(os.Args) == 3 {
				return
			}
			switch os.Args[3] {
			case "running":
				s.Service.List.IsRunning = true
				return
			case "stopped":
				s.Service.List.IsStopped = true
				return
			case "starting":
				s.Service.List.IsStarting = true
				return
			case "errored":
				s.Service.List.IsErrored = true
				return
			}
		case "start", "stop", "status":
			if len(os.Args) == 3 {
				s.IsHelp = true
				s.Help.IsService = true
				return
			}
			switch os.Args[2] {
			case "start":
				s.Service.IsStart = true
				s.Service.Services = append(s.Service.Services, os.Args[3:]...)
				return
			case "stop":
				s.Service.IsStop = true
				s.Service.Services = append(s.Service.Services, os.Args[3:]...)
				return
			case "status":
				s.Service.IsStatus = true
				s.Service.Services = append(s.Service.Services, os.Args[3:]...)
				return
			case "restart":
				s.Service.IsRestart = true
				s.Service.Services = append(s.Service.Services, os.Args[3:]...)
				return
			}
		default:
			s.IsHelp = true
			s.Help.IsService = true
			return
		}
		return
	case "autostart":
		if len(os.Args) == 2 {
			s.IsHelp = true
			s.Help.IsAutostart = true
			return
		}
		s.IsAutostart = true

		switch os.Args[2] {
		case "list":
			s.Autostart.IsList = true
			if len(os.Args) == 3 {
				return
			}
			s.Autostart.Services = append(s.Autostart.Services, os.Args[3:]...)
		case "add", "remove":
			if len(os.Args) == 3 {
				s.IsHelp = true
				s.Help.IsAutostart = true
				return
			}

			s.Autostart.Services = append(s.Autostart.Services, os.Args[3:]...)

			switch os.Args[2] {
			case "add":
				s.Autostart.IsAdd = true
				return
			case "remove":
				s.Autostart.IsRemove = true
				return
			}
		default:
			s.IsHelp = true
			s.Help.IsAutostart = true
			return
		}
	default:
		s.IsHelp = true
	}
}
