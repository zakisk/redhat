package handlers

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/zakisk/redhat/server/helpers"
	"github.com/zakisk/redhat/server/models"
	"golang.org/x/crypto/blake2b"
)

func (h *Handler) CheckSumFile(rw http.ResponseWriter, r *http.Request) {
	checksum := r.URL.Query().Get("checksum")
	if len(checksum) == 0 {
		msg := "checksum of file doesn't exist in query, ensure sending it"
		h.log.Error().Msg(msg)
		res := &models.Response{
			Success: false,
			Message: msg,
		}
		rw.WriteHeader(http.StatusBadRequest)
		helpers.ToJSON(res, rw)
		return
	}

	entries, err := os.ReadDir("./assets")
	if err != nil {
		h.log.Error().Msg(err.Error())
		res := &models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to read directory\nerror: %s", err.Error()),
		}
		rw.WriteHeader(http.StatusInternalServerError)
		helpers.ToJSON(res, rw)
		return
	}

	hasher, _ := blake2b.New256(nil)

	for _, e := range entries {
		info, _ := e.Info()
		f, err := os.Open("./assets/" + info.Name())
		if err != nil {
			h.log.Error().Msg(err.Error())
			res := &models.Response{
				Success: false,
				Message: fmt.Sprintf("Failed to open file\nerror: %s", err.Error()),
			}
			rw.WriteHeader(http.StatusInternalServerError)
			helpers.ToJSON(res, rw)
			return
		}
		if _, err = io.Copy(hasher, f); err != nil {
			h.log.Error().Msg(err.Error())
			res := &models.Response{
				Success: false,
				Message: fmt.Sprintf("Failed to copy file content\nerror: %s", err.Error()),
			}
			rw.WriteHeader(http.StatusInternalServerError)
			helpers.ToJSON(res, rw)
			return
		}

		hash := hasher.Sum(nil)
		fileChecksum := hex.EncodeToString(hash)
		if fileChecksum == checksum {
			res := &models.Response{
				Success:  true,
				Metadata: map[string]interface{}{"checksum_exists": true},
			}
			rw.WriteHeader(http.StatusOK)
			helpers.ToJSON(res, rw)
			return
		}
	}

	res := &models.Response{
		Success:  true,
		Metadata: map[string]interface{}{"checksum_exists": false},
	}
	rw.WriteHeader(http.StatusOK)
	helpers.ToJSON(res, rw)
}
