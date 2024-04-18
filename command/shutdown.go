package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type Shutdown struct{}

func (cmd *Shutdown) Parse(args []resp.Arg) error {
	opts := map[string]bool{"SAVE": true, "NOSAVE": true, "NOW": true, "FORCE": true}

	for _, arg := range args {
		if _, ok := opts[arg.String()]; ok {
			continue
		}

		if arg.IsOption("ABORT") {
			return resp.ErrUnsupportedInRedtable("SHUTDOWN with ABORT")
		}

		return resp.ErrSyntax()
	}
	return nil
}

func (cmd *Shutdown) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	return nil, resp.ErrShutdown
}
