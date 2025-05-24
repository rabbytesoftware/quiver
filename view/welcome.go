package view

import (
	"github.com/rabbytesoftware/quiver/info"
	"github.com/rabbytesoftware/quiver/logger"
)

func Welcome() {
	log := logger.NewLogger("welcome")

	log.Info(` _____       _                `)
	log.Info(`|  _  |     (_)               `)
	log.Info(`| | | |_   _ ___   _____ _ __ `)
	log.Info(`| | | | | | | \ \ / / _ \ '__|`)
	log.Info(`\ \/' / |_| | |\ V /  __/ |   `)
	log.Info(` \_/\_\\__,_|_| \_/ \___|_|   `)
	log.Info(`                              `)
	log.Info(`~~~~~~~~~ Quiver SDK ~~~~~~~~~`)
	log.Info("Developed by the people behind %s.", info.DevelopedBy)
	log.Info("Quiver and Quiver SDK is license under %s.", info.License)
	log.Info("Version: %s", info.Version)
	log.Info("Maintainers: %s", info.Maintainers)
	log.Info("")
	log.Info("")
}
