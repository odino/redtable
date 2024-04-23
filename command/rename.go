package command

import (
	"context"
	"errors"
	"strconv"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/redtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Rename struct {
	Src  string
	Dest string
}

func (cmd *Rename) Parse(args []resp.Arg) error {
	if len(args) != 2 {
		return resp.ErrNumArgs("rename")
	}

	cmd.Src = args[0].String()
	cmd.Dest = args[1].String()

	return nil
}

func (cmd *Rename) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	row, ok, err := util.GetRow(ctx, cmd.Src, tbl)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("no such key")
	}

	row.Key = cmd.Dest

	keys := []string{cmd.Src, cmd.Dest}
	del := bigtable.NewMutation()
	del.DeleteRow()
	write := bigtable.NewMutation()
	write.Set(redtable.COLUMN_FAMILY, redtable.STRING_VALUE_COLUMN, bigtable.ServerTime, []byte(row.Value))

	if !row.Expiry.IsZero() {
		write.Set(redtable.COLUMN_FAMILY, redtable.EXPIRY_COLUMN, bigtable.ServerTime, []byte(strconv.Itoa(int(row.Expiry.UnixMilli()))))
	}

	_, err = tbl.ApplyBulk(ctx, keys, []*bigtable.Mutation{del, write})

	return resp.OK, err
}
