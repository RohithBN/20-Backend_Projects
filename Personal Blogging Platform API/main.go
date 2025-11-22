package main

import (
	"github.com/RohithBN/database"
	"github.com/RohithBN/routes"
)


func init(){
	database.DbConnect()
}

func main() {

	r := routes.SetupRoutes()
	r.Run(":8080")

}
