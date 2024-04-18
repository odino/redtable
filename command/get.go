package command

import (
	"context"
	"errors"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Get struct {
	Key string
}

func (cmd *Get) Parse(args []resp.Arg) error {
	if len(args) < 1 {
		return errors.New("GET requires at least a key")
	}

	cmd.Key = args[0].String()

	return nil
}

func (cmd *Get) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	row, err := tbl.ReadRow(ctx, cmd.Key, bigtable.RowFilter(bigtable.LatestNFilter(1)))

	if err != nil {
		return nil, err
	}

	val, ok := util.ReadBTValue(row)

	if !ok {
		return nil, nil
	}

	return val, nil
}
