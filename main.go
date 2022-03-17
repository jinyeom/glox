package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

func usage() {
	fmt.Println("Usage: lox [script]")
	os.Exit(64)
}

func handleErr(line int, msg string) {
	reportErr(line, "", msg)
}

func reportErr(line int, where, msg string) {
	fmt.Printf("[line %d] Error%s: %s", line, where, msg)
}

func run(source string) bool {
	var s scanner.Scanner
	s.Init(strings.NewReader(source))
	for t := s.Scan(); t != scanner.EOF; t = s.Scan() {
		fmt.Println(s.TokenText())
	}
	return true
}

func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err) // FIXME: handle this error better
	}
	if ok := run(string(data)); !ok {
		os.Exit(65)
	}
}

func runPrompt() {
	fmt.Println("Lox (gLox) 0.0.1")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
		}
		line = strings.Trim(line, "\n")
		run(line)
	}
}

func main() {
	args := os.Args
	switch {
	case len(args) > 2:
		usage()
	case len(args) == 2:
		runFile(args[1])
	default:
		runPrompt()
	}
}
