package command

import (
	"context"
	"strconv"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type IncrBy struct {
	Incr       Incr
	Multiplier int
}

func (cmd *IncrBy) Parse(args []resp.Arg) error {
	if len(args) != 2 {
		return resp.ErrNumArgs(cmd.Incr.Name)
	}

	cmd.Incr.Key = args[0].String()

	i, err := strconv.Atoi(args[1].String())

	if err != nil {
		return resp.ErrInt
	}

	cmd.Incr.Delta = i * cmd.Multiplier

	return nil
}

func (cmd *IncrBy) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	return cmd.Incr.Run(ctx, tbl)
}
