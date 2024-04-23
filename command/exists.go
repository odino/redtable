package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Exists struct {
	Keys []string
}

func (cmd *Exists) Parse(args []resp.Arg) error {
	if len(args) < 1 {
		return resp.ErrNumArgs("exists")
	}

	for _, arg := range args {
		cmd.Keys = append(cmd.Keys, arg.String())
	}

	return nil
}

func (cmd *Exists) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	_, l, err := util.GetRows(ctx, cmd.Keys, tbl)

	return resp.SimpleInt(l), err
}
