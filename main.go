package main

import (
	"fmt"
	"log"
	"net/http"
	"psq-project/router"
)

func main() {
	r := router.Router()

	fmt.Println("Server is Listening on the port 8080....")

	log.Fatal(http.ListenAndServe(":8080", r))
}
