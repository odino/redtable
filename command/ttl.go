package command

import (
	"context"
	"math"
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
	row, err := util.GetRow(ctx, cmd.Key, tbl)

	if err != nil {
		return "", err
	}

	if !row.Found {
		return resp.SimpleInt(-2), nil
	}

	if row.Timestamp.Unix() == NO_EXPIRY_TS.Unix() {
		return resp.SimpleInt(-1), nil
	}

	i := int(math.RoundToEven(time.Until(row.Timestamp).Seconds()))
	return resp.SimpleInt(i), nil
}
