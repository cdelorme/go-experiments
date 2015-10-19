package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// if you need to trace the full path to the executable
// including resolving symbolic links, such as for asset files
// you can use the EvalSymlinks method from the filepath package

func main() {

	// get abs path to self
	path, e := filepath.EvalSymlinks(os.Args[0])
	if e != nil {
		fmt.Println("Fail Eval Symlink!")
		return
	}

	appAbs, e := filepath.Abs(path)
	if e != nil {
		fmt.Println("Fail Abs!")
		return
	}

	fmt.Println(appAbs)

}
