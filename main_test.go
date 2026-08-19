package main

import (
	"errors"
	"strings"
	"testing"
)

func TestRunFormatsStdinToStdout(t *testing.T) {
	stdin := strings.NewReader("| a | b |\n|---|---|\n| 1 | 22 |\n")
	var stdout strings.Builder
	err := run(stdin, &stdout)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, stdout.String(), "| a | b  |\n|---|----|\n| 1 | 22 |\n")
}

func TestRunReportsReadError(t *testing.T) {
	var stdout strings.Builder
	err := run(failingReader{}, &stdout)
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected read error, got: %v", err)
	}
	assertEqual(t, stdout.String(), "")
}

func TestRunReportsWriteError(t *testing.T) {
	err := run(strings.NewReader("content"), failingWriter{})
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected write error, got: %v", err)
	}
}

var errBoom = errors.New("boom")

type failingReader struct{}

func (this failingReader) Read([]byte) (int, error) {
	return 0, errBoom
}

type failingWriter struct{}

func (this failingWriter) Write([]byte) (int, error) {
	return 0, errBoom
}
