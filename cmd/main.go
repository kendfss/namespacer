package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	. "github.com/kendfss/namespacer"
)

const (
	noError errorCode = iota
	userInputError
	internalError
)

var (
	format string
	path   string
	index  int
)

type (
	errorCode int
)

func (this errorCode) Error() string {
	return this.String()
}

func (this errorCode) String() string {
	return []string{"NO ERROR", "USER INPUT ERROR", "INTERNAL ERROR"}[this]
}

func (code errorCode) abort(message error) {
	fmt.Fprint(os.Stderr, fmt.Errorf("[%s]: %w\n", code, message))
	flag.Usage()
	os.Exit(int(code))
}

func main() {
	// flag.StringVar(&path, "path", "", "Path to the file/folder you would like to space")
	flag.StringVar(&format, "format", "_%v", "Format string to append to the path's name. Must contain %v, %d, %b (binary), or %x (hex)")
	flag.IntVar(&index, "index", 2, "First index to attempt incrementing to")
	flag.Parse()

	if len(path) == 0 {
		if flag.Arg(0) == "" {
			userInputError.abort(errors.New("no args given"))
		}
		path = flag.Arg(0)
	}

	// for args, ctr := flag.Args(), 1; len(args)-1 > 0 && ctr < 3; ctr++ {
	// 	args = args[1:]
	// 	arg := args[0]
	// 	err := fmtOrInt(arg)
	// 	if err != nil {
	// 		userInputError.abort(fmt.Errorf("couldn't parse argument #%v (%v) because %w\n", ctr, arg, err))
	// 	}
	// }
	for _, path := range flag.Args() {
		path, err := SpacedName3(path, format, uint(index))
		if err != nil {
			internalError.abort(err)
		}
		fmt.Println(path)
	}
}

// check if an Arg is a format string or integer and assign it to format or index, respectively
func fmtOrInt(arg string) error {
	if arg != "" {
		if b, err := isNumeric(arg); b {
			if err != nil {
				return err
			}
			i, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			index = i
		} else if b, err := hasIntableFormat(arg); b {
			if err != nil {
				return err
			}
			format = arg
		} else {
			return errors.New("Argument is not an integer and cannot be used to format one")
		}
		return nil
	}
	return errors.New("Cannot parse empty argument")
}

// check if a string is purely digital
func isNumeric(str string) (bool, error) {
	return regexp.MatchString("^[0-9]+$", str)
}

// check if a string has an integer-friendly format parameter
func hasIntableFormat(str string) (bool, error) {
	return regexp.MatchString("(%[bdvx]){1}", str)
}
