package command

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigtable"
)

type FlushAll struct{}

func (cmd *FlushAll) Parse(args []string) error {
	return nil
}

func (cmd *FlushAll) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	keys := []string{}
	muts := []*bigtable.Mutation{}

	err := tbl.ReadRows(context.Background(), bigtable.InfiniteRange(""), func(r bigtable.Row) bool {
		keys = append(keys, r.Key())
		mut := bigtable.NewMutation()
		mut.DeleteRow()
		muts = append(muts, mut)

		return true
	})

	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "OK", nil
	}

	_, err = tbl.ApplyBulk(ctx, keys, muts)

	return fmt.Sprintf("deleted %d keys", len(keys)), err
}
