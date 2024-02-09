package models

import (
	"net/http"
)

type HandlerInterface interface {
	CheckSumFile(wr http.ResponseWriter, r *http.Request)
	StoreFile(wr http.ResponseWriter, r *http.Request)
	RemoveFile(wr http.ResponseWriter, r *http.Request)
	UpdateFile(wr http.ResponseWriter, r *http.Request)
	ListFiles(wr http.ResponseWriter, r *http.Request)
}
