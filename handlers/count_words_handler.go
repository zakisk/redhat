package handlers

import (
	"net/http"

	"github.com/zakisk/redhat/server/helpers"
	"github.com/zakisk/redhat/server/models"
)

func (h *Handler) CountAllWords(rw http.ResponseWriter, r *http.Request) {
	wordCount, err := h.fileOps.CountAllWords()
	if err != nil {
		res := &models.WordCountResponse{
			Success:           false,
			Message:           err.Error(),
			AllFilesProcessed: 0,
			AllWordsCount:     0,
		}
		rw.WriteHeader(http.StatusInternalServerError)
		helpers.ToJSON(res, rw)
	}

	res := &models.WordCountResponse{
		Success:           true,
		AllFilesProcessed: wordCount.TotalFileCount,
		AllWordsCount:     wordCount.TotalWordsCount,
	}
	rw.WriteHeader(http.StatusOK)
	helpers.ToJSON(res, rw)
}
