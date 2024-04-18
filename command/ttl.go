package command

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type TTL struct {
	Key string
}

func (cmd *TTL) Parse(args []resp.Arg) error {
	if len(args) < 1 {
		errors.New("TTL requires at least a key")
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
			ts, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(c.Value))

			if err != nil {
				break
			}

			if time.Until(ts) >= 0 {
				val = int(time.Until(ts).Round(time.Second).Seconds())
				break
			}
		}
	}

	return resp.SimpleInt(val), nil
}
