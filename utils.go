package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"os/exec"
	"regexp"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func getFileName(url string) (audioFileName string, coverFileName string) {
	cmd := exec.Command(
		"yt-dlp",
		"--get-filename",
		"-o", "%(title)s.%(ext)s",
		"-x",
		"--audio-format", "mp3",
		url,
	)

	fName, err := cmd.Output()
	if err != nil {
		log.Panic(err)
	}

	match := regexp.MustCompile(`\.[^.]+$`)
	audioFileName = match.ReplaceAllString(string(fName), ".mp3")
	coverFileName = match.ReplaceAllString(string(fName), ".png")
	return
}

func startDownload(download DownloadPackage, callback func()) {
	os.Chdir(downloadPath)
	coverPath := "cover.png"

	finalFileName, coverFileName := getFileName(download.URL)

	exec.Command(
		"yt-dlp",
		"-o", "%(title)s.%(ext)s",
		"-x",
		"--audio-format", "mp3",
		"--write-thumbnail",
		"--convert-thumbnails", "png",
		download.URL,
	).Run()

	title := fmt.Sprintf("%s - %s", download.Title, download.Artists)
	if download.IsCover {
		title = fmt.Sprintf("%s (covered by %s)", download.Title, download.Artists)
	}

	err := os.Rename(coverFileName, coverPath)
	if err != nil {
		log.Fatal(err)
	}
	convertCover(coverPath)

	addMetadataAndCover(OutputMetadata{
		FileName:  finalFileName,
		Title:     title,
		Artists:   download.Artists,
		Album:     download.Title,
		CoverPath: coverPath,
	})

	os.Remove(finalFileName)
	os.Remove(coverPath)

	callback()
}

func getDimension(imgPath string) (width int, height int, err error) {
	reader, err := os.Open(imgPath)
	if err == nil {
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			log.Fatal(err)
		}
		return im.Width, im.Height, nil
	}
	return 0, 0, err
}

func convertCover(coverPath string) {
	width, height, err := getDimension(coverPath)
	if err != nil {
		log.Fatal(err)
	}
	margin := (width - height) / 2
	finalPath := "cover.done.png"

	fmt.Println(width, height, err)

	ffmpeg.Filter([]*ffmpeg.Stream{
		ffmpeg.Input(coverPath).Video().
			Filter("crop", ffmpeg.Args{fmt.Sprintf("%d:%d:%d:0", height, height, margin)}).
			Filter("gblur", ffmpeg.Args{"sigma=50"}).
			Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", width, width)}),
		ffmpeg.Input(coverPath),
	},
		"overlay",
		ffmpeg.Args{fmt.Sprintf("0:%d", margin)},
	).
		OverWriteOutput().
		ErrorToStdOut().
		Output(finalPath).
		Run()

	os.Rename(finalPath, coverPath)
}

func addMetadataAndCover(metadata OutputMetadata) {
	ffmpeg.Output(
		[]*ffmpeg.Stream{
			ffmpeg.Input(metadata.FileName),
			ffmpeg.Input(metadata.CoverPath),
		},
		fmt.Sprintf("%s.mp3", metadata.Title),
		ffmpeg.KwArgs{
			"codec":         "copy",
			"id3v2_version": "3",
			"metadata": []string{
				fmt.Sprintf("title=%s", metadata.Title),
				fmt.Sprintf("artist=%s", metadata.Artists),
				fmt.Sprintf("Album=%s", metadata.Album),
			},
			"metadata:s:v": []string{
				`title="Album cover"`,
				`comment="Cover (front)"`,
			},
		},
	).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
}
