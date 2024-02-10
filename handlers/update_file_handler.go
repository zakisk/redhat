package handlers

import (
	"fmt"
	"net/http"

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

	err = h.fileOps.UpdateFile(header.Filename, file)
	if err != nil {
		h.log.Error().Msg(err.Error())
		res := &models.Response{
			Success: false,
			Message: err.Error(),
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
