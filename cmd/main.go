package main

import (
	"fmt"
	"io"
	"os"

	"github.com/thebenkogan/grep/internal/regex"
)

// Usage: echo <input_text> | ./build/main -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Println("usage: mygrep -E <pattern>")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	ok, err := executePattern(os.Stdin, pattern)
	if err != nil {
		fmt.Println("failed to execute pattern on input:", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func executePattern(reader io.Reader, pattern string) (bool, error) {
	regex, err := regex.NewRegex(pattern)
	if err != nil {
		return false, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return false, err
	}
	line := string(data)

	return regex.Match(line), nil
}
