package main

// verify walking a path and printing only "normal" files
// was mostly to see what accounts for "abnormal"

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	cwd, _ := os.Getwd()
	filepath.Walk(cwd, walk)
}

func walk(path string, f os.FileInfo, err error) error {
	if f != nil && f.Mode().IsRegular() {
		fmt.Println(f.Name())
		fmt.Println(path)
	}
	return err
}
