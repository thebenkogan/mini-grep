package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Usage: echo <input_text> | your_grep.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Println("usage: mygrep -E <pattern>")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	ok, err := executePattern(os.Stdin, pattern)
	if err != nil {
		fmt.Println("failed to execute pattern on input", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func executePattern(reader io.Reader, pattern string) (bool, error) {
	line, err := io.ReadAll(reader) // assume we're only dealing with a single line
	if err != nil {
		return false, err
	}

	if len(pattern) != 1 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	return bytes.ContainsAny(line, pattern), nil
}
