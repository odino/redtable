package command

import (
	"context"
	"errors"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Append struct {
	Key   string
	Value string
}

func (cmd *Append) Parse(args []resp.Arg) error {
	if len(args) != 2 {
		return errors.New("wrong number of arguments for 'append' command")
	}

	cmd.Key = args[0].String()
	cmd.Value = args[1].String()

	return nil
}

func (cmd *Append) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	mut := bigtable.NewReadModifyWrite()
	mut.AppendValue("_values", "value", []byte(cmd.Value))
	row, err := tbl.ApplyReadModifyWrite(ctx, cmd.Key, mut)

	if err != nil {
		return nil, err
	}

	val, ok := util.ReadBTValue(row)

	if !ok {
		return nil, nil
	}

	return resp.SimpleInt(len(val)), nil
}
