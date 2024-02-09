package router

import (
	// "fmt"
	// "net/http"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/zakisk/redhat/server/models"
)

type Router struct {
	SM *mux.Router
}

func NewRouter(handler models.HandlerInterface) (*Router, error) {
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/file_exists", handler.CheckSumFile)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/upload", handler.StoreFile)

	return &Router{SM: sm}, nil
}
