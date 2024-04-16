package command

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
)

type Set struct {
	Key     string
	Value   string
	NX      bool
	XX      bool
	Get     bool
	EX      time.Time
	KeepTTL bool
}

func (cmd *Set) Parse(args []resp.Arg) error {
	if len(args) < 2 {
		return errors.New("wrong number of arguments for 'set' command")
	}

	var skip bool
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}

		if i == 0 {
			cmd.Key = string(arg)
			continue
		}

		if i == 1 {
			cmd.Value = string(arg)
			continue
		}

		if arg.IsOption("EXAT") || arg.IsOption("PEXAT") {
			return resp.ErrUnsupportedInRedtable(fmt.Sprintf("SET with %s", arg))
		}

		if arg.IsOption("EX") || arg.IsOption("PX") {
			// an expiry was already set, wtf are we doing?
			if !cmd.EX.IsZero() {
				return resp.ErrSyntax()
			}

			ts := time.Now()
			val, err := strconv.Atoi(string(args[i+1]))
			skip = true

			if err != nil {
				return errors.New("value is not an integer or out of range")
			}

			unit := time.Second

			if arg.IsOption("PX") {
				unit = time.Millisecond
			}

			cmd.EX = ts.Add(unit * time.Duration(val))
			continue
		}

		if arg.IsOption("KEEPTTL") {
			cmd.KeepTTL = true
			continue
		}

		if arg.IsOption("NX") {
			cmd.NX = true
			continue
		}

		if arg.IsOption("XX") {
			cmd.XX = true
			continue
		}

		if arg.IsOption("GET") {
			cmd.Get = true
			continue
		}

		return resp.ErrSyntax()
	}

	if !cmd.EX.IsZero() {
		if cmd.KeepTTL {
			return resp.ErrSyntax()
		}
	}

	if cmd.NX && cmd.XX {
		return resp.ErrSyntax()
	}

	return nil
}

func (cmd *Set) Run(ctx context.Context, tbl *bigtable.Table) (any, error) {
	mut := bigtable.NewMutation()
	mut.Set("_values", "value", bigtable.ServerTime, []byte(cmd.Value))

	// unless KEEPTTL is specfied, we nuke the current expiry
	if !cmd.KeepTTL {
		mut.DeleteCellsInColumn("_values", "exp")
	}

	// an expiry is set
	if !cmd.EX.IsZero() {
		mut.Set("_values", "exp", bigtable.ServerTime, []byte(cmd.EX.Format("2006-01-02 15:04:05.999999999 -0700 MST")))
	}

	var ret any

	// standard return value
	ret = resp.OK

	// these options assume you have to read
	// the current value
	if cmd.NX || cmd.XX || cmd.Get {
		// read
		get := &Get{Key: cmd.Key}
		val, err := get.Run(ctx, tbl)

		// if error, we cant go through
		if err != nil {
			return nil, err
		}

		// NX: no-op if value is set
		if val != nil && cmd.NX {
			return nil, nil
		}

		// XX: no-op if value is not set
		if val == nil && cmd.XX {
			return nil, nil
		}

		// GET: return previous value
		if cmd.Get {
			ret = val
		}
	}

	err := tbl.Apply(ctx, cmd.Key, mut)

	return ret, err
}
