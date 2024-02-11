package handlers

import (
	"net/http"
	"strconv"

	"github.com/zakisk/redhat/server/helpers"
	"github.com/zakisk/redhat/server/models"
)

func (h *Handler) GetMostFrequentWords(rw http.ResponseWriter, r *http.Request) {
	w := r.URL.Query().Get("words")
	words, _ := strconv.Atoi(w)
	if words <= 0 {
		res := &models.FrequentWordsResponse{
			Success: false,
			Message: "words count should be greater than zero",
		}
		rw.WriteHeader(http.StatusBadRequest)
		helpers.ToJSON(res, rw)
		return
	}

	wordCount, err := h.fileOps.CountAllWords()
	if err != nil {
		res := &models.FrequentWordsResponse{
			Success: false,
			Message: err.Error(),
		}
		rw.WriteHeader(http.StatusInternalServerError)
		helpers.ToJSON(res, rw)
	}

	frequentWords := make(map[string]int, words)
	for i := 0; i < words; i++ {
		maxK, maxV := "", 0
		for k, v := range wordCount.WordsCountMap {
			if v > maxV {
				maxK = k
				maxV = v
			}
		}
		frequentWords[maxK] = maxV
		delete(wordCount.WordsCountMap, maxK)
	}

	msg := ""
	if wordCount.TotalFileCount == 0 {
		msg = "There is no file on server"
	}

	res := &models.FrequentWordsResponse{
		Success: true,
		Message: msg,
		Words:   frequentWords,
	}
	rw.WriteHeader(http.StatusOK)
	helpers.ToJSON(res, rw)
}
