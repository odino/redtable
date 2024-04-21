package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigtable"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewResult(r bigtable.Row) Result {
	res := Result{Row: r, Key: r.Key()}

	v, ok := r["_values"]

	if !ok {
		return res
	}

	for _, c := range v {
		if c.Column == "_values:value" {
			res.Timestamp = time.UnixMicro(int64(c.Timestamp))
			res.Found = true
			res.Value = string(c.Value)
		}
	}

	return res
}

type Result struct {
	Row       bigtable.Row
	Key       string
	Timestamp time.Time
	Value     string
	Found     bool
}

func GetRow(ctx context.Context, key string, tbl *bigtable.Table) (Result, error) {
	row, err := tbl.ReadRow(
		ctx,
		key,
		bigtable.RowFilter(
			bigtable.ChainFilters(
				bigtable.LatestNFilter(1),
				bigtable.TimestampRangeFilter(time.Now(), time.Time{}),
			),
		),
	)

	return NewResult(row), err
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

// The Go equivalent of "sorry not sorry"
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

func Gc(tbl *bigtable.Table) {
	keys := []string{}
	muts := []*bigtable.Mutation{}

	err := ScanTable(tbl, func(r bigtable.Row) bool {
		row := NewResult(r)

		if !row.Found {
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
