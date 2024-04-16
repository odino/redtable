package command

import (
	"context"
	"errors"

	"cloud.google.com/go/bigtable"
)

type Del struct {
	Key string
}

func (cmd *Del) Parse(args []string) error {
	if len(args) < 1 {
		return errors.New("SET requires at least a key and a value")
	}

	cmd.Key = args[0]

	return nil
}

func (cmd *Del) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	mut := bigtable.NewMutation()
	mut.DeleteRow()
	err := tbl.Apply(ctx, cmd.Key, mut)

	return 1, err
}
