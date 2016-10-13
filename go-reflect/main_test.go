package main

import (
	"testing"
)

type expected struct {
	NoMap          string
	NumberToNumber int    `json:"numberToNumber,omitempty"`
	BoolToBool     bool   `json:"boolToBool,omitempty"`
	StringToString string `json:"stringToString,omitempty"`
	ForcedPrivate  string `json:"-"`
}

type embedded struct {
	NonConflicting string
}

type embedcast struct {
	StringExpected string `json:"stringExpected,omitempty"`
}

type unexpected struct {
	expected
	embedcast
	Composite          embedded
	Unmapped           string
	StringToBool       bool   `json:"stringToBool,omitempty"`
	NumberToBool       bool   `json:"numberToBool,omitempty"`
	BoolToString       string `json:"boolToString,omitempty"`
	NumberToString     string `json:"numberToString,omitempty"`
	StringToNumber     int    `json:"stringToNumber,omitempty"`
	BoolToNumber       int    `json:"boolToNumber,omitempty"`
	StringBoolToNumber int    `json:"stringBoolToNumber,omitempty"`
}

var expectedMap = map[string]interface{}{
	"NoMap":          "yarp",
	"NumberToNumber": 42,
	"boolToBool":     true,
	"stringToString": "yarp",
	"StringToString": "narp", // this is ignored in favor of tag match
}

var unexpectedMap = map[string]interface{}{
	"NoMap":              "yarp",
	"Unmapped":           "yarp",
	"NumberToNumber":     42,
	"boolToBool":         true,
	"stringToString":     "yarp",
	"stringToBool":       "true",
	"numberToBool":       1,
	"boolToString":       true,
	"numberToString":     42,
	"stringToNumber":     "42",
	"boolToNumber":       true,
	"stringBoolToNumber": "true",
	"Composite":          map[string]interface{}{"NonConflicting": 42},
	"stringExpected":     42,
}

func TestTo(t *testing.T) {
	// t.Logf("A common data map: %#v\n", expectedMap)

	var o expected
	to(&o, expectedMap)
	// t.Logf("json marshal into expected structure: %#v\n", o)

	if o.BoolToBool != true {
		t.Logf("failed bool to bool conversion")
		t.Fail()
	}

	if o.NumberToNumber != 42 {
		t.Logf("failed number to number conversion")
		t.Fail()
	}

	if o.StringToString != "yarp" {
		t.Logf("failed string to string conversion")
		t.Fail()
	}

	if o.NoMap != "yarp" {
		t.Logf("failed to catch property without tag")
		t.Fail()
	}
}

func TestEdge(t *testing.T) {
	// t.Logf("A map containing edge cases: %#v\n", unexpectedMap)

	var o unexpected

	to(&o, unexpectedMap)
	// t.Logf("results of casting using old method: %#v\n", o)
	if o.StringToBool == true {
		t.Logf("string to bool was magically solved by json...")
		t.Fail()
	}
	if o.NumberToBool == true {
		t.Logf("number to bool was magically solved by json...")
		t.Fail()
	}
	if o.BoolToString == "true" {
		t.Logf("bool to string was magically solved by json...")
		t.Fail()
	}
	if o.NumberToString == "42" {
		t.Logf("number to string was magically solved by json...")
		t.Fail()
	}
	if o.StringToNumber == 42 {
		t.Logf("string to number was magically solved by json...")
		t.Fail()
	}
	if o.BoolToNumber == 1 {
		t.Logf("bool to number was magically solved by json...")
		t.Fail()
	}
	if o.StringBoolToNumber == 1 {
		t.Logf("string bool to number was magically solved by json...")
		t.Fail()
	}

	edge(&o, unexpectedMap)
	// t.Logf("results of casting using new method: %#v\n", o)
	if o.StringToBool != true {
		t.Logf("string to bool has not been solved for...")
		t.Fail()
	}
	if o.NumberToBool != true {
		t.Logf("number to bool has not been solved for...")
		t.Fail()
	}
	if o.BoolToString != "true" {
		t.Logf("bool to string has not been solved for...")
		t.Fail()
	}
	if o.NumberToString != "42" {
		t.Logf("number to string has not been solved for...")
		t.Fail()
	}
	if o.StringToNumber != 42 {
		t.Logf("string to number has not been solved for...")
		t.Fail()
	}
	if o.BoolToNumber != 1 {
		t.Logf("bool to number has not been solved for...")
		t.Fail()
	}
	if o.StringBoolToNumber != 1 {
		t.Logf("string bool to number has not been solved for...")
		t.Fail()
	}
	if o.Composite.NonConflicting != "42" {
		t.Logf("named composition does not work...")
		t.Fail()
	}
	if o.StringExpected != "42" {
		t.Logf("implicit composition with casting does not work...")
		t.Fail()
	}
}

func BenchmarkEdgeCases(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var o unexpected
		edge(&o, unexpectedMap)
	}
}
