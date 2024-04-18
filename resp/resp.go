package resp

import (
	"errors"
	"fmt"
	"strings"

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

func SimpleInt(i int) redcon.SimpleInt {
	return redcon.SimpleInt(i)
}

var OK = SimpleString("OK")

type Arg string

func (a Arg) IsOption(s string) bool {
	return strings.EqualFold(s, a.String())
}

func (a Arg) String() string {
	return string(a)
}
