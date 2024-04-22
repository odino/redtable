package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Del struct {
	Keys []string
}

func (cmd *Del) Parse(args []resp.Arg) error {
	if len(args) < 1 {
		return resp.ErrSyntax()
	}

	for _, arg := range args {
		cmd.Keys = append(cmd.Keys, arg.String())
	}

	return nil
}

func (cmd *Del) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	_, l, err := util.GetRows(ctx, cmd.Keys, tbl)

	if err != nil {
		return nil, err
	}

	err = util.DeleteRows(ctx, cmd.Keys, tbl)

	if err != nil {
		return nil, err
	}

	return resp.SimpleInt(l), err
}
