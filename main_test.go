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

func checkData(t *testing.T, actual *PluginData) {
	expected := getExpectedData()
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
	result, err := ZipProcessor("plugin.zip")
	if err != nil {
		t.FailNow()
	}

	checkData(t, result)
}

func BenchmarkZip(b *testing.B) {
	for b.Loop() {
		_, err := ZipProcessor("plugin.zip")

		if err != nil {
			// Should never error
			b.FailNow()
		}
	}
}

func TestTar(t *testing.T) {
	result, err := TarProcessor("plugin.tar")
	if err != nil {
		t.FailNow()
	}

	checkData(t, result)
}

func BenchmarkTar(b *testing.B) {
	for b.Loop() {
		_, err := TarProcessor("plugin.tar")

		if err != nil {
			// Should never error
			b.FailNow()
		}
	}
}

func TestTarGz(t *testing.T) {
	result, err := TarGzProcessor("plugin.tar.gz")
	if err != nil {
		t.FailNow()
	}

	checkData(t, result)
}

func BenchmarkTarGz(b *testing.B) {
	for b.Loop() {
		_, err := TarGzProcessor("plugin.tar.gz")

		if err != nil {
			// Should never error
			b.FailNow()
		}
	}
}
