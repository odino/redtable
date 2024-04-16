package command

import (
	"context"

	"cloud.google.com/go/bigtable"
)

type Ping struct{}

func (cmd *Ping) Run(ctx context.Context, args []string, tbl *bigtable.Table) (any, error) {
	return "PONG", nil
}
