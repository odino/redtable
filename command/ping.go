package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type Ping struct{}

func (cmd *Ping) Parse(args []resp.Arg) error {
	return nil
}

func (cmd *Ping) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	return resp.SimpleString("PONG"), nil
}
