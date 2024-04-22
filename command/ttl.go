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

	if err != nil {
		return nil, err
	}

	if !ok {
		return resp.SimpleInt(-2), nil
	}

	if time.Until(row.Expiry) >= 0 {
		return resp.SimpleInt(int(time.Until(row.Expiry).Round(time.Second).Seconds())), nil
	}

	return resp.SimpleInt(-1), nil
}
