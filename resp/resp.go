package resp

import (
	"errors"
	"fmt"

	"github.com/tidwall/redcon"
)

func ErrSyntax() error {
	return errors.New("syntax error")
}

func ErrUnsupportedInRedtable(s string) error {
	return fmt.Errorf("unsupported in redtable: %s -- open an issue at https://github.com/odino/redtable/issues", s)
}

func SimpleString(s string) redcon.SimpleString {
	return redcon.SimpleString(s)
}

var OK = SimpleString("OK")
