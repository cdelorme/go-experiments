package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// demonstrate custom MarshalJSON on a self-defined data type that extends
// another data type, augmenting what gets displayed

// in this case, when printed, the json.Marshal process fills in the Time value
// using the difference with the global started value, without having to directly
// add that computation

// there are more complex or "obvious" benefits to this on entire custom defined objects
// however there are many gotchas and this is somewhat indirect or abstracted making it
// possibly difficult to troubleshoot when a problem occurs

var started = time.Now()

type uptime time.Duration

func (u uptime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Now().Sub(started).String() + `"`), nil
}

type CustomObject struct {
	Data string `json:"data"`
	Time uptime `json:"time"`
}

func main() {

	n := CustomObject{}
	n.Data = "Some test data"

	fmt.Printf("Go Object: %+v\n", n)
	time.Sleep(time.Second * 3)
	j, _ := json.Marshal(n)
	fmt.Printf("Json: %s\n", j)
}
