package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var Version string = "v0.0.0-dev"

func main() {
	exitIfErr := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if len(os.Args) != 2 {
		fmt.Printf("%s <filename>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	log := &log{name: os.Args[1], lineLimit: 1000000}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		_, err := fmt.Println(line)
		exitIfErr(err)
		_, err = log.Println(line)
		exitIfErr(err)
	}
	exitIfErr(scanner.Err())
}

type log struct {
	lock      sync.Mutex
	name      string
	lineLimit int
	file      *os.File
	lineCount int
}

func (lg *log) Println(line string) (int, error) {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	if lg.file == nil {
		var err error
		lg.file, err = os.Create(lg.name)
		if err != nil {
			return 0, err
		}
	}
	if n, err := lg.file.Write([]byte(line)); err != nil {
		return n, err
	}
	if n, err := lg.file.Write([]byte("\n")); err != nil {
		return n, err
	}
	lg.lineCount++
	if lg.lineCount >= lg.lineLimit {
		if err := lg.file.Close(); err != nil {
			return len([]byte(line)) + 1, err
		}
		if err := os.Rename(lg.name, lg.name+time.Now().Format(".2006-01-02.150405")); err != nil {
			return len([]byte(line)) + 1, err
		}
		var err error
		lg.file, err = os.Create(lg.name)
		if err != nil {
			return len([]byte(line)) + 1, err
		}
	}
	return len([]byte(line)) + 1, nil
}
