package main

var _ ArchiveProcessor = TarProcessorV2

func TarProcessorV2Big(path string) (*PluginData, error) {
	manifestData, err := getFileFromTarFile(path, "plugin.yaml")
	if err != nil {
		return nil, err
	}

	wasmData, err := getFileFromTarFile(path, "big.file")
	if err != nil {
		return nil, err
	}

	return &PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}, nil
}
