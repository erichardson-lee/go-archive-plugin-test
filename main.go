package main

import "fmt"

func main() {
	fmt.Print("Hello world!")
}

type PluginData struct {
	manifest []byte
	wasm     []byte
}

type PluginFiles struct {
	manifestPath string
	wasmPath     string
}

type ArchiveProcessor func(archivePath string, files *PluginFiles) (data *PluginData, err error)
