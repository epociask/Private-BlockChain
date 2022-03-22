package main

import (
	"chain/uri/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"chain/uri"
)

var (
	Router = uri.NewRouter()
)

func main() {

	var port string
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "6689"
	}

	handlers.InitSelfAddress(port)

	fmt.Println("running: " + "localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, Router))
}
