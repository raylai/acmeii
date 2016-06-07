package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: acmeiiwin dir\n", os.Args[0])
		os.Exit(1)
	}
	dir := os.Args[1]
	cchannel := tailFile(fmt.Sprintf("%s/out", dir))
	cuser := readUser()
	infile, err := os.OpenFile(fmt.Sprintf("%s/in", dir),
		os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case msg := <-cchannel:
			fmt.Println(msg)
		case msg := <-cuser:
			fmt.Fprintln(infile, msg)
		}
	}
}
