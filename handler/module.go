package handler

import (
	"log"

	"github.com/rudbast/my-net/core"
)

type (
	Module struct {
		Logger  *log.Logger
		Service *core.Service
	}
)

func New(lg *log.Logger, svc *core.Service) *Module {
	return &Module{
		Logger:  lg,
		Service: svc,
	}
}
