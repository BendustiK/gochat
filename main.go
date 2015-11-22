package main

import (
	"server"
	"utils"
)

func main() {
	utils.InitLogger()

	server := server.InitServer()
	server.Start()
}
