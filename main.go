// Command md-table-fmt reads markdown from standard input and writes it to
// standard output with the columns of each table vertically aligned. All
// other content, including table-like content inside fenced code blocks,
// passes through unchanged.
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "md-table-fmt:", err)
		os.Exit(1)
	}
}

func run(stdin io.Reader, stdout io.Writer) error {
	content, err := io.ReadAll(stdin)
	if err != nil {
		return err
	}
	_, err = io.WriteString(stdout, Format(string(content)))
	return err
}
