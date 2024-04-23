package command

import (
	"context"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type Time struct{}

func (cmd *Time) Parse(args []resp.Arg) error {
	if len(args) != 0 {
		return resp.ErrNumArgs("time")
	}

	return nil
}

func (cmd *Time) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	now := time.Now()
	secs := now.Unix()
	sub := now.Sub(time.Unix(secs, 0)).Microseconds()

	return []any{secs, sub}, nil
}
