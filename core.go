package namespacer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	defaultFormat = "_%v"
	defaultIndex  = 2
)

// DefaultFormat is the value the library uses to initialize name spacers
func DefaultFormat() string {
	return defaultFormat
}

// DefaultIndex is the value the library uses to initialize name spacers
func DefaultIndex() uint {
	return defaultIndex
}

// Get the spaced version of a file path
func MustSpace(path string) string {
	path, err := SpacedName(path)
	if err != nil {
		panic(err)
	}
	return path
}

// Get the spaced version of a file path
func SpacedName(path string) (string, error) {
	return SpacedName2(path, defaultFormat)
}

// SpacedName2 gets the spaced version of a file path
// This version allows you to specify the format string that will differentiate collisions
func SpacedName2(path, format string) (string, error) {
	return SpacedName3(path, format, defaultIndex)
}

// SpacedName3 gets the spaced version of a file path
// This version allows you to specify the format string that will differentiate collisions
// as well as the starting index
func SpacedName3(path, format string, index uint) (string, error) {
	return SpacedName4(path, format, index, LocalFS{})
}

// SpacedName4 gets the spaced version of a file path on a custom file system
// we recommend using 2 as an index
func SpacedName4(path, format string, index uint, fileSystem fs.FS) (string, error) {
	ns := NameSpacer{
		Format: format,
		Index:  index,
		FS:     fileSystem,
	}
	return ns.Space(path)
}

type (
	NameSpacer struct {
		// Handles the spacing of path names by checking for existing matches
		Format string // Pattern that spaced names should follow. "_%v" is easier to parse but " (%v)" is common to major OSs.
		Index  uint   // Starting index for collision evasion. For most uses it is best to set this equal to 2.
		FS     fs.FS
	}
)

// Apply the spacing operation on the path
func (this *NameSpacer) Space(path string) (string, error) {
	_, err := fs.Stat(this.FS, path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		return path, nil
	}
	n := this.new(path)
	_, err = fs.Stat(this.FS, n)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		return n, nil
	} else {
		this.Index++
		return this.Space(path)
	}
}

// new is responsible for generating a file name without clobbering the extension
func (this NameSpacer) new(path string) string {
	ext := filepath.Ext(path)
	name := path[:len(path)-len(ext)]
	return (name + fmt.Sprintf(this.Format, this.Index) + ext)
}

type LocalFS struct{} // LocalFS is a local file system, just a wrapper on os.Open for now

var _ fs.FS = LocalFS{}

// Open attempts to retrieve a file from the disk
func (LocalFS) Open(name string) (fs.File, error) { return os.Open(name) }
