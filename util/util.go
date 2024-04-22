package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/bigtable"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Row struct {
	Expiry time.Time
	Value  string
}

// Fetches a Row from bigtable, and returns:
// - the actual Row
// - whether the row was found to begin with
// - any underlying error
func GetRow(ctx context.Context, key string, tbl *bigtable.Table) (Row, bool, error) {
	row, err := tbl.ReadRow(ctx, key, bigtable.RowFilter(bigtable.LatestNFilter(1)))

	if err != nil {
		return Row{}, false, err
	}

	r, ok := ParseRow(row)

	return r, ok, nil
}

// ParseRow converts a bigtable.Row result
// into our own Row. Since the data models
// are quite different (eg. multi-columns, multi-cells, cell-timestamp)
// we want to try to make sure we can parse a BT structure
// into something that resembles a simple kv structure.
func ParseRow(row bigtable.Row) (Row, bool) {
	r := Row{}

	v, ok := row["_values"]

	if !ok {
		return r, false
	}

	var hasValue bool
	var isExpired bool

	for _, c := range v {
		if c.Column == "_values:value" {
			r.Value = string(c.Value)
			hasValue = true
		}

		if c.Column == "_values:exp" {
			ts, err := strconv.Atoi(string(c.Value))

			if err != nil {
				isExpired = true
				continue
			}

			t := time.UnixMilli(int64(ts))
			r.Expiry = t

			if err != nil {
				continue
			}

			if time.Until(t) <= 0 {
				isExpired = true
			}
		}
	}

	if !hasValue || isExpired {
		return r, false
	}

	return r, true
}

func CreateTable(project string, instance string, table string) error {
	admin, err := bigtable.NewAdminClient(context.Background(), project, instance)

	if err != nil {
		return err
	}

	err = admin.CreateTableFromConf(context.Background(), &bigtable.TableConf{
		TableID: table,
		ColumnFamilies: map[string]bigtable.Family{
			"_values": {GCPolicy: bigtable.MaxVersionsGCPolicy(1)},
		},
	})

	if err != nil {
		if status.Code(err) != codes.AlreadyExists {
			return err
		}
	}

	return nil
}

func GetTable(project string, instance string, table string) (*bigtable.Table, error) {
	client, err := bigtable.NewClient(context.Background(), project, instance)

	if err != nil {
		return nil, err
	}

	return client.Open(table), nil
}

// HandleNotHandle is the Go equivalent of "sorry not sorry"
func HandleNotHandle(err error) {
	if err != nil {
		panic(err)
	}
}

func Getenv(key string, defaults ...string) string {
	v, ok := os.LookupEnv(key)

	if !ok {
		if len(defaults) == 0 {
			panic(fmt.Sprintf("must provide env var '%s'", key))
		}

		v = defaults[0]
	}

	return v
}

// Shall we even have GC?
// Because of the distributed nature of BT, we could be reading
// an expired row while another request is writing its value,
// and we'd end up deleting the row incorrectly.
func Gc(tbl *bigtable.Table) {
	keys := []string{}
	muts := []*bigtable.Mutation{}

	err := ScanTable(tbl, func(r bigtable.Row) bool {
		_, ok := ParseRow(r)

		if !ok {
			keys = append(keys, r.Key())
			mut := bigtable.NewMutation()
			mut.DeleteRow()
			muts = append(muts, mut)
		}

		return true
	})

	if err != nil {
		log.Print(err)
		return
	}

	_, err = tbl.ApplyBulk(context.Background(), keys, muts)

	if err != nil {
		log.Print(err)
	}
}

type Scanner func(r bigtable.Row) bool

func ScanTable(tbl *bigtable.Table, f Scanner) error {
	return tbl.ReadRows(context.Background(), bigtable.InfiniteRange(""), f)
}
