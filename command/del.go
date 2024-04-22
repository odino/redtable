package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Del struct {
	Key string
}

func (cmd *Del) Parse(args []resp.Arg) error {
	if len(args) < 1 {
		return resp.ErrSyntax()
	}

	cmd.Key = args[0].String()

	return nil
}

func (cmd *Del) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	ok, err := util.DeleteRow(ctx, cmd.Key, tbl)

	if err != nil {
		return nil, err
	}

	var res int

	if ok {
		res = 1
	}

	return resp.SimpleInt(res), err
}
