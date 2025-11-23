package main

import (
	"github.cim/RohithBN/lib"
	"github.cim/RohithBN/routes"
)

func init(){
	lib.DbConnect()
}

func main(){

	r:= routes.SetupRoutes()
	r.Run(":8081")

}