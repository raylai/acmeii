package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func win(path string, f os.FileInfo, err error) error {
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
	return exec.Command("win", "acmeiiwin", path).Start()
}

func watchDir(dir string, depth int) {
	inFile := fmt.Sprintf("%s/in", dir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Name == inFile {
					// "in" is recreated often (by ii?)
					continue
				}
				win(event.Name, nil, nil)
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	if depth > 0 {
		f, err := os.Stat(dir)
		if err != nil {
			log.Fatal(err)
		}
		if f.Mode().IsDir() {
			watchDir(dir, depth-1)
		}
	}
	<-done
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

	if err := filepath.Walk(dir, win); err != nil {
		log.Fatal(err)
	}
	watchDir(dir, 1)
}
