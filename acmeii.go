package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/hpcloud/tail"
)

func readUser() <-chan string {
	c := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			c <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}
	}()
	return c
}

func tailFile(file string) <-chan string {
	c := make(chan string)
	go func() {
		t, err := tail.TailFile(file, tail.Config{Follow: true})
		if err != nil {
			log.Fatalln(err)
		}
		for line := range t.Lines {
			c <- line.Text
		}
		err = t.Wait()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	return c
}

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
