package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _ ArchiveProcessor = TarProcessorV2

func TarProcessorV2(path string) (*PluginData, error) {
	manifestData, err := getFileFromTarFile(path, "plugin.yaml")
	if err != nil {
		return nil, err
	}

	wasmData, err := getFileFromTarFile(path, "plugin.wasm")
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}

// Now reads the full tar file, skipping everything except the target file
// You may be thinking, why do a full read twice?
//
// We don't have the wasm path until we read the config, so we need to find the
// config and then find the wasm file after
func getFileFromTarFile(tarPath string, path string) ([]byte, error) {
	tarFile, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()

	return getFileFromTarReader(tarFile, path)
}

func getFileFromTarReader(reader io.Reader, path string) ([]byte, error) {
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("file not found: %v", path)
		}
		if err != nil {
			return nil, fmt.Errorf("error getting next file: %w", err)
		}

		fileName := filepath.Clean(header.Name)

		if fileName != path {
			// Skip to next file
			continue
		}

		fileData := make([]byte, header.Size)

		if _, err := io.ReadFull(tarReader, fileData); err != nil {
			return nil, fmt.Errorf("error reading file %v: %w", fileName, err)
		}

		return fileData, nil
	}
}
