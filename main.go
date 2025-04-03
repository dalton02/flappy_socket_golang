package main

import (
	"loja-back/core/server"
	"loja-back/core/server/shared"
)

func main() {

	database, _ := server.InitConnection()
	defer database.Close()
	shared.SetDB(database)
	server.MainServer()

}
