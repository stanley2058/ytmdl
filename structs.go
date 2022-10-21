package main

import "github.com/rivo/tview"

type DownloadPackage struct {
	URL     string
	Title   string
	Artists string
	IsCover bool
}

type OutputMetadata struct {
	FileName  string
	Title     string
	Artists   string
	Album     string
	CoverPath string
}

type TUI struct {
	Application   *tview.Application
	Form          *tview.Form
	Queue         *tview.List
	DownloadQueue []string
}
