package util

import (
	"context"
	"fmt"
	"os"
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
			ts, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(c.Value))

			if err != nil {
				continue
			}

			if time.Until(ts) <= 0 {
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
