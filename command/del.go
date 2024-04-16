package command

import (
	"context"
	"errors"

	"cloud.google.com/go/bigtable"
)

type Del struct {
	Key string
}

func (cmd *Del) Run(ctx context.Context, args []string, tbl *bigtable.Table) (any, error) {
	if len(args) < 1 {
		return "", errors.New("SET requires at least a key and a value")
	}

	mut := bigtable.NewMutation()
	mut.DeleteRow()
	err := tbl.Apply(ctx, args[0], mut)

	return 1, err
}
