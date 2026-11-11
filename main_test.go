package main

import (
	"bytes"
	"os"
	"testing"
)

func getTestData() (manifest []byte, wasm []byte) {
	yamlData, err := os.ReadFile("./plugin/plugin.yaml")
	if err != nil {
		panic("AAAAAAAAA!!!!")
	}
	wasmData, err := os.ReadFile("./plugin/plugin.wasm")
	if err != nil {
		panic("PANIC!")
	}

	return yamlData, wasmData
}

func testArchiveMethod(t *testing.T, method ArchiveProcessor, file string) {
	eYaml, eWasm := getTestData()

	yaml, wasm, err := method(file)

	if err != nil {
		t.Errorf("error processing: %s", err.Error())
		t.FailNow()
	}
	if !bytes.Equal(yaml, eYaml) {
		t.Error("manifest incorrect")
		t.FailNow()
	}
	if !bytes.Equal(wasm, eWasm) {
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
