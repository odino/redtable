package command

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
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
		fmt.Println(1)
		return true
	},
		bigtable.RowFilter(
			bigtable.ChainFilters(
				bigtable.LatestNFilter(1),
				bigtable.TimestampRangeFilter(time.Now(), time.Time{}),
			),
		), bigtable.WithFullReadStats(func(frs *bigtable.FullReadStats) {
			i = int(frs.ReadIterationStats.RowsReturnedCount)
		}))

	return resp.SimpleInt(i), err
}
