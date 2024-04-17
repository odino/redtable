package util

import (
	"time"

	"cloud.google.com/go/bigtable"
)

func ReadBTValue(r bigtable.Row) (string, bool) {
	v, ok := r["_values"]

	if !ok {
		return "", false
	}

	var hasValue bool
	var value string
	var isExpired bool

	for _, c := range v {
		if c.Column == "_values:value" {
			value = string(c.Value)
			hasValue = true
		}

		if c.Column == "_values:exp" {
			ts, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(c.Value))

			if err != nil {
				continue
			}

			if time.Until(ts) <= 0 {
				isExpired = true
			}
		}
	}

	if !hasValue || isExpired {
		return "", false
	}

	return value, true
}
