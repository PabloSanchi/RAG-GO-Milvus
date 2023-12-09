package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"github.com/pablosanchi/datastore/core/ports"
	"github.com/pablosanchi/datastore/core/services"
	"github.com/pablosanchi/datastore/repositories/datastore"
	"github.com/pablosanchi/datastore/handler"
	restful "github.com/emicklei/go-restful/v3"
)

var(
	binding string
)

func init() {
	flag.StringVar(&binding, "binding", "localhost:8080", "Binding address")
	flag.Parse()
}

func main() {
	fmt.Println("Initializing server...")

	var datastoreRepository ports.DatastoreRepository = datastore.NewDatastoreMilvusRepository()
	
	datastoreService := services.NewDatastoreService(datastoreRepository)
	
	ws := new(restful.WebService)
	ws.Path("/api")
	handler.NewDatastoreHandler(datastoreService, ws)
	restful.Add(ws)

	fmt.Println("Listening on", binding)
	if err := http.ListenAndServe(binding, restful.DefaultContainer); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
