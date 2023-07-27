//go:generate go run github.com/vektra/mockery/v2 --all --with-expecter --inpackage
package main

import (
	"net/http"
	"notes-server/config"

	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.Load()
	port := viper.GetString("PORT")
	logrus.Infof("Service running on port: %s", port)
	err := http.ListenAndServe(":"+port, ChiRouter().InitRouter())
	if err != nil {
		logrus.Warn("failed to setup service", err)
	}
}
