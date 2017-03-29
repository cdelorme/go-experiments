package main

// demonstration of using the reflect package to pre-emptively
// cast map[string]string into expected types for json marshal
// to populate a final structure with correct data

// for reference:
// @link: http://blog.golang.org/laws-of-reflection
// @link: https://golang.org/pkg/reflect/#StructTag.Get

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func main() {}

func to(o interface{}, m map[string]interface{}) {
	jd, _ := json.Marshal(m)
	json.Unmarshal(jd, o)
}

func isNumeric(t reflect.Kind) bool {
	switch t {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func reCast(o interface{}, m map[string]interface{}) {
	c := reflect.ValueOf(o).Elem()
	var skipstep bool
	for i := 0; i < c.NumField(); i++ {
		t := strings.Split(c.Type().Field(i).Tag.Get("json"), ",")[0]
		if t == "-" {
			continue
		}
		n := c.Type().Field(i).Name
		at := c.Field(i).Kind()
		for k, v := range m {
			if (t == "" && n == k) || t == k {
				skipstep = true
				in := reflect.TypeOf(v).Kind()
				if at == in || (isNumeric(at) && isNumeric(in)) {
					continue
				}
				if in == reflect.String && at == reflect.Bool {
					m[k], _ = strconv.ParseBool(v.(string))
				} else if in == reflect.Bool && at == reflect.String {
					m[k] = strconv.FormatBool(v.(bool))
				} else if isNumeric(in) && at == reflect.Bool {
					m[k] = (v == 1)
				} else if at == reflect.String && isNumeric(in) {
					m[k] = strconv.FormatFloat(reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0))).Float(), 'G', -1, 64)
				} else if in == reflect.Bool && isNumeric(at) && v.(bool) {
					m[k] = 1
				} else if in == reflect.String && isNumeric(at) {
					var err error
					if m[k], err = strconv.ParseFloat(v.(string), 64); err != nil {
						if b, _ := strconv.ParseBool(v.(string)); b {
							m[k] = 1
						}
					}
				} else if in == reflect.Map && at == reflect.Struct {
					if p, ok := v.(map[string]interface{}); ok {
						reCast(c.Field(i).Addr().Interface(), p)
					}
				}
			}
		}
		if !skipstep {
			if t == "" && at == reflect.Struct {
				reCast(reflect.New(c.Field(i).Type()).Interface(), m)
			}
		}
		skipstep = false
	}
}

func edge(o interface{}, m map[string]interface{}) {
	reCast(o, m)
	to(o, m)
}
