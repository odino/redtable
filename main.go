package main

import (
	"context"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/command"
	"github.com/tidwall/redcon"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var project = "redtable-test-project"
var instance = "redtable-test-instance"
var tableName = "redtable"

func main() {
	port := getPort()

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
	err = redcon.ListenAndServe(port,
		func(conn redcon.Conn, cmd redcon.Command) {
			cmds := []string{}

			for _, c := range cmd.Args {
				cmds = append(cmds, string(c))
			}

			res, err := command.Process(strings.ToLower(cmds[0]), cmds[1:], tbl)

			if err != nil {
				conn.WriteError(err.Error())
				return
			}

			conn.WriteAny(res)

		},
		func(conn redcon.Conn) bool {
			// Use this function to accept or deny the connection.
			// log.Printf("accept: %s", conn.RemoteAddr())
			return true
		},
		func(conn redcon.Conn, err error) {
			// This is called when the connection has been closed
			// log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
		},
	)

	if err != nil {
		panic(err)
	}
}

func getPort() string {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "6379"
	}

	return ":" + port
}
