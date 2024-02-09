package router

import (
	// "fmt"
	// "net/http"

	"net/http"

	"github.com/gorilla/mux"
	// "github.com/tus/tusd/v2/pkg/filestore"
	// tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/zakisk/redhat/server/models"
)

type Router struct {
	SM *mux.Router
}

func NewRouter(handler models.HandlerInterface) (*Router, error) {
	sm := mux.NewRouter()

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/upload", handler.StoreFile)

	// store := filestore.FileStore{
	// 	Path: "./assets",
	// }
	// composer := tusd.NewStoreComposer()
	// store.UseIn(composer)
	// corsConfig := tusd.DefaultCorsConfig
	// corsConfig.Disable = true
	// tusHandler, err := tusd.NewHandler(tusd.Config{
	// 	BasePath:              "/file/",
	// 	StoreComposer:         composer,
	// 	NotifyCompleteUploads: true,
	// 	Cors:                  &corsConfig,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// go func() {
	// 	for {
	// 		event := <-tusHandler.CompleteUploads
	// 		fmt.Printf("Upload %s finished\n", event.Upload.ID)
	// 	}
	// }()

	// sm.Handle("/file/upload", http.StripPrefix("/file/", tusHandler))

	return &Router{SM: sm}, nil
}
