package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type DbSize struct {
}

func (cmd *DbSize) Parse(args []resp.Arg) error {
	if len(args) != 0 {
		return resp.ErrNumArgs("dbsize")
	}

	return nil
}

func (cmd *DbSize) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	var i int

	err := tbl.ReadRows(context.Background(), bigtable.InfiniteRange(""), func(r bigtable.Row) bool {
		if _, ok := util.ParseRow(r); ok {
			i++
		}

		return true
	})

	return resp.SimpleInt(i), err
}
