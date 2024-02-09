package handlers

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/crypto/blake2b"
)

func (h *Handler) CheckSumFile(rw http.ResponseWriter, r *http.Request) {
	checksum := r.FormValue("checksum")
	if len(checksum) == 0 {
		http.Error(rw, "checksum of file doesn't exist, ensure sending it in request", http.StatusBadRequest)
		return
	}

	entries, err := os.ReadDir("./assets")
	if err != nil {
		http.Error(rw,
			fmt.Sprintf("Failed to read directory\nerror: %s", err.Error()),
			http.StatusInternalServerError)
		return
	}

	hasher, _ := blake2b.New256(nil)
	for _, e := range entries {
		info, _ := e.Info()
		f, err := os.Open("./assets/" + info.Name())
		if err != nil {
			http.Error(rw,
				fmt.Sprintf("Failed to create new file\nerror: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}
		if _, err = io.Copy(hasher, f); err != nil {
			http.Error(rw,
				fmt.Sprintf("Failed to copy file content\nerror: %s", err.Error()),
				http.StatusInternalServerError)
			return
		}

		hash := hasher.Sum(nil)
		fileChecksum := hex.EncodeToString(hash[:])
		if fileChecksum == checksum {
			http.Error(rw,
				fmt.Sprintf("Failed to copy file content\nerror: %s", err.Error()),
				http.StatusConflict) // duplicate record
			return
		}
	}
}

func (h *Handler) StoreFile(rw http.ResponseWriter, r *http.Request) {

}

func (h *Handler) RemoveFile(rw http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UpdateFile(rw http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ListFiles(rw http.ResponseWriter, r *http.Request) {

}
