package command

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
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
	row, err := tbl.ReadRow(ctx, cmd.Key, bigtable.RowFilter(bigtable.LatestNFilter(1)))

	if err != nil {
		return "", err
	}

	v, ok := row["_values"]

	if !ok {
		return nil, nil
	}

	val := -1

	for _, c := range v {
		if c.Column == "_values:exp" {
			ts, err := strconv.Atoi(string(c.Value))

			if err != nil {
				return nil, resp.ErrBrokenKey
			}

			t := time.UnixMilli(int64(ts))

			if err != nil {
				break
			}

			if time.Until(t) >= 0 {
				val = int(time.Until(t).Round(time.Second).Seconds())
				break
			}
		}
	}

	return resp.SimpleInt(val), nil
}
