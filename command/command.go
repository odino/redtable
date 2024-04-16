package command

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigtable"
)

type Command interface {
	Parse([]string) error
	Run(context.Context, *bigtable.Table) (any, error)
}

func getCmd(s string, args []string) (Command, error) {
	var cmd Command
	var err error

	switch strings.ToLower(s) {
	case "ping":
		cmd = &Ping{}
	case "flushall":
		cmd = &FlushAll{}
	case "set":
		cmd = &Set{}
	case "get":
		cmd = &Get{}
	case "del":
		cmd = &Del{}
	case "ttl":
		cmd = &TTL{}
	default:
		err = fmt.Errorf("unknown command '%s', with args beginning with: %s", s, strings.Join(args, ","))
	}

	if err != nil {
		return nil, err
	}

	err = cmd.Parse(args)

	return cmd, err
}

func Process(cmd string, args []string, tbl *bigtable.Table) (any, error) {
	c, err := getCmd(cmd, args)

	if err != nil {
		return "", err
	}

	return c.Run(context.TODO(), tbl)
}
