package main

import (
	"flag"
	"fmt"
	"os"

	. "github.com/kendfss/namespacer"
)

var (
	format string
	index  uint
)

func init() {
	flag.StringVar(&format, "format", DefaultFormat(), "Format string to append to the path's name. Must contain %v, %d, %b (binary), or %x (hex). For more information visit https://pkg.go.dev/fmt.")
	flag.UintVar(&index, "index", DefaultIndex(), "First index used for incrementation")
}

func main() {
	flag.Parse()

	for _, path := range flag.Args() {
		path, err := SpacedName3(path, format, uint(index))
		if err != nil {
			internalError.abort(err)
		}
		fmt.Println(path)
	}
}

type errorCode int

const (
	noError errorCode = iota
	userInputError
	internalError
)

func (this errorCode) String() string {
	return []string{"NO ERROR", "USER INPUT ERROR", "INTERNAL ERROR"}[this]
}

func (code errorCode) abort(message error) {
	fmt.Fprint(os.Stderr, fmt.Errorf("[%s]: %w\n", code, message))
	flag.Usage()
	os.Exit(int(code))
}
