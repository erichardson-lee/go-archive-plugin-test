package main

import "fmt"

func main() {
	fmt.Print("Hello world!")
}

type PluginData struct {
	manifest []byte
	wasm     []byte
}

type ArchiveProcessor func(archivePath string) (data *PluginData, err error)
