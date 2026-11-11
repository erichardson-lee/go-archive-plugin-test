package main

import (
	"compress/gzip"
	"os"
)

var _ ArchiveProcessor = TarGzProcessorV1

func TarGzProcessorV1(path string) (*PluginData, error) {
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
