package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func acmeiiwin(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if f == nil {
		f, err = os.Stat(path)
		if err != nil {
			return err
		}
	}
	if !f.Mode().IsDir() {
		return nil
	}
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	// Make sure this is the right kind of directory
	hasIn, hasOut := false, false
	names, err := d.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, name := range names {
		if name == "in" {
			hasIn = true
		} else if name == "out" {
			hasOut = true
		}
	}
	if !(hasIn && hasOut) {
		return nil
	}
	args := fmt.Sprintf("label \"%s\"; exec acmeiiwin \"%s\"", path, path)
	cmd := exec.Command("win", "sh", "-c", args)
	return cmd.Start()
}

func watchDir(dir string) {
	inFile := fmt.Sprintf("%s/in", dir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Name == inFile {
					// "in" is recreated often (by ii?)
					continue
				}
				acmeiiwin(event.Name, nil, nil)
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var dir string

	switch len(os.Args) {
	case 1:
		dir = fmt.Sprintf("%s/irc", os.Getenv("HOME"))
	case 2:
		dir = os.Args[1]
	default:
		fmt.Fprintln(os.Stderr, "usage: acmeii [dir]")
		os.Exit(1)
	}

	if err := filepath.Walk(dir, acmeiiwin); err != nil {
		log.Fatal(err)
	}
	watchDir(dir)
}
