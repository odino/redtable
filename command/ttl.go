package command

import (
	"context"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type TTL struct {
	Key string
}

func (cmd *TTL) Parse(args []resp.Arg) error {
	if len(args) < 1 {
		return resp.ErrNumArgs("ttl")
	}

	cmd.Key = args[0].String()

	return nil
}

func (cmd *TTL) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	row, ok, err := util.GetRow(ctx, cmd.Key, tbl)

	val := -1

	if err != nil {
		return nil, err
	}

	if !ok {
		val = -2
	}

	if row.Expiry.IsZero() {
		val = -1
	}

	if time.Until(row.Expiry) >= 0 {
		val = int(time.Until(row.Expiry).Round(time.Second).Seconds())
	}

	return resp.SimpleInt(val), nil
}
