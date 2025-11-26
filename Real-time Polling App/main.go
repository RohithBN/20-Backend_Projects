package main

import (
	"fmt"
	"net/http"

	"github.com/RohithBN/handler"
)

func main() {

	// create a ws endpoint
	http.HandleFunc("/ws", handler.PollHandler)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
	
}
