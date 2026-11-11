package main

import (
	"os"
)

var _ ArchiveProcessor = TarProcessorV1
var _ ArchiveProcessor = TarProcessorV2

func TarProcessorV1(path string, files *PluginFiles) (*PluginData, error) {
	tarFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()

	tarFs, err := getTarFS(tarFile)
	if err != nil {
		return nil, err
	}

	manifestData, err := tarFs.Open(files.manifestPath)
	if err != nil {
		return nil, err
	}

	wasmData, err := tarFs.Open(files.wasmPath)
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}

func TarProcessorV2(path string, files *PluginFiles) (*PluginData, error) {
	manifestData, err := getFileFromTarFile(path, files.manifestPath)
	if err != nil {
		return nil, err
	}

	wasmData, err := getFileFromTarFile(path, files.wasmPath)
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}
