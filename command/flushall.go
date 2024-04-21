package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type FlushAll struct {
	Async bool
}

func (cmd *FlushAll) Parse(args []resp.Arg) error {
	if len(args) > 1 {
		return resp.ErrSyntax()
	}

	for _, arg := range args {
		if arg.IsOption("ASYNC") {
			cmd.Async = true
			continue
		}

		if arg.IsOption("SYNC") {
			continue
		}

		return resp.ErrSyntax()
	}

	return nil
}

func (cmd *FlushAll) doRun(ctx context.Context, tbl *bigtable.Table) (any, error) {
	keys := []string{}
	muts := []*bigtable.Mutation{}

	err := util.ScanTable(tbl, func(r bigtable.Row) bool {
		keys = append(keys, r.Key())
		mut := bigtable.NewMutation()
		mut.DeleteRow()
		muts = append(muts, mut)

		return true
	})

	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return resp.OK, nil
	}

	_, err = tbl.ApplyBulk(ctx, keys, muts)

	return resp.OK, err
}

func (cmd *FlushAll) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	if cmd.Async {
		go cmd.doRun(context.Background(), tbl)
		return resp.OK, nil
	}

	return cmd.doRun(ctx, tbl)
}
