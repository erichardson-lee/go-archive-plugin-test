package main

import (
	"archive/zip"
	"io"
	"io/fs"
)

var _ ArchiveProcessor = ZipProcessorV1
var _ ArchiveProcessor = ZipProcessorV2

func ZipProcessorV1(path string, files *PluginFiles) (*PluginData, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	manifestData, err := fs.ReadFile(zipFile, files.manifestPath)
	if err != nil {
		return nil, err
	}

	wasmData, err := fs.ReadFile(zipFile, files.wasmPath)
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}

func ZipProcessorV2(path string, files *PluginFiles) (*PluginData, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	var filesizes = make(map[string]uint64)
	for _, file := range zipFile.File {
		filesizes[file.Name] = file.UncompressedSize64
	}

	var manifestData = make([]byte, filesizes[files.manifestPath])
	manifestFile, err := zipFile.Open(files.manifestPath)
	if err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(manifestFile, manifestData); err != nil {
		return nil, err
	}

	var wasmData = make([]byte, filesizes[files.wasmPath])
	wasmFile, err := zipFile.Open(files.wasmPath)
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
