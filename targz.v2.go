package main

import (
	"compress/gzip"
	"os"
)

var _ ArchiveProcessor = TarGzProcessorV2

func TarGzProcessorV2(path string) (*PluginData, error) {

	manifestData, err := getFileFromTarGzFile(path, "plugin.yaml")
	if err != nil {
		return nil, err
	}

	wasmData, err := getFileFromTarGzFile(path, "plugin.wasm")
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}

// Now reads the full tar.gz file, skipping everything except the target file
// You may be thinking, why do a full read twice?
//
// We don't have the wasm path until we read the config, so we need to find the
// config and then find the wasm file after
func getFileFromTarGzFile(tarGzPath string, path string) ([]byte, error) {
	gzFile, err := os.Open(tarGzPath)
	if err != nil {
		return nil, err
	}
	defer gzFile.Close()

	tarFile, err := gzip.NewReader(gzFile)
	if err != nil {
		return nil, err
	}

	return getFileFromTarReader(tarFile, path)
}
