# Gobackup

Gobackup takes in a directory or file name as an argument and recursively compresses all the contents to a tarball (`.tar.gz`). For now, the argument must be the absolute path to the file or directory. Support for more archive filetypes will be coming in the future.

## Installation and Usage

If you have a Go environment set up, just run `go get github.com/rugglcon/gobackup` and Gobackup will be placed in your $GOPATH.

Otherwise, you can either clone this repository and run `go build gobackup.go` and place the binary wherever you would like, or download one of the prebuilt binaries over on the [Releases](https://github.com/rugglcon/gobackup/releases) page.