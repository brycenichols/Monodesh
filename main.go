package main

import (
	"io"
	"os"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func sanitize(s string) string {
	return strings.TrimSpace(s)
}

func transFormVal(vals map[string]interface{}) (interface{}, bool) {
	count := 0
	var out interface{} = nil
	var valid bool = false
	for k, val := range vals {
		// Only take the first recognized field
		if count == 1 {
			break
		}
		count = 1
		switch sanitize(k) {
		case "S":
			v := sanitize(val.(string))
			// TODO: RFC3339
			t, err := time.Parse(time.RFC3339, v)
			if err == nil {
				out = t.Unix()
				valid = true
			} else {
				if v != "" {
					out = v
					valid = true
				}
			}
		case "N":
			v := sanitize(val.(string))
			f, err := strconv.ParseFloat(v, 64)
			if err == nil {
				out = f
				valid = true
			} else {

				i, err := strconv.ParseInt(v, 64, 10)
				if err == nil {
					out = i
					valid = true
				}
			}
		case "BOOL":
			v := sanitize(val.(string))
			valid = true
			switch v {
			case "1":
				out = true
			case "t":
				out = true
			case "T":
				out = true
			case "TRUE":
				out = true
			case "true":
				out = true
			case "True":
				out = true
			case "0":
				out = false
			case "f":
				out = false
			case "F":
				out = false
			case "FALSE":
				out = false
			case "false":
				out = false
			case "False":
				out = false
			default:
				valid = false
				out = nil
			}
		case "NULL":
			v := sanitize(val.(string))
			switch v {
			case "1":
				out = nil
			case "t":
				out = nil
				valid = true
			case "T":
				out = nil
				valid = true
			case "TRUE":
				out = nil
				valid = true
			case "true":
				out = nil
				valid = true
			case "True":
				out = nil
				valid = true
			default:
				valid = false
			}
		case "L":
			var l []interface{}
			switch val.(type) {
			case []interface{}:
				inp_list := val.([]interface{})
				for _, iv := range inp_list {
					switch iv.(type) {
					case map[string]any:
						res, valid := transFormVal(iv.(map[string]any))
						if valid {
							l = append(l, res)
						}
					}
				}
				if len(l) > 0 {
					out = l
					valid = true
				} else {
					out = nil
					valid = false
				}
			default:
				out = nil
				valid = false
			}
		case "M":
			var m = make(map[string]interface{})
			for k, mv := range val.(map[string]any) {
				sk := sanitize(k)
				if sk != "" {
					res, valid := transFormVal(mv.(map[string]any))
					if valid {
						m[sk] = res
					}
				}
			}
			if len(m) > 0 {
				out = m
				valid = true
			}
		default:
			count = 0
		}
	}
	return out, valid
}

func main() {
	stdin, err := io.ReadAll(os.Stdin)

	var result map[string]any

	if err == nil {
		json.Unmarshal([]byte(stdin), &result)
		var a []map[string]interface{}
		var m = make(map[string]interface{})
		for k, mv := range result {
			sk := sanitize(k)
			if sk != "" {
				res, valid := transFormVal(mv.(map[string]interface{}))
				if valid {
					m[sk] = res
				}
			}
		}
		a = append(a, m)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(a)
	}
}
