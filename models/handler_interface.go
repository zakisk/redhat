package models

import (
	"net/http"
)

type HandlerInterface interface {
	CheckSumFile(rw http.ResponseWriter, r *http.Request)
	StoreFile(rw http.ResponseWriter, r *http.Request)
	RemoveFile(rw http.ResponseWriter, r *http.Request)
	UpdateFile(rw http.ResponseWriter, r *http.Request)
	ListFiles(rw http.ResponseWriter, r *http.Request)
	CountAllWords(rw http.ResponseWriter, r *http.Request)
	GetMostFrequentWords(rw http.ResponseWriter, r *http.Request)
}
