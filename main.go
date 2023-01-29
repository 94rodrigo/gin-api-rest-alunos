package main

import (
	"github.com/94rodrigo/api-go-gin/database"
	"github.com/94rodrigo/api-go-gin/routes"
)

func main() {
	database.ConectaComBancoDeDados()
	routes.HandleRequests()
}
