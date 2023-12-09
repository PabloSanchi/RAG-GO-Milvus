package handler

import (
	"github.com/pablosanchi/datastore/core/ports"
	"github.com/pablosanchi/datastore/handler/command"
	restful "github.com/emicklei/go-restful/v3"
    "net/http"
)

type DatastoreHandler struct {
	datastoreService ports.DatastoreService
}

func NewDatastoreHandler(datastoreService ports.DatastoreService, ws *restful.WebService) *DatastoreHandler {
	handler := &DatastoreHandler{
		datastoreService: datastoreService,
	}

	ws.Route(ws.POST("/datastore/create").To(handler.CreateCollection).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.DELETE("/datastore/delete").To(handler.DeleteCollection).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/datastore/list").To(handler.List).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.POST("/datastore/upsert").To(handler.UpsertDocuments).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))
	ws.Route(ws.POST("/datastore/search").To(handler.Search).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON))

	return handler
}

func (dsh *DatastoreHandler) CreateCollection(req *restful.Request, res *restful.Response) {
	var collectionCommand command.CollectionCommand

    if err := req.ReadEntity(&collectionCommand); err != nil {
        res.WriteError(http.StatusBadRequest, err)
        return
    }

	err := dsh.datastoreService.CreateCollection(collectionCommand.CollectionName)

	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (dsh *DatastoreHandler) DeleteCollection(req *restful.Request, res *restful.Response) {
	var collectionCommand command.CollectionCommand

    if err := req.ReadEntity(&collectionCommand); err != nil {
        res.WriteError(http.StatusBadRequest, err)
        return
    }

	err := dsh.datastoreService.DeleteCollection(collectionCommand.CollectionName)

	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (dsh *DatastoreHandler) List(req *restful.Request, res *restful.Response) {
	collections, err := dsh.datastoreService.List()

	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}

	res.WriteEntity(collections)

}


func (dsh *DatastoreHandler) UpsertDocuments(req *restful.Request, res *restful.Response) {
	
	var upsertDocumentsCommand command.UpsertDocumentsCommand;

    if err := req.ReadEntity(&upsertDocumentsCommand); err != nil {
        res.WriteError(http.StatusBadRequest, err)
        return
    }


	err := dsh.datastoreService.UpsertDocuments(upsertDocumentsCommand.CollectionName, upsertDocumentsCommand.Documents)

	if err != nil { 
		res.WriteError(http.StatusInternalServerError, err)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (dsh *DatastoreHandler) Search(req *restful.Request, res *restful.Response) {
	var searchCommand command.SearchCommand

	if err := req.ReadEntity(&searchCommand); err != nil {
		res.WriteError(http.StatusBadRequest, err)
		return
	}

	results, err := dsh.datastoreService.Search(searchCommand.CollectionName, searchCommand.Query)

	if err != nil {
		res.WriteError(http.StatusInternalServerError, err)
		return
	}

	res.WriteAsJson(results)
}