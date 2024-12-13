package namespacer

import (
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ._-"

// touch creates a file on the local system
func touch(name string) error {
	return os.WriteFile(name, []byte(name), os.ModePerm|0777)
}

// randStr generates a random string using the alphabet above and the functions below
func randStr(n int) string {
	return getLetters(takeN(n, upto(len(alphabet))...)...)
}

// upto returns a slice of indices in [0, n)
func upto(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

// takeN randomly selects a given number of the variadic terms
func takeN(amount int, indices ...int) []int {
	out := make([]int, amount)
	for i := range out {
		out[i] = indices[rand.Intn(len(indices))]
	}
	return out
}

// get letters fetches the letters at the given indices from the alphabet string
func getLetters(indices ...int) string {
	buf := make([]byte, len(indices))
	for i, arg := range indices {
		buf[i] = alphabet[arg]
	}
	return string(buf)
}

// randBool returns a random boolean value
func randBool() bool {
	return []bool{false, true}[rand.Intn(2)]
}

func TestSpaceLocal(t *testing.T) {
	tests := []struct {
		name   string
		levels int
		dir    bool
	}{
		{name: randStr(10), levels: rand.Intn(10) + 1, dir: randBool()},
		{name: randStr(10), levels: rand.Intn(10) + 1, dir: randBool()},
		{name: randStr(10), levels: rand.Intn(10) + 1, dir: randBool()},
		{name: randStr(10), levels: rand.Intn(10) + 1, dir: randBool()},
	}

	dir := t.TempDir()
	for _, test := range tests {
		name := filepath.Join(dir, test.name)
		t.Run(name, func(t *testing.T) {
			var err error
			if test.dir {
				err = os.Mkdir(name, 0777&fs.ModePerm)
			} else {
				err = touch(name)
			}
			if err != nil {
				t.Fatal(err)
			}

			for range make([]int, test.levels) {
				next := MustSpace(name)
				if name == next {
					t.Fatalf("new name (%q) == old name (%q)", next, name)
				}
				err := touch(next)
				if err != nil {
					t.Fatal(err)
				}
				name = next
			}
		})
	}
}
