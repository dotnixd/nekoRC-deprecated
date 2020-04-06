package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/logrusorgru/aurora"
)

// Run command
func Run(cmd string, args ...string) error {
	os.Setenv("PATH", "/sbin:/usr/bin:/bin:/usr/sbin:/usr/local/bin")
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	return c.Run()
}

// RunBackground spawn command in background
func RunBackground(cmd string, args ...string) (int, error) {
	os.Setenv("PATH", "/sbin:/usr/bin:/bin:/usr/sbin:/usr/local/bin")
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr

	err := c.Start()

	return c.Process.Pid, err
}

// Remount partition
func Remount(dir string) {
	die := false
	if dir == "/" {
		die = true
	}
	_check(Run("mount", "-o", "remount,rw", dir), die)
}

func setupFiles() {
	if _, err := os.Stat("/run/nekoRC"); os.IsNotExist(err) {
		err := os.Mkdir("/run/nekoRC", 755)
		_check(err, false)
	}
}

func doImportantThings() {
	fmt.Println(prefixInfo, aurora.Green("Entering stage \"IMPORTANT\""))
	fmt.Println(prefixInfo, aurora.White("Re-mounting root partition..."))
	Remount("/")
	fmt.Println(prefixInfo, aurora.White("Setting up nekoRC files..."))
	setupFiles()
}

func contains(s *[]string, e string) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}
