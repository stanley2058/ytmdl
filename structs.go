package main

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
