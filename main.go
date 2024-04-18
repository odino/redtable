package main

import (
	"log"
	"sync/atomic"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/odino/redtable/command"
	"github.com/odino/redtable/resp"
	"github.com/odino/redtable/util"
	"github.com/tidwall/redcon"
)

var ready = &atomic.Bool{}

func main() {
	// Give redtable some time to wake up
	time.Sleep(1 * time.Second)

	// Get config from env
	port := util.Getenv("PORT", "6379")
	project := util.Getenv("PROJECT")
	instance := util.Getenv("INSTANCE")
	table := util.Getenv("TABLE")

	// Initializations
	shutdown := make(chan bool)
	err := util.CreateTable(project, instance, table)
	util.HandleNotHandle(err)

	tbl, err := util.GetTable(project, instance, table)
	util.HandleNotHandle(err)

	// HERE COMES THE FUN!
	log.Printf("starting redtable server at %s", port)
	server := redcon.NewServerNetwork("tcp", ":"+port, getHandler(tbl, shutdown), onConnect, nil)
	go func() {
		err = server.ListenAndServe()
		util.HandleNotHandle(err)
	}()
	ready.Store(true)

	// shutting down
	<-shutdown
	handleShutdown(ready, server)
}

func handleShutdown(ready *atomic.Bool, server *redcon.Server) {
	log.Printf("shutdown sequence initiated 🚀")
	ready.Store(false)
	time.Sleep(time.Second * 1)
	err := server.Close()

	if err != nil {
		log.Print(err)
	}

	log.Printf("goodbye chief 👋")
}

type handler func(conn redcon.Conn, cmd redcon.Command)

func getHandler(tbl *bigtable.Table, shutdown chan bool) handler {
	return func(conn redcon.Conn, cmd redcon.Command) {
		if !ready.Load() {
			conn.WriteError("server not ready")
			return
		}

		cmds := []resp.Arg{}

		for _, c := range cmd.Args {
			cmds = append(cmds, resp.Arg(c))
		}

		res, err := command.Process(string(cmds[0]), cmds[1:], tbl)

		if err != nil {
			if err == resp.ErrShutdown {
				conn.WriteString("AS YOU WISH")
				shutdown <- true
				return
			}

			conn.WriteError(err.Error())
			return
		}

		conn.WriteAny(res)
	}
}

func onConnect(conn redcon.Conn) bool {
	if !ready.Load() {
		log.Printf("reject: %s (server not ready)", conn.RemoteAddr())
		return false
	}

	log.Printf("accept: %s", conn.RemoteAddr())
	return true
}
