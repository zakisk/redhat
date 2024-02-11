package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/zakisk/redhat-server/helpers"
	"github.com/zakisk/redhat-server/models"
)

func (h *Handler) GetMostFrequentWords(rw http.ResponseWriter, r *http.Request) {
	w := r.URL.Query().Get("words")
	order := r.URL.Query().Get("order")
	if len(order) == 0 {
		order = "asc" //default
	}
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
	fmt.Printf("order: %s\n", order)
	if order != "asc" && order != "dsc" {
		res := &models.FrequentWordsResponse{
			Success: false,
			Message: fmt.Sprintf("order `%s` is invalid", order),
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

	sortedFrequentWords := sortMap(frequentWords, order)

	msg := ""
	if wordCount.TotalFileCount == 0 {
		msg = "There is no file on server"
	}

	res := &models.FrequentWordsResponse{
		Success: true,
		Message: msg,
		Words:   sortedFrequentWords,
	}
	rw.WriteHeader(http.StatusOK)
	helpers.ToJSON(res, rw)
}

func sortMap(wordMap map[string]int, order string) map[string]int {
	newMap := map[string]int{}
	keys := make([]string, 0, len(wordMap))

	for k := range wordMap {
		keys = append(keys, k)
	}

	if order == "asc" {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}

	for _, key := range keys {
		newMap[key] = wordMap[key]
	}

	return newMap
}
