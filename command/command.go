package command

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type Command interface {
	Parse([]resp.Arg) error
	Run(context.Context, *bigtable.Table) (any, error)
}

func getCmd(s string, args []resp.Arg) (Command, error) {
	var cmd Command
	var err error

	switch strings.ToLower(s) {
	case "ping":
		cmd = &Ping{}
	case "flushall", "flushdb":
		cmd = &FlushAll{}
	case "set":
		cmd = &Set{}
	case "get":
		cmd = &Get{}
	case "del":
		cmd = &Del{}
	case "ttl":
		cmd = &TTL{}
	case "append":
		cmd = &Append{}
	case "shutdown":
		cmd = &Shutdown{}
	case "dbsize":
		cmd = &DbSize{}
	case "bitcount":
		cmd = &BitCount{}
	case "getdel":
		cmd = &GetDel{}
	case "copy":
		cmd = &Copy{}
	case "echo":
		cmd = &Echo{}
	case "time":
		cmd = &Time{}
	case "rename":
		cmd = &Rename{}
	case "incr":
		cmd = &Incr{Name: "incr", Delta: 1}
	case "decr":
		cmd = &Incr{Name: "decr", Delta: -1}
	case "incrby":
		cmd = &IncrBy{Incr: Incr{Name: "incrby"}, Multiplier: 1}
	case "decrby":
		cmd = &IncrBy{Incr: Incr{Name: "decrby"}, Multiplier: -1}
	case "exists":
		cmd = &Exists{}
	case "keys":
		cmd = &Keys{}
	default:
		fmtargs := []string{}

		for _, s := range args {
			fmtargs = append(fmtargs, fmt.Sprintf("'%s' ", s))
		}

		err = fmt.Errorf("unknown command '%s', with args beginning with: %s", s, strings.Join(fmtargs, ""))
	}

	if err != nil {
		return nil, err
	}

	err = cmd.Parse(args)

	return cmd, err
}

func Process(cmd string, args []resp.Arg, tbl *bigtable.Table) (any, error) {
	c, err := getCmd(cmd, args)

	if err != nil {
		return "", err
	}

	return c.Run(context.TODO(), tbl)
}
