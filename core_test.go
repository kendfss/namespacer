package namespacer

import (
	"math"
	"math/rand"
	"os"
	"testing"
	"unicode"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ._-"

func touch(name string) error {
	return os.WriteFile(name, []byte{}, os.ModePerm)
}

func runePred(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsPunct(r)
}

func randStr() string {
	return getLetters(takeN(rand.Intn(math.MaxUint8), upto(len(alphabet))...)...)
}

func upto(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func takeN(count int, args ...int) []int {
	out := make([]int, count)
	for i := range out {
		out[i] = args[rand.Intn(len(args))]
	}
	return out
}

func getLetters(args ...int) string {
	buf := make([]byte, len(args))
	for i, arg := range args {
		buf[i] = alphabet[arg]
	}
	return string(buf)
}

func TestSpace(t *testing.T) {
	type test struct {
		name   string
		levels int
	}

	tests := []test{
		{name: randStr(), levels: rand.Intn(0)},
	}
}
