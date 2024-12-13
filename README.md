namespacer
---
  
CLI to avoid over-writing file-sytem names which already exist.  
It is insensitive to argument type (file/directory).  
  
### Installation
```shell
git clone https://github.com/kendfss/namespacer
cd namespacer
make install
```
- You need a [Go](https://go.dev) compiler for this

### Usage
It accepts three arguments:  
- `path`: to the file/directory  
- `format`: string that specifies the pattern that spaced names should extend with  
- `index`: the first number to try if the path argument already exists  

```shell
namespacer main.go 
# returns "main_2.go" if main.go already exists
```
  
```shell
namespacer -format " (%v)" -index 3 main.go
# "main (3).go"
```
  
```shell
namespacer -format " v%d" -index 40 main.go
# main v40.go
```
  
```shell
namespacer -format " v%x" -index 40 main.go
# main v28.go
```
  
```shell
namespacer -format " v%b" -index 40 main.go
# main v101000.go
```
  


### API
```go
package namespacer // import "github.com/kendfss/namespacer"


// FUNCTIONS

func DefaultFormat() string
    // DefaultFormat is the value the library uses to initialize name spacers

func DefaultIndex() uint
    // DefaultIndex is the value the library uses to initialize name spacers

func MustSpace(path string) string
    // Get the spaced version of a file path

func SpacedName(path string) (string, error)
    // Get the spaced version of a file path

func SpacedName2(path, format string) (string, error)
    // SpacedName2 gets the spaced version of a file path This version allows you
    // to specify the format string that will differentiate collisions

func SpacedName3(path, format string, index uint) (string, error)
    // SpacedName3 gets the spaced version of a file path This version allows you
    // to specify the format string that will differentiate collisions as well as
    // the starting index

func SpacedName4(path, format string, index uint, fileSystem fs.FS) (string, error)
    // SpacedName4 gets the spaced version of a file path on a custom file system
    // we recommend using 2 as an index


// TYPES

type LocalFS struct{} // LocalFS is a local file system, just a wrapper on os.Open for now

func (LocalFS) Open(name string) (fs.File, error)
    // Open attempts to retrieve a file from the disk

type NameSpacer struct {
	// Handles the spacing of path names by checking for existing matches
	Format string // Pattern that spaced names should follow. "_%v" is easier to parse but " (%v)" is common to major OSs.
	Index  uint   // Starting index for collision evasion. For most uses it is best to set this equal to 2.
	FS     fs.FS
}

func (this *NameSpacer) Space(path string) (string, error)
    // Apply the spacing operation on the path

```
