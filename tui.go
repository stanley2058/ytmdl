package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func drawForm(app *tview.Application, modal *tview.Modal) *tview.Form {
	curPackage := &DownloadPackage{}
	form := tview.NewForm()
	form.
		AddInputField("Download location", downloadPath, 30, nil, func(t string) { downloadPath = t }).
		AddInputField("YouTube URL", "", 40, nil, func(t string) { curPackage.URL = t }).
		AddInputField("Title", "", 20, nil, func(t string) { curPackage.Title = t }).
		AddInputField("Artist(s)", "", 20, nil, func(t string) { curPackage.Artists = t }).
		AddCheckbox("Is cover", false, func(b bool) { curPackage.IsCover = b }).
		AddButton("Start", func() {
			form.Clear(true)
			app.SetRoot(modal, true)
			go startDownload(*curPackage, func() {
				newForm := drawForm(app, modal)
				app.SetRoot(newForm, true)

				newForm.SetBackgroundColor(tcell.ColorBlack)
				app.Draw()
				newForm.SetBackgroundColor(tcell.ColorDefault)
				app.Draw()
			})
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetLabelColor(tcell.ColorGold).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorDarkCyan).
		SetButtonTextColor(tcell.ColorYellowGreen).
		SetButtonBackgroundColor(tcell.ColorBlack)

	form.SetFocus(1)
	form.SetBorder(true).
		SetTitle(" YouTube Video To MP3 ").
		SetTitleAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDefault)
	return form
}

func createTui() {
	app := tview.NewApplication()
	modal := tview.NewModal()

	modal.SetText("Downloading...")
	form := drawForm(app, modal)

	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		log.Panic(err)
	}
}
