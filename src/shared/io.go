package shared

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func GetOrCreateDir(directory string) error {
	if err := os.Mkdir(directory, os.ModeDir); err != nil && !os.IsExist(err) {
		return err
	}

	isDir, err := IsDirectory(directory)
	if err != nil {
		return err
	} else if !isDir {
		return fmt.Errorf("expected path<%s> to point to a directory but found a file instead", directory)
	}

	return nil
}

func CleanDirectory(directory string) error {
	isDir, err := IsDirectory(directory)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil && !isDir {
		return fmt.Errorf("expected path<%s> to point to a directory but found a file instead", directory)
	}

	dir, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(directory, name))
		if err != nil {
			return err
		}
	}

	return nil
}

func GetZipContents(filepath string) ([]string, error) {
	read, err := zip.OpenReader(filepath)
	if err != nil {
		return nil, err
	}
	defer read.Close()

	var result []string
	for _, file := range read.File {
		result = append(result, file.Name)
	}

	return result, nil
}

func UnzipArchive(input string, output string) error {
	if err := GetOrCreateDir(output); err != nil {
		return err
	}

	read, err := zip.OpenReader(input)
	if err != nil {
		return err
	}
	defer read.Close()

	for _, file := range read.File {
		fc, err := file.Open()
		if err != nil {
			return err
		}
		defer fc.Close()

		fpath := filepath.Join(output, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, file.Mode())
			continue
		}

		var fdir string
		if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
			fdir = fpath[:lastIndex]
		}

		err = os.MkdirAll(fdir, file.Mode())
		if err != nil {
			return err
		}

		f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, fc)
		if err != nil {
			return err
		}
	}

	return nil
}
