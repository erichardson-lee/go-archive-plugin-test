package main

import (
	"compress/gzip"
	"os"
)

var _ ArchiveProcessor = TarGzProcessorV1
var _ ArchiveProcessor = TarGzProcessorV2

func TarGzProcessorV1(path string, files *PluginFiles) (*PluginData, error) {
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

func TarGzProcessorV2(path string, files *PluginFiles) (*PluginData, error) {
	manifestData, err := getFileFromTarGzFile(path, files.manifestPath)
	if err != nil {
		return nil, err
	}

	wasmData, err := getFileFromTarGzFile(path, files.wasmPath)
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}
