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

func ReadBTValue(r bigtable.Row) (string, bool) {
	v, ok := r["_values"]

	if !ok {
		return "", false
	}

	var hasValue bool
	var value string
	var isExpired bool

	for _, c := range v {
		if c.Column == "_values:value" {
			value = string(c.Value)
			hasValue = true
		}

		if c.Column == "_values:exp" {
			ts, err := strconv.Atoi(string(c.Value))

			if err != nil {
				isExpired = true
				continue
			}

			t := time.UnixMilli(int64(ts))

			if err != nil {
				continue
			}

			if time.Until(t) <= 0 {
				isExpired = true
			}
		}
	}

	if !hasValue || isExpired {
		return "", false
	}

	return value, true
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
		_, ok := ReadBTValue(r)

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
