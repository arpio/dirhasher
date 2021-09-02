# dirhasher

### Command-line interface to Go's `dirhash` module

`dirhasher` computes a hash of a directory (or a ZIP archive) using go's
built-in `dirhash` module and prints it to stdout.

When computing the hash of a directory, the `prefix` argument passed
to `dirhash.HashDir` is always `""`.

### Usage

Use it on a ZIP archive:

    dirhasher my-archive.zip

Or on a directory:

    dirhasher /tmp/some-directory 

Package documentation is available at
[godoc](https://godoc.org/github.com/arpio/dirhasher).

### See Also

https://pkg.go.dev/golang.org/x/mod/sumdb/dirhash
