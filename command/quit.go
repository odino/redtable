package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type Quit struct{}

func (cmd *Quit) Parse(args []resp.Arg) error {
	return nil
}

func (cmd *Quit) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	return nil, resp.ErrQuit
}
