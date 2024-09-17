package trud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"snomed/src/shared"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func downloadPackage(prefix string, release *Release, directory string) error {
	fileURL := release.Metadata.ArchiveFileURL
	fileSize := release.Metadata.ArchiveFileSizeBytes
	fileName := release.Metadata.ArchiveFileName
	filePath := path.Join(directory, fileName)

	res, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"package<%s> at url<%s> returned status<%d> with message: %s",
			fileName, fileURL, res.StatusCode, res.Status,
		)
	}

	tmp, err := os.OpenFile(fmt.Sprintf("%s.tmp", filePath), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer tmp.Close()

	desc := fmt.Sprintf("[cyan][%s][reset] Downloading %s...", prefix, fileName)
	progress := progressbar.NewOptions64(
		int64(fileSize),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println("\n======= Download completed ==========")
		}),
	)

	_, err = io.Copy(tmp, io.TeeReader(res.Body, progress))
	if err != nil {
		return err
	}

	if err := os.Rename(filePath+".tmp", filePath); err != nil {
		return err
	}

	return nil
}

func DownloadPackages(ctx context.Context, category Category, apiKey string, directory string) error {
	if err := shared.GetOrCreateDir(directory); err != nil {
		return err
	}

	releases, err := getReleases(category, apiKey)
	if err != nil {
		return err
	}

	total := len(releases)
	index := 0
	for index < total {
		release := releases[index]

		exists, err := release.HasRelease(directory)
		if err != nil {
			return err
		} else if exists {
			total--
			releases = append(releases[:index], releases[index+1:]...)
			fmt.Printf("[%d] Skipping ReleasePackage<%s> since it already exists\n", index, release.Metadata.Name)
			continue
		}
		index++
	}

	for _, release := range releases {
		err := downloadPackage(fmt.Sprintf("%d/%d", index, total), release, directory)
		if err != nil {
			return err
		}

		output := path.Join(directory, release.Metadata.Name)
		filePath := path.Join(directory, release.Metadata.ArchiveFileName)
		if err := shared.UnzipArchive(filePath, output); err != nil {
			return err
		}
	}

	return nil
}
