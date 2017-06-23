package logging

import (
	"io"
	"os"
	"strings"
)

type WriterFunc func(p []byte) (n int, err error)

func (w WriterFunc) Write(p []byte) (n int, err error) {
	return w(p)
}

func NewLogger(debug bool) io.Writer {
	return WriterFunc(func(p []byte) (n int, err error) {
		if !debug && strings.Contains(string(p), "[DEBUG]") {
			return 0, nil
		}

		return os.Stdout.Write(p)
	})
}

type writerFunc func(p []byte) (n int, err error)

func (w writerFunc) Write(p []byte) (n int, err error) {
	return w(p)
}
