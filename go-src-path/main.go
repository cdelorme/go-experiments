package main

// using `runtime.Caller` we can trace the src path
// useful when resources are included with the source

// while XDG pathing would be preferred this is very
// useful for including a `public/` with a project

// another option is game assets which probably
// should exist with the game source

// in general this would be best combined with other options
// such as a flag to set the path for resources, or a check
// against XDG, with a fallback to the src path

// when you copy an executable, the src path does not change
// from the time it compiled, rendering it useless for shared
// binaries/executables

import (
	"fmt"
	"runtime"
)

func main() {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		fmt.Println(file)
	}
}
