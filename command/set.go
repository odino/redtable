package command

import (
	"context"
	"errors"
	"strconv"
	"time"

	"cloud.google.com/go/bigtable"
)

func errsyntax() (any, error) {
	return nil, errors.New("syntax error")
}

type Set struct {
	Key     string
	Value   string
	NX      bool
	XX      bool
	EX      time.Time
	KeepTTL bool
}

func (cmd *Set) Run(ctx context.Context, args []string, tbl *bigtable.Table) (any, error) {
	if len(args) < 2 {
		return "", errors.New("wrong number of arguments for 'set' command")
	}

	var skip bool
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}

		if i == 0 {
			cmd.Key = arg
			continue
		}

		if i == 1 {
			cmd.Value = arg
			continue
		}

		if arg == "EX" || arg == "PX" {
			// an expiry was already set, wtf are we doing?
			if !cmd.EX.IsZero() {
				return errsyntax()
			}

			ts := time.Now()
			val, err := strconv.Atoi(args[i+1])
			skip = true

			if err != nil {
				return nil, errors.New("value is not an integer or out of range")
			}

			unit := time.Second

			if arg == "PX" {
				unit = time.Millisecond
			}

			cmd.EX = ts.Add(unit * time.Duration(val))
			continue
		}

		if arg == "KEEPTTL" {
			cmd.KeepTTL = true
			continue
		}

		if arg == "NX" {
			cmd.NX = true
			continue
		}

		if arg == "XX" {
			cmd.XX = true
			continue
		}

		return errsyntax()
	}

	mut := bigtable.NewMutation()
	mut.Set("_values", "value", bigtable.ServerTime, []byte(cmd.Value))

	if !cmd.KeepTTL {
		mut.DeleteCellsInColumn("_values", "exp")
	}

	if !cmd.EX.IsZero() {
		if cmd.KeepTTL {
			return errsyntax()
		}

		mut.Set("_values", "exp", bigtable.ServerTime, []byte(cmd.EX.Format("2006-01-02 15:04:05.999999999 -0700 MST")))
	}

	if cmd.NX || cmd.XX {
		if cmd.NX && cmd.XX {
			return errsyntax()
		}

		row, err := tbl.ReadRow(ctx, cmd.Key)

		if err != nil {
			return nil, err
		}

		if len(row) != 0 && cmd.NX {
			return nil, nil
		}

		if len(row) == 0 && cmd.XX {
			return nil, nil
		}
	}

	err := tbl.Apply(ctx, cmd.Key, mut)

	return "OK", err
}
