package main

import (
	"walletApp/server"
)

func main() {
	app := server.NewApp()
	app.Start()
}
