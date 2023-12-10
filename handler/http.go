package handler

import (
	"github.com/pablosanchi/datastore/core/domain"
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

	ws.Route(ws.POST("/datastore/create").To(handler.CreateCollection).
		Doc("Create a new collection").
		Operation("createCollection").
		Reads(command.CollectionCommand{}).
		Writes(http.StatusCreated).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))

	ws.Route(ws.DELETE("/datastore/delete").To(handler.DeleteCollection).
		Doc("Delete a collection").
		Operation("deleteCollection").
		Reads(command.CollectionCommand{}).
		Writes(http.StatusOK).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))

		
	ws.Route(ws.GET("/datastore/list").To(handler.List).
		Doc("List all collections").
		Operation("listCollections").
		Writes([]string{}).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))

	ws.Route(ws.POST("/datastore/upsert").To(handler.UpsertDocuments).
		Doc("Upsert documents into a collection").
		Operation("upsertDocuments").
		Reads(command.UpsertDocumentsCommand{}).
		Writes(http.StatusCreated).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	
	ws.Route(ws.POST("/datastore/search").To(handler.Search).
		Doc("Search within a collection").
		Operation("searchCollection").
		Reads(command.SearchCommand{}).
		Writes([]domain.Document{}).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	

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