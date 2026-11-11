package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _ ArchiveProcessor = TarProcessorV1

func TarProcessorV1(path string) (*PluginData, error) {
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
		fileData := make([]byte, header.Size)

		if _, err := io.ReadFull(tarReader, fileData); err != nil {
			return nil, fmt.Errorf("error reading file %v: %w", fileName, err)
		}

		tarFs[fileName] = fileData
	}

	return tarFs, nil
}
