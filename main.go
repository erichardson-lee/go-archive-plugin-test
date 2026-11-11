package main

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing/fstest"
)

func main() {
	fmt.Print("Hello world!")
}

type ArchiveProcessor func(archivePath string) (manifest []byte, data []byte, err error)

var (
	_ ArchiveProcessor = ZipProcessor
	_ ArchiveProcessor = TarProcessor
	_ ArchiveProcessor = TarGzProcessor
)

func ZipProcessor(path string) ([]byte, []byte, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, nil, err
	}
	defer zipFile.Close()

	manifestData, err := fs.ReadFile(zipFile, "plugin.yaml")
	if err != nil {
		return nil, nil, err
	}

	wasmData, err := fs.ReadFile(zipFile, "plugin.wasm")
	if err != nil {
		return nil, nil, err
	}

	return manifestData, wasmData, nil
}

func TarProcessor(path string) ([]byte, []byte, error) {
	tarFile, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer tarFile.Close()

	tarFs, err := getTarFS(tarFile)

	manifestData, err := fs.ReadFile(tarFs, "plugin.yaml")
	if err != nil {
		return nil, nil, err
	}

	wasmData, err := fs.ReadFile(tarFs, "plugin.wasm")
	if err != nil {
		return nil, nil, err
	}

	return manifestData, wasmData, nil
}

func TarGzProcessor(path string) ([]byte, []byte, error) {
	return nil, nil, errors.New("not implemented")
}

func getTarFS(tarFile io.Reader) (fs.FS, error) {
	tarReader := tar.NewReader(tarFile)

	var tarFs = make(fstest.MapFS)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		fileName := filepath.Clean(header.Name)
		fileData, err := io.ReadAll(tarFile)

		if err != nil {
			return nil, err
		}

		tarFs[fileName] = &fstest.MapFile{
			Data:    fileData,
			Mode:    header.FileInfo().Mode(),
			ModTime: header.ModTime,
			Sys:     header.FileInfo().Sys(),
		}
	}

	return tarFs, nil
}
