package command

import (
	"context"
	"strconv"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type BitCount struct {
	Key      string
	Start    int
	HasStart bool
	End      int
	HasEnd   bool
	Index    string
}

func (cmd *BitCount) Parse(args []resp.Arg) error {
	if len(args) < 1 || len(args) > 4 {
		return resp.ErrNumArgs("bitcount")
	}

	for i, arg := range args {
		if i == 0 {
			cmd.Key = arg.String()
		}

		if arg.IsOption("BYTE") || arg.IsOption("BIT") {
			cmd.Index = arg.String()
			break
		}

		if i == 1 || i == 2 {
			pos, err := strconv.Atoi(arg.String())

			if err != nil {
				return resp.ErrInt
			}

			if i == 1 {
				cmd.HasStart = true
				cmd.Start = pos
				continue
			}

			cmd.HasEnd = true
			cmd.End = pos
			continue
		}
	}

	return nil
}

func bitSetCount(v byte) byte {
	v = (v & 0x55) + ((v >> 1) & 0x55)
	v = (v & 0x33) + ((v >> 2) & 0x33)
	return (v + (v >> 4)) & 0xF
}

func (cmd *BitCount) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	row, ok, err := util.GetRow(ctx, cmd.Key, tbl)

	if err != nil {
		return nil, err
	}

	if !ok {
		return resp.SimpleInt(0), nil
	}

	c := 0

	for i, b := range []byte(row.Value) {
		if cmd.HasStart && i < cmd.Start {
			continue
		}

		if cmd.HasEnd && i > cmd.End {
			continue
		}

		c += int(bitSetCount(b))
	}

	return resp.SimpleInt(c), nil

}
