package command

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type Get struct {
	Key string
}

func (cmd *Get) Parse(args []resp.Arg) error {
	if len(args) < 1 {
		return errors.New("GET requires at least a key")
	}

	cmd.Key = string(args[0])

	return nil
}

func (cmd *Get) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	row, err := tbl.ReadRow(ctx, cmd.Key, bigtable.RowFilter(bigtable.LatestNFilter(1)))

	if err != nil {
		return nil, err
	}

	v, ok := row["_values"]

	if !ok {
		return nil, nil
	}

	var hasValue bool
	var value string
	var isExpired bool

	for _, c := range v {
		if c.Column == "_values:value" {
			value = string(c.Value)
			hasValue = true
		}

		if c.Column == "_values:exp" {
			ts, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(c.Value))

			if err != nil {
				continue
			}

			if time.Until(ts) <= 0 {
				isExpired = true
			}
		}
	}

	if !hasValue || isExpired {
		return nil, nil
	}

	return value, nil
}
