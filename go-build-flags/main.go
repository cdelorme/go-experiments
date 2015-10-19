package main

import (
	"fmt"
)

// demonstrate by printing variables
// ldflags are only sent once at build time

func main() {
	fmt.Printf("%s: %-30s\n", "abc", abc)
	fmt.Printf("%s: %-30s\n", "softwareName", softwareName)
	fmt.Printf("%s: %-30s\n", "buildVersion", buildVersion)
	fmt.Printf("%s: %-30s\n", "buildDate", buildDate)
}
