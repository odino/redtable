package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Keys struct{}

func (cmd *Keys) Parse(args []resp.Arg) error {
	if len(args) != 1 {
		return resp.ErrNumArgs("keys")
	}

	if args[0].String() != "*" {
		return resp.ErrUnsupportedInRedtable("KEYS with a pattern other than *")
	}

	return nil
}

func (cmd *Keys) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	keys := []string{}

	err := util.ScanTable(tbl, func(r bigtable.Row) bool {
		keys = append(keys, r.Key())
		return true
	})

	if err != nil {
		return nil, err
	}

	return keys, err
}
