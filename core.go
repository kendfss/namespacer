package namespacer

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	format string
	path   string
	index  int
)

// Get the spaced version of a file path
func MustSpace(path string) string {
	ns := NameSpacer{
		Format: "_%v",
		Index:  2,
	}
	out, err := ns.Space(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return out
}

// Get the spaced version of a file path
func SpacedName(path string) (string, error) {
	ns := NameSpacer{
		Format: "_%v",
		Index:  2,
	}
	return ns.Space(path)
}

// Get the spaced version of a file path
// This version allows you to specify the format string that will differentiate collisions
func SpacedName2(path, format string) (string, error) {
	ns := NameSpacer{
		Format: format,
		Index:  2,
	}
	return ns.Space(path)
}

// Get the spaced version of a file path
// This version allows you to specify the format string that will differentiate collisions
// as well as the starting index
func SpacedName3(path, format string, index uint) (string, error) {
	ns := NameSpacer{
		Format: format,
		Index:  index,
	}
	return ns.Space(path)
}

type (
	NameSpacer struct {
		// Handles the spacing of path names by checking for existing matches
		Format string // Pattern that spaced names should follow. "_%v" is easier to parse but " (%v)" is common to major OSs.
		Index  uint   // Starting index of duplicates. For most uses it is best to set this equal to 2
	}
)

// Apply the spacing operation on the path
func (this *NameSpacer) Space(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return path, nil
		}
		return "", err
	} else {
		n := this.new(path)
		_, err = os.Stat(n)
		if err != nil {
			if os.IsNotExist(err) {
				return n, nil
			}
			return "", err
		} else {
			this.Index++
			return this.Space(path)
		}
	}
}

func (this NameSpacer) new(path string) string {
	ext := filepath.Ext(path)
	name := path[:len(path)-len(ext)]
	return (name + fmt.Sprintf(this.Format, this.Index) + ext)
}

type (
	LocalFS struct{}
	HttpFS  struct {
		url.URL
		open func(name string) (fs.File, error)
	}
	httpFileInfo struct {
		name    func() string      // base name of the file
		size    func() int64       // length in bytes for regular files; system-dependent for others
		mode    func() fs.FileMode // file mode bits
		modtime func() time.Time   // modification time
		isdir   func() bool        // abbreviation for Mode().IsDir()
		sys     func() interface{} // underlying data source (can return nil)
	}
	httpFile struct {
		resp http.Response
		stat func() (fs.FileInfo, error) // unimplemented
	}
)

func (LocalFS) Open(name string) (fs.File, error) { return os.Open(name) }

func (self HttpFS) Open(name string) (fs.File, error) { return self.open(name) }
func (self httpFile) Stat() (fs.FileInfo, error)      { return self.stat() }
func (self httpFile) Read(buf []byte) (int, error)    { return io.ReadFull(self.resp.Body, buf) }
func (self httpFile) Close() error                    { return self.resp.Body.Close() }

func (self *httpFileInfo) of(src httpFile) *httpFileInfo { return self }

// base name of the file
func (self httpFileInfo) Name() string { return self.name() }

// length in bytes for regular files; system-dependent for others
func (self httpFileInfo) Size() int64 { return self.size() }

// file mode bits
func (self httpFileInfo) Mode() fs.FileMode { return self.mode() }

// modification time
func (self httpFileInfo) ModTime() time.Time { return self.modtime() }

// abbreviation for Mode().IsDir()
func (self httpFileInfo) Isdir() bool { return self.isdir() }

// underlying data source (can return nil)
func (self httpFileInfo) Sys() interface{} { return self.sys() }
