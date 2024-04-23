package command

import (
	"context"
	"encoding/binary"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/redtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
	"google.golang.org/grpc/status"
)

type Incr struct {
	Key   string
	Delta int
	Name  string
}

func (cmd *Incr) Parse(args []resp.Arg) error {
	if len(args) != 1 {
		return resp.ErrNumArgs(cmd.Name)
	}

	cmd.Key = args[0].String()

	return nil
}

func (cmd *Incr) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	mut := bigtable.NewReadModifyWrite()
	mut.Increment(redtable.COLUMN_FAMILY, redtable.STRING_VALUE_COLUMN, int64(cmd.Delta))
	row, err := tbl.ApplyReadModifyWrite(ctx, cmd.Key, mut)

	if err != nil {
		if status.Convert(err).Message() == "increment on non-64-bit value" {
			return nil, resp.ErrInt
		}

		return nil, err
	}

	r, ok := util.ParseRow(row)

	if !ok {
		return nil, nil
	}

	return resp.SimpleInt(int(binary.BigEndian.Uint64(r.RawValue))), nil
}
