package controllers

import (
	"notes-server/interfaces"
	"notes-server/loggers"
)

type NotesController struct {
	service interfaces.INotesService
	logger  *loggers.Logger
}

type LoginController struct {
	service interfaces.ILoginService
	logger  *loggers.Logger
}

func NewLoginController(logger *loggers.Logger, service interfaces.ILoginService) LoginController {
	return LoginController{
		service: service,
		logger:  logger,
	}
}

func NewNotesController(logger *loggers.Logger, service interfaces.INotesService) NotesController {
	return NotesController{
		service: service,
		logger:  logger,
	}
}
