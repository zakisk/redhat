package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/zakisk/redhat/server/helpers"
	"github.com/zakisk/redhat/server/models"
)

func (h *Handler) UpdateFile(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200 << 20) // limit 200MB
	if err != nil {
		h.log.Error().Msg(err.Error())
		res := &models.Response{
			Success: false,
			Message: "Failed to parse multipart form",
		}
		rw.WriteHeader(http.StatusBadRequest)
		helpers.ToJSON(res, rw)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.log.Error().Msg(err.Error())
		res := &models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to get file from form data\nerror: %s", err.Error()),
		}
		rw.WriteHeader(http.StatusBadRequest)
		helpers.ToJSON(res, rw)
		return
	}
	defer file.Close()

	currentFile, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		h.log.Error().Msg(err.Error())
		res := &models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to open destination file\nerror: %s", err.Error()),
		}
		rw.WriteHeader(http.StatusInternalServerError)
		helpers.ToJSON(res, rw)
		return
	}
	defer currentFile.Close()

	_, err = io.Copy(currentFile, file)
	if err != nil {
		h.log.Error().Msg(err.Error())
		res := &models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to copy data into file\nerror: %s", err.Error()),
		}
		rw.WriteHeader(http.StatusInternalServerError)
		helpers.ToJSON(res, rw)
		return
	}

	res := &models.Response{
		Success: true,
		Message: fmt.Sprintf("File's content (`%s`) is updated successfully", header.Filename),
	}
	rw.WriteHeader(http.StatusOK)
	helpers.ToJSON(res, rw)
}
