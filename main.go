package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	fmt.Print("Hello world!")
}

type PluginData struct {
	manifest []byte
	wasm     []byte
}

type ArchiveProcessor func(archivePath string) (data *PluginData, err error)

var (
	_ ArchiveProcessor = ZipProcessor
	_ ArchiveProcessor = TarProcessor
	_ ArchiveProcessor = TarGzProcessor
)

func ZipProcessor(path string) (*PluginData, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	manifestData, err := fs.ReadFile(zipFile, "plugin.yaml")
	if err != nil {
		return nil, err
	}

	wasmData, err := fs.ReadFile(zipFile, "plugin.wasm")
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}

func TarProcessor(path string) (*PluginData, error) {
	tarFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()

	tarFs, err := getTarFS(tarFile)
	if err != nil {
		return nil, err
	}

	manifestData, err := tarFs.Open("plugin.yaml")
	if err != nil {
		return nil, err
	}

	wasmData, err := tarFs.Open("plugin.wasm")
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}

func TarGzProcessor(path string) (*PluginData, error) {
	gzFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer gzFile.Close()

	tarFile, err := gzip.NewReader(gzFile)
	if err != nil {
		return nil, err
	}

	tarFs, err := getTarFS(tarFile)
	if err != nil {
		return nil, err
	}

	manifestData, err := tarFs.Open("plugin.yaml")
	if err != nil {
		return nil, err
	}

	wasmData, err := tarFs.Open("plugin.wasm")
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}

type TarMemFs map[string][]byte

func (fs TarMemFs) Open(path string) ([]byte, error) {
	data, ok := fs[filepath.Clean(path)]

	if !ok {
		return nil, fmt.Errorf("file not found: %v: valid paths: %v", path, fs)
	}

	return data, nil
}

func getTarFS(tarFile io.Reader) (TarMemFs, error) {
	tarReader := tar.NewReader(tarFile)

	var tarFs = make(TarMemFs)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error getting next file: %w", err)
		}

		fileName := filepath.Clean(header.Name)
		fileData, err := io.ReadAll(tarReader)

		if err != nil {
			return nil, fmt.Errorf("error reading file %v: %w", fileName, err)
		}

		tarFs[fileName] = fileData
	}

	return tarFs, nil
}
