package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"github.com/pablosanchi/datastore/core/ports"
	"github.com/pablosanchi/datastore/core/ports/secondary"
	"github.com/pablosanchi/datastore/repositories/utils"
	"github.com/pablosanchi/datastore/core/services"
	"github.com/pablosanchi/datastore/repositories/datastore"
	"github.com/pablosanchi/datastore/handler"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/emicklei/go-restful-openapi/v2"
)

var(
	hostname string
	port string
)

func main() {
	flag.StringVar(&hostname, "hostname", "localhost", "hostname address")
	flag.StringVar(&port, "port", "8080", "Port to bind")
	flag.Parse()

	address := fmt.Sprintf("%s:%s", hostname, port)

	var encoder secondary.TextEncoder = utils.NewEncoder()
	var datastoreRepository ports.DatastoreRepository = datastore.NewDatastoreMilvusRepository(encoder)
	datastoreService := services.NewDatastoreService(datastoreRepository)
	
	ws := new(restful.WebService)
	ws.Path("/api")

	handler.NewDatastoreHandler(datastoreService, ws)
	restful.Add(ws)

	// Swagger JSON
	config := restfulspec.Config{
		WebServices:    restful.DefaultContainer.RegisteredWebServices(),
		WebServicesURL: fmt.Sprintf("http://%s:%s", hostname, port),
		APIPath:        "/apidocs.json",
	}

	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))


	log.Println("Listening on", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Failed to start server: %v",  	 err)
	}
}