package view

import (
	"github.com/pterm/pterm"
)

func ProgressLoader(
	Title string,
	Prefix string,
	List []struct {
		Item     string
		Callback func() (bool, error)
	},
) {
	p, _ := pterm.DefaultProgressbar.WithTotal(len(List)).WithTitle(Title).Start()

	for i := 0; i < p.Total; i++ {
		p.UpdateTitle(Prefix + " " + List[i].Item)

		_, err := List[i].Callback()
		if err != nil {
			pterm.Error.Println(Prefix + " " + List[i].Item + ": " + err.Error())
		}

		p.Increment()
	}
}