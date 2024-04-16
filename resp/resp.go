package resp

import (
	"errors"

	"github.com/tidwall/redcon"
)

func ErrSyntax() error {
	return errors.New("syntax error")
}

func SimpleString(s string) redcon.SimpleString {
	return redcon.SimpleString(s)
}

var OK = SimpleString("OK")
