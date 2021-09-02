package main

import (
	"fmt"
	"golang.org/x/mod/sumdb/dirhash"
	"os"
)

func main() {
	os.Exit(realMain(os.Args))
}

func realMain(args []string) int {
	if len(args) != 2 {
		errorf("usage: %s archive.zip|directory\n", args[0])
		return 1
	}
	path := args[1]

	fi, err := os.Stat(path)
	if err != nil {
		errorln(err)
		return 1
	}

	var hash string
	if fi.Mode().IsDir() {
		hash, err = dirhash.HashDir(path, "", dirhash.Hash1)
	} else {
		hash, err = dirhash.HashZip(path, dirhash.Hash1)
	}

	if err != nil {
		errorln(err)
		return 1
	}

	fmt.Println(hash)
	return 0
}

func errorln(err error) bool {
	_, e := fmt.Fprintln(os.Stderr, err)
	return e == nil
}

func errorf(format string, a ...interface{}) bool {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	return err == nil
}
