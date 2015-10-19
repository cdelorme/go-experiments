package imports

import (

	// a lot of examples show this on stack overflow
	"./bigger"

	// the go blog explicitly says **not to do this**
	// @link: http://blog.golang.org/organizing-go-code

	// neither of these work
	// "bigger"
	// "go-imports/bigger"

	"fmt"
)

func main() {
	fmt.Println(bigger.Bigger)
}
