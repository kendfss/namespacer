namespacer
---
  
CLI to avoid over-writing file-sytem names which already exist.  
It is insensitive to argument type (file/directory).  
  
### Installation
```shell
go install github.com/kendfss/namespacer
```
Or
```shell
git clone github.com/kendfss/namespacer
cd namespacer
go install
```
Or
```shell
git clone github.com/kendfss/namespacer
cd namespacer
make install
```

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
namespacer -path main.go -format " (%v)" -index 3 
# "main (3).go"
```
  
```shell
namespacer main.go " v%d" 40 
# main v40.go
```
  
```shell
namespacer main.go " v%x" 40 
# main v28.go
```
  
```shell
namespacer main.go " v%b" 40 
# main v101000.go
```
