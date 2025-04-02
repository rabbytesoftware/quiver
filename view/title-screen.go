package view

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func titleScreen() {	
	s, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("W", pterm.FgLightRed.ToStyle()),
		putils.LettersFromStringWithStyle("atcher", pterm.FgDarkGray.ToStyle()),
	).Srender()

	pterm.Println(s)
}