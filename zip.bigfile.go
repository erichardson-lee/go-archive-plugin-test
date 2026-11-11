package main

import (
	"archive/zip"
	"io/fs"
)

var _ ArchiveProcessor = ZipProcessorV1

func ZipProcessorV1Big(path string) (*PluginData, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	manifestData, err := fs.ReadFile(zipFile, "plugin.yaml")
	if err != nil {
		return nil, err
	}

	wasmData, err := fs.ReadFile(zipFile, "big.file")
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}
