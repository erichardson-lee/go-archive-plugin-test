package main

var _ ArchiveProcessor = TarGzProcessorV2

func TarGzProcessorV2Big(path string) (*PluginData, error) {

	manifestData, err := getFileFromTarGzFile(path, "plugin.yaml")
	if err != nil {
		return nil, err
	}

	wasmData, err := getFileFromTarGzFile(path, "big.file")
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}
