package handlers

import (
	"github.com/rs/zerolog"
	fileops "github.com/zakisk/redhat-server/file_ops"
	"github.com/zakisk/redhat-server/models"
)

type Handler struct {
	log       zerolog.Logger
	fileOps   *fileops.FileOps
}

func NewHandlerInstance(log zerolog.Logger, fileOps *fileops.FileOps) models.HandlerInterface {
	return &Handler{log: log, fileOps: fileOps}
}
