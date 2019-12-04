package main

import (
	"./uri"
	"fmt"
	"log"
	"net/http"
	"os"
	"./uri/handlers"
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

	fmt.Println("running: "+ "localhost:"+port)
	log.Fatal(http.ListenAndServe(":"+port, Router))
}