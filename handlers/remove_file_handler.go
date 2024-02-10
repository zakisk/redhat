package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/zakisk/redhat/server/helpers"
	"github.com/zakisk/redhat/server/models"
)

func (h *Handler) RemoveFile(rw http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file_name")
	if len(fileName) == 0 {
		msg := "File name doesn't exist in query, ensure sending it"
		h.log.Error().Msg(msg)
		res := &models.Response{
			Success: false,
			Message: msg,
		}
		rw.WriteHeader(http.StatusBadRequest)
		helpers.ToJSON(res, rw)
		return
	}

	err := h.fileOps.RemoveFile(fileName)
	if err != nil {
		h.log.Error().Msg(err.Error())
		res := &models.Response{Success: false}
		if errors.Is(err, os.ErrNotExist) {
			rw.WriteHeader(http.StatusNotFound)
			res.Message = fmt.Sprintf("No such file `%s`", fileName)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			res.Message = fmt.Sprintf("Failed to retrieve information of file\nerror: `%s`", err.Error())
		}
		helpers.ToJSON(res, rw)
		return
	}

	res := &models.Response{
		Success: true,
		Message: fmt.Sprintf("File `%s` is removed successfully", fileName),
	}
	rw.WriteHeader(http.StatusOK)
	helpers.ToJSON(res, rw)
}
