package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
	os.Exit(0)

}

func nextToken(pattern *string) string {
	switch (*pattern)[0] {
	case '\\':
		return (*pattern)[0:2]
	case '^':
		return (*pattern)[0:1]
	case '[':
		k := 1

		for (*pattern)[k] != ']' {
			k += 1
		}

		return (*pattern)[0 : k+1]
	default:
		return (*pattern)[0:1]
	}

}

func match(line []byte, pattern string) bool {
	if len(pattern) == 0 {
		return true
	}
	if len(line) == 0 {
		return false
	}

	token := nextToken(&pattern)

	var matched bool

	switch token {
	case "\\d":
		if bytes.ContainsAny(line[0:1], "0123456789") {
			matched = true
		}
	case "\\w":
		if bytes.ContainsAny(line[0:1], "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_") {
			matched = true
		}
	default:
		if token[0] == '[' {
			if token[1] == '^' {
				if !bytes.ContainsAny(line[0:1], token[2:len(token)-1]) {
					matched = true
				}

			} else {
				if bytes.ContainsAny(line[0:1], token[1:len(token)-1]) {
					matched = true
				}
			}

		} else {
			matched = token[0] == line[0]
		}
	}

	if matched {
		return match(line[1:], pattern[len(token):])
	}

	return false
}

func matchLine(line []byte, pattern string) (bool, error) {

	if nextToken(&pattern) == "^" {
		return match(line, pattern[1:]), nil
	}

	for i := 0; i < len(line); i++ {
		if match(line[i:], pattern) {
			return true, nil
		}
	}

	return false, nil
}
