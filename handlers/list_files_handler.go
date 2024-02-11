package handlers

import (
	"net/http"

	"github.com/zakisk/redhat-server/helpers"
	"github.com/zakisk/redhat-server/models"
)

func (h *Handler) ListFiles(rw http.ResponseWriter, r *http.Request) {
	files, err := h.fileOps.ListFile("./assets")
	if err != nil {
		h.log.Error().Msg(err.Error())
		res := &models.ListResponse{
			Success: false,
			Message: err.Error(),
			Count:   0,
			Files:   nil,
		}
		rw.WriteHeader(http.StatusInternalServerError)
		helpers.ToJSON(res, rw)
		return
	}

	res := &models.ListResponse{
		Success: true,
		Count:   len(files),
		Files:   files,
	}

	rw.WriteHeader(http.StatusOK)
	helpers.ToJSON(res, rw)
}
