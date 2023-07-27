package main

import (
	"notes-server/controllers"
	"notes-server/db"
	"notes-server/loggers"
	"notes-server/repositories"
	"notes-server/services"
	"sync"

	logrus "github.com/sirupsen/logrus"
)

type IServiceContainer interface {
	InjectNotesController() controllers.NotesController
	InjectLoginController() controllers.LoginController
}

type kernel struct{}

func (k *kernel) InjectNotesController() controllers.NotesController {
	logrus.Infof("Notes service successfully connected!")
	logger := loggers.NewLogger()
	notesRepository := repositories.NewNotesRepository(db.NewDB(), logger)
	notesService := services.NewNotesService(logger, notesRepository)
	notesController := controllers.NewNotesController(logger, notesService)
	return notesController
}

func (k *kernel) InjectLoginController() controllers.LoginController {
	logrus.Infof("Login service successfully connected!")
	logger := loggers.NewLogger()
	loginRepository := repositories.NewLoginRepository(db.NewDB(), logger)
	loginService := services.NewLoginService(logger, loginRepository)
	loginController := controllers.NewLoginController(logger, loginService)
	return loginController
}

var (
	k             *kernel
	containerOnce sync.Once
)

func ServiceContainer() IServiceContainer {
	if k == nil {
		containerOnce.Do(func() {
			k = &kernel{}
		})
	}
	return k
}
