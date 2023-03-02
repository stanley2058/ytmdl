package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func getFileName(url string) (id string, audioFileName string, coverFileName string) {
	id = uuid.New().String()
	audioFileName = fmt.Sprintf("%s.mp3", string(id))
	coverFileName = fmt.Sprintf("%s.png", string(id))
	return
}

func startDownload(download DownloadPackage, callback func()) {
	os.Chdir(downloadPath)

	id, finalFileName, coverFileName := getFileName(download.URL)

	exec.Command(
		"yt-dlp",
		"-o", fmt.Sprintf("%s.%%(ext)s", id),
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

	coverPath := fmt.Sprintf("%s.cover.png", title)
	err := os.Rename(coverFileName, coverPath)
	if err != nil {
		log.Fatal(err)
	}
	convertCover(coverPath, title, download.SquareCover)

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

func convertCover(coverPath string, title string, isSquare bool) {
	width, height, err := getDimension(coverPath)
	if err != nil {
		log.Fatal(err)
	}
	margin := (width - height) / 2
	finalPath := fmt.Sprintf("%s.cover.done.png", title)

	image := ffmpeg.
		Input(coverPath).Video().
		Filter("crop", ffmpeg.Args{fmt.Sprintf("%d:%d:%d:0", height, height, margin)})

	if !isSquare {
		image = image.
			Filter("gblur", ffmpeg.Args{"sigma=50"}).
			Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", width, width)})
	}

	filter := image
	if !isSquare {
		stream := []*ffmpeg.Stream{
			image,
			ffmpeg.Input(coverPath),
		}
		filter = ffmpeg.Filter(
			stream,
			"overlay",
			ffmpeg.Args{fmt.Sprintf("0:%d", margin)},
		)
	}

	filter.
		OverWriteOutput().
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
		Run()
}
