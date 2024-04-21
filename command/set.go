package command

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
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

var NO_EXPIRY_TS = time.Now().Add(time.Hour * 24 * 365 * 100)

func (cmd *Set) Parse(args []resp.Arg) error {
	if len(args) < 2 {
		return resp.ErrNumArgs("set")
	}

	cmd.EX = NO_EXPIRY_TS

	var skip bool
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}

		if i == 0 {
			cmd.Key = arg.String()
			continue
		}

		if i == 1 {
			cmd.Value = arg.String()
			continue
		}

		if arg.IsOption("EXAT") || arg.IsOption("PEXAT") {
			return resp.ErrUnsupportedInRedtable(fmt.Sprintf("SET with %s", arg))
		}

		if arg.IsOption("EX") || arg.IsOption("PX") {
			// an expiry was already set, wtf are we doing?
			if cmd.EX != NO_EXPIRY_TS {
				return resp.ErrSyntax()
			}

			ts := time.Now()
			val, err := strconv.Atoi(args[i+1].String())
			skip = true

			if err != nil {
				return resp.ErrInt
			}

			if val <= 0 {
				return resp.ErrInvalidExpire("set")
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

	if cmd.EX != NO_EXPIRY_TS {
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
	var ret any

	// standard return value
	ret = resp.OK

	// these options assume you have to read
	// the current value
	if cmd.NX || cmd.XX || cmd.Get || cmd.KeepTTL {
		// read
		res, err := util.GetRow(ctx, cmd.Key, tbl)

		// if error, we cant go through
		if err != nil {
			return nil, err
		}

		// NX: no-op if value is set
		if res.Found && cmd.NX {
			return nil, nil
		}

		// XX: no-op if value is not set
		if !res.Found && cmd.XX {
			return nil, nil
		}

		// GET: return previous value
		if cmd.Get {
			if !res.Found {
				ret = nil
			} else {
				ret = res.Value
			}

		}

		if cmd.KeepTTL && res.Found {
			cmd.EX = res.Timestamp
		}
	}

	mut := bigtable.NewMutation()
	mut.DeleteRow()
	mut.Set("_values", "value", bigtable.Timestamp(cmd.EX.UnixMicro()), []byte(cmd.Value))

	err := tbl.Apply(ctx, cmd.Key, mut)

	return ret, err
}
