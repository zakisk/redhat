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

func NewRouter(handler models.HandlerInterface) *Router {
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/file_exists", handler.CheckSumFile).Queries("checksum", "{[0-9a-fA-F]+}")
	getRouter.HandleFunc("/list_files", handler.ListFiles)
	getRouter.HandleFunc("/count_all_words", handler.CountAllWords)
	getRouter.HandleFunc("/get_most_frequent_words", handler.GetMostFrequentWords)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/upload_file", handler.StoreFile)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/remove_file", handler.RemoveFile).Queries("file_name", "")

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/update_file", handler.UpdateFile)

	return &Router{SM: sm}
}
