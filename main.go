package main

import (
	"log"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/command"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
	"github.com/tidwall/redcon"
)

func main() {
	// Give redtable some time to wake up
	time.Sleep(1 * time.Second)

	// Get config from env
	port := util.Getenv("PORT", "6379")
	project := util.Getenv("PROJECT")
	instance := util.Getenv("INSTANCE")
	table := util.Getenv("TABLE")

	// Initializations
	err := util.CreateTable(project, instance, table)
	util.HandleNotHandle(err)

	tbl, err := util.GetTable(project, instance, table)
	util.HandleNotHandle(err)

	// HERE COMES THE FUN!
	log.Printf("starting redtable server at %s", port)

	err = redcon.ListenAndServe(":"+port, getHandler(tbl), onConnect, onClose)

	if err != nil {
		panic(err)
	}
}

type handler func(conn redcon.Conn, cmd redcon.Command)

func getHandler(tbl *bigtable.Table) handler {
	return func(conn redcon.Conn, cmd redcon.Command) {
		cmds := []resp.Arg{}

		for _, c := range cmd.Args {
			cmds = append(cmds, resp.Arg(c))
		}

		res, err := command.Process(string(cmds[0]), cmds[1:], tbl)

		if err != nil {
			conn.WriteError(err.Error())
			return
		}

		conn.WriteAny(res)
	}
}

func onConnect(conn redcon.Conn) bool {
	log.Printf("accept: %s", conn.RemoteAddr())
	return true
}

func onClose(conn redcon.Conn, err error) {
	log.Printf("closed: %s, err: %v", conn.RemoteAddr(), err)
}
