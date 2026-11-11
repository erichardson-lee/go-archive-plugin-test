package main

import (
	"archive/zip"
	"io"
)

var _ ArchiveProcessor = ZipProcessorV2

func ZipProcessorV2(path string) (*PluginData, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	var filesizes = make(map[string]uint64)
	for _, file := range zipFile.File {
		filesizes[file.Name] = file.UncompressedSize64
	}

	var manifestData = make([]byte, filesizes["plugin.yaml"])
	manifestFile, err := zipFile.Open("plugin.yaml")
	if err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(manifestFile, manifestData); err != nil {
		return nil, err
	}

	var wasmData = make([]byte, filesizes["plugin.wasm"])
	wasmFile, err := zipFile.Open("plugin.wasm")
	if err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(wasmFile, wasmData); err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}
