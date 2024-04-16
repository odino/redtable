package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/command"
	"github.com/tidwall/redcon"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	time.Sleep(1 * time.Second)
	port := getenv("PORT", "6379")
	project := getenv("PROJECT")
	instance := getenv("INSTANCE")
	tableName := getenv("TABLE")

	admin, err := bigtable.NewAdminClient(context.Background(), project, instance)

	if err != nil {
		panic(err)
	}

	err = admin.CreateTableFromConf(context.Background(), &bigtable.TableConf{
		TableID: tableName,
		ColumnFamilies: map[string]bigtable.Family{
			"_values": {GCPolicy: bigtable.MaxVersionsGCPolicy(1)},
		},
	})

	if err != nil {
		if status.Code(err) != codes.AlreadyExists {
			panic(err)
		}
	}

	client, err := bigtable.NewClient(context.Background(), project, instance)

	if err != nil {
		panic(err)
	}

	tbl := client.Open(tableName)

	log.Printf("starting redtable server at %s", port)

	err = redcon.ListenAndServe(":"+port,
		func(conn redcon.Conn, cmd redcon.Command) {
			cmds := []string{}

			for _, c := range cmd.Args {
				cmds = append(cmds, string(c))
			}

			res, err := command.Process(cmds[0], cmds[1:], tbl)

			if err != nil {
				conn.WriteError(err.Error())
				return
			}

			conn.WriteAny(res)

		},
		func(conn redcon.Conn) bool {
			log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)

	if err != nil {
		panic(err)
	}
}

func getenv(key string, defaults ...string) string {
	v, ok := os.LookupEnv(key)

	if !ok {
		if len(defaults) == 0 {
			panic(fmt.Sprintf("must provide env var '%s'", key))
		}

		v = defaults[0]
	}

	return v
}
