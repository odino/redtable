package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
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
	mut := bigtable.NewMutation()
	mut.DeleteRow()
	err := tbl.Apply(ctx, cmd.Key, mut)

	return 1, err
}
