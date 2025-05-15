package view

import (
	"rounds.com.ar/watcher/kerner"
	"rounds.com.ar/watcher/view/logger"
)

func Welcome(
	log *logger.Logger,
) {
	log.Info("welcome", ` ___       __   ________  _________  ________  ___  ___  _______   ________     `)
	log.Info("welcome", `|\  \     |\  \|\   __  \|\___   ___\\   ____\|\  \|\  \|\  ___ \ |\   __  \    `)
	log.Info("welcome", `\ \  \    \ \  \ \  \|\  \|___ \  \_\ \  \___|\ \  \\\  \ \   __/|\ \  \|\  \   `)
	log.Info("welcome", ` \ \  \  __\ \  \ \   __  \   \ \  \ \ \  \    \ \   __  \ \  \_|/_\ \   _  _\  `)
	log.Info("welcome", `  \ \  \|\__\_\  \ \  \ \  \   \ \  \ \ \  \____\ \  \ \  \ \  \_|\ \ \  \\  \| `)
	log.Info("welcome", `   \ \____________\ \__\ \__\   \ \__\ \ \_______\ \__\ \__\ \_______\ \__\\ _\ `)
	log.Info("welcome", `    \|____________|\|__|\|__|    \|__|  \|_______|\|__|\|__|\|_______|\|__|\|__|`)
	log.Info("welcome", `                                                                                `)
	log.Info("welcome", `~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ Welcome to Watcher ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`)
	log.Info("welcome", "Developed by the people behind " + kerner.DevelopedBy + ".")
	log.Info("welcome", "Watcher and Watcher SDK is license under " + kerner.License + ".")
	log.Info("welcome", "Version: " 		+ kerner.Version)
	log.Info("welcome", "Maintainers: " 	+ kerner.Maintainers)
	log.Info("welcome", "")
	log.Info("welcome", "")
}
