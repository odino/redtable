package command

import (
	"context"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
)

type Copy struct {
	Source      string
	Destination string
	Replace     bool
}

func (cmd *Copy) Parse(args []resp.Arg) error {
	if len(args) < 2 || len(args) > 3 {
		return resp.ErrNumArgs("copy")
	}

	cmd.Source = args[0].String()
	cmd.Destination = args[1].String()

	if len(args) == 3 {
		if !args[2].IsOption("REPLACE") {
			return resp.ErrSyntax()
		}

		cmd.Replace = true
	}

	return nil
}

func (cmd *Copy) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	src, ok, err := util.GetRow(ctx, cmd.Source, tbl)

	if err != nil {
		return nil, err
	}

	if !ok {
		return resp.SimpleInt(0), nil
	}

	opts := []util.QueryOption{}
	var exists bool

	if !cmd.Replace {
		opts = append(opts, util.IfNotExistsQueryOption(&exists))
	}

	err = util.WriteRow(
		ctx,
		util.Row{Key: cmd.Destination, Value: src.Value, Expiry: src.Expiry},
		tbl,
		opts...,
	)

	res := 1

	if exists && !cmd.Replace {
		res = 0
	}

	return resp.SimpleInt(res), err
}
