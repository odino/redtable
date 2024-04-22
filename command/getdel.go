package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type GetDel struct {
	Key string
}

func (cmd *GetDel) Parse(args []resp.Arg) error {
	if len(args) != 1 {
		return resp.ErrNumArgs("getdel")
	}

	cmd.Key = args[0].String()

	return nil
}

func (cmd *GetDel) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	row, ok, err := util.GetRow(ctx, cmd.Key, tbl)

	if !ok {
		return nil, err
	}

	_, err = util.DeleteRow(ctx, cmd.Key, tbl)

	return row.Value, err
}
