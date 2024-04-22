package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Get struct {
	Key string
}

func (cmd *Get) Parse(args []resp.Arg) error {
	if len(args) < 1 {
		return resp.ErrNumArgs("get")
	}

	cmd.Key = args[0].String()

	return nil
}

func (cmd *Get) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	row, ok, err := util.GetRow(ctx, cmd.Key, tbl)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return row.Value, nil
}
