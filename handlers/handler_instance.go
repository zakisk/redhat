package handlers

import (
	"github.com/rs/zerolog"
	"github.com/zakisk/redhat/server/models"
)

type Handler struct {
	log       zerolog.Logger
}

func NewHandlerInstance(log zerolog.Logger) models.HandlerInterface {
	return &Handler{log: log}
}
