package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createTui() TUI {
	tui := TUI{
		Application:   tview.NewApplication(),
		Form:          tview.NewForm(),
		Queue:         tview.NewList(),
		DownloadQueue: []string{},
	}

	tui.Queue.SetBorder(true).SetTitle("Queue").SetBackgroundColor(tcell.ColorDefault)
	tui.drawForm()

	flex := tview.NewFlex().
		AddItem(tui.Form, 0, 1, true).
		AddItem(tui.Queue, 30, 1, false)

	if err := tui.Application.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Panic(err)
	}

	return tui
}

func (tui *TUI) drawForm() {
	curPackage := &DownloadPackage{}

	tui.Form.
		AddInputField("Download location", downloadPath, 25, nil, func(t string) { downloadPath = t }).
		AddInputField("YouTube URL", "", 30, nil, func(t string) { curPackage.URL = t }).
		AddInputField("Title", "", 20, nil, func(t string) { curPackage.Title = t }).
		AddInputField("Artist(s)", "", 20, nil, func(t string) { curPackage.Artists = t }).
		AddCheckbox("Is cover", false, func(b bool) { curPackage.IsCover = b }).
		AddButton("Start", func() {
			url := tui.enqueue(curPackage)
			go startDownload(*curPackage, func() {
				tui.Application.QueueUpdateDraw(func() {
					tui.dequeue(url)
				})
			})

			tui.Form.Clear(true)
			tui.drawForm()
		}).
		AddButton("Quit", func() {
			tui.Application.Stop()
		}).
		SetLabelColor(tcell.ColorGold).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorDarkCyan).
		SetButtonTextColor(tcell.ColorYellowGreen).
		SetButtonBackgroundColor(tcell.ColorBlack)

	tui.Form.SetFocus(1)
	tui.Form.SetBorder(true).
		SetTitle(" YouTube Video To MP3 ").
		SetTitleAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDefault)
	tui.Application.SetFocus(tui.Form)
}

func (tui *TUI) enqueue(download *DownloadPackage) string {
	title := fmt.Sprintf("%s - %s", download.Title, download.Artists)
	if download.IsCover {
		title = fmt.Sprintf("%s (covered by %s)", download.Title, download.Artists)
	}
	tui.Queue.AddItem(title, download.URL, 0, nil)
	tui.DownloadQueue = append(tui.DownloadQueue, download.URL)
	return download.URL
}

func (tui *TUI) dequeue(url string) {
	index := 0
	for ; index < len(tui.DownloadQueue); index++ {
		if tui.DownloadQueue[index] == url {
			tui.DownloadQueue = append(tui.DownloadQueue[:index], tui.DownloadQueue[index+1:]...)
			break
		}
	}
	tui.Queue.RemoveItem(index)
}
