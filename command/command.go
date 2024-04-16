package command

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigtable"
)

type Command interface {
	Run(context.Context, []string, *bigtable.Table) (any, error)
}

type Commander func() Command

var commands map[string]Commander

func init() {
	commands = map[string]Commander{
		"ping":     func() Command { return &Ping{} },
		"flushall": func() Command { return &FlushAll{} },
		"set":      func() Command { return &Set{} },
		"del":      func() Command { return &Del{} },
		"get":      func() Command { return &Get{} },
		"ttl":      func() Command { return &TTL{} },
	}
}

func Process(cmd string, args []string, tbl *bigtable.Table) (any, error) {
	c, ok := commands[cmd]

	if !ok {
		return "", fmt.Errorf("unknown command '%s'", cmd)
	}

	return c().Run(context.TODO(), args, tbl)
}
