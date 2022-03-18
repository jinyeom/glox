package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	fmt.Println("Usage: lox [script]")
	os.Exit(64)
}

func handleErr(line int, msg string) {
	reportErr(line, "", msg)
}

func reportErr(line int, where, msg string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, msg)
}

func run(src string) bool {
	scanner := NewScanner(src)
	tokens, err := scanner.Scan()
	if err != nil {
		handleErr(scanner.line, err.Error())
		return false
	}
	for _, token := range tokens {
		fmt.Println(token)
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
