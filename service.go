package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"syscall"

	"github.com/logrusorgru/aurora"

	"github.com/goccy/go-yaml"
)

// ServiceState is service status
type ServiceState uint8

func (s *ServiceState) String() string {
	switch *s {
	case StatusServiceDown:
		return "down"
	case StatusServiceStarting:
		return "starting"
	case StatusServiceUp:
		return "up"
	case StatusServiceError:
		return "errored"
	}

	return "undefined"
}

const (
	// StatusServiceDown service down
	StatusServiceDown = ServiceState(iota)
	// StatusServiceStarting service starting
	StatusServiceStarting
	// StatusServiceUp service up
	StatusServiceUp
	// StatusServiceError error occured when service was starting
	StatusServiceError
)

// Services is array of services
var Services []Service

// Service declares service
type Service struct {
	Name         string   `yaml:"Name"`         // Name of service
	Description  string   `yaml:"Description"`  // Description of service
	Dependencies []string `yaml:"Dependencies"` // Needed services for service
	Startup      string   `yaml:"Startup"`      // Command to start service
	Shutdown     string   `yaml:"Shutdown"`     // Command to shutdown service
	Wait         bool     `yaml:"Wait"`         // Should nekoRC wait when the service starts?
	Status       ServiceState
	PID          int
}

// Load loads service
func (s *Service) Load(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, s)
}

// Start service
func (s *Service) Start(buff *bytes.Buffer) {
	if s.Status == StatusServiceUp {
		return
	}

	if s.Status == StatusServiceStarting {
		for s.Status == StatusServiceStarting {
		}

		return
	}

	s.Status = StatusServiceStarting

	s.SatisfyDependencies()

	fmt.Fprint(buff, prefixInfo, aurora.White(" Starting up "+s.Name+"... "))

	if s.Wait {
		if err := Run("sh", "-c", s.Startup); err != nil {
			s.Status = StatusServiceError
			fmt.Fprintln(buff, prefixWarning, aurora.White(err))
			return
		}
	} else {
		pid, err := RunBackground("sh", "-c", s.Startup)
		if err != nil {
			s.Status = StatusServiceError
			fmt.Fprintln(buff, prefixWarning, aurora.White(err))
			return
		}

		s.PID = pid
	}

	s.Status = StatusServiceUp
	fmt.Fprintln(buff, aurora.Green("DONE"))
}

// Stop service
func (s *Service) Stop(buff *bytes.Buffer) {
	if s.Shutdown != "" {
		fmt.Fprint(buff, prefixInfo, aurora.White(" Killing "+s.Name+"... "))
		syscall.Kill(s.PID, syscall.SIGKILL)
	} else {
		fmt.Fprint(buff, prefixInfo, aurora.White(" Stopping "+s.Name+"... "))
		_check(Run(s.Shutdown), false)
	}

	fmt.Fprintln(buff, aurora.Green("DONE"))
}

// SatisfyDependencies satisfy dependencies
func (s *Service) SatisfyDependencies() {
	for _, service := range Services {
		if contains(&s.Dependencies, service.Name) {
			var buff bytes.Buffer
			service.Start(&buff)
			fmt.Print(buff.String())
		}
	}
}

func loadServices() {
	content, err := ioutil.ReadFile("/etc/nekoRC/inittab.neko.yml")
	_check(err, true)

	var tab []string
	_check(yaml.Unmarshal(content, &tab), true)

	var buff bytes.Buffer

	for _, sym := range tab {
		var current Service
		_check(current.Load("/etc/nekoRC/services/"+sym+".neko.yml"), false)
		current.Status = StatusServiceDown
		current.Start(&buff)
		fmt.Print(buff.String())
		Services = append(Services, current)
	}

	fmt.Println(prefixNeko, aurora.Cyan("DONE UwU!!!"))
}

func stopServices() {
	var buff bytes.Buffer
	for _, s := range Services {
		s.Stop(&buff)
	}

	fmt.Println(prefixInfo, aurora.White("Sending SIGKILL to remaining services..."))
	for _, s := range Services {
		s.Stop(&buff)
	}
}

func searchService(name string) (*Service, error) {
	for _, s := range Services {
		if s.Name == name {
			return &s, nil
		}
	}

	var tmp Service
	if err := tmp.Load("/etc/nekoRC/services/" + name + ".neko.yml"); err != nil {
		return nil, errors.New(name + " - no such service")
	}

	Services = append(Services, tmp)

	return searchService(tmp.Name)
}
