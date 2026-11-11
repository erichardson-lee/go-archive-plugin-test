package main

import (
	"bytes"
	"os"
	"testing"
)

func getExpectedData() PluginData {
	manifestData, err := os.ReadFile("./plugin/plugin.yaml")
	if err != nil {
		panic("AAAAAAAAA!!!!")
	}
	wasmData, err := os.ReadFile("./plugin/plugin.wasm")
	if err != nil {
		panic("PANIC!")
	}

	return PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}
}

func testArchiveMethod(t *testing.T, method ArchiveProcessor, file string) {
	expected := getExpectedData()

	actual, err := method(file)

	if err != nil {
		t.Errorf("error processing: %s", err.Error())
		t.FailNow()
	}
	if !bytes.Equal(expected.manifest, actual.manifest) {
		t.Error("manifest incorrect")
		t.FailNow()
	}
	if !bytes.Equal(expected.wasm, actual.wasm) {
		t.Error("wasm incorrect")
		t.FailNow()
	}
}

func TestZip(t *testing.T) {
	testArchiveMethod(t, ZipProcessor, "plugin.zip")
}
func TestTar(t *testing.T) {
	testArchiveMethod(t, TarProcessor, "plugin.tar")
}
func TestTarGz(t *testing.T) {
	testArchiveMethod(t, TarGzProcessor, "plugin.tar.gz")
}
