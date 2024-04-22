package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type Echo struct {
	String string
}

func (cmd *Echo) Parse(args []resp.Arg) error {
	if len(args) != 1 {
		return resp.ErrNumArgs("echo")
	}

	cmd.String = args[0].String()

	return nil
}

func (cmd *Echo) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	return cmd.String, nil
}
