package view

import (
	"rounds.com.ar/watcher/kerner"
	"rounds.com.ar/watcher/view/logger"
)

func Welcome(
	log *logger.Logger,
) {
	log.Info(` ___       __   ________  _________  ________  ___  ___  _______   ________     `)
	log.Info(`|\  \     |\  \|\   __  \|\___   ___\\   ____\|\  \|\  \|\  ___ \ |\   __  \    `)
	log.Info(`\ \  \    \ \  \ \  \|\  \|___ \  \_\ \  \___|\ \  \\\  \ \   __/|\ \  \|\  \   `)
	log.Info(` \ \  \  __\ \  \ \   __  \   \ \  \ \ \  \    \ \   __  \ \  \_|/_\ \   _  _\  `)
	log.Info(`  \ \  \|\__\_\  \ \  \ \  \   \ \  \ \ \  \____\ \  \ \  \ \  \_|\ \ \  \\  \| `)
	log.Info(`   \ \____________\ \__\ \__\   \ \__\ \ \_______\ \__\ \__\ \_______\ \__\\ _\ `)
	log.Info(`    \|____________|\|__|\|__|    \|__|  \|_______|\|__|\|__|\|_______|\|__|\|__|`)
	log.Info(`                                                                                `)
	log.Info(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ Welcome to Watcher ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`)
	log.Info("Developed by the people behind " + kerner.DevelopedBy + ".")
	log.Info("Watcher and Watcher SDK is license under " + kerner.License + ".")
	log.Info("Version: " 		+ kerner.Version)
	log.Info("Maintainers: " 	+ kerner.Maintainers)
	log.Info("")
	log.Info("")
}
