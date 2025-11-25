package main

import (
	"github.com/RohithBN/lib"
	"github.com/RohithBN/routes"
)

func init(){
	lib.DbConnect()
}

func main(){
	r:=routes.SetupRoutes()
	r.Run(":8080")

}