package command

import (
	"context"

	"cloud.google.com/go/bigtable"
)

type Ping struct{}

func (cmd *Ping) Parse(args []string) error {
	return nil
}

func (cmd *Ping) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	return "PONG", nil
}
