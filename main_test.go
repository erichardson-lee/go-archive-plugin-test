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

// impl is a named ArchiveProcessor version. Tests and benchmarks run each one
// as a sub-test, e.g. BenchmarkZip/V1, BenchmarkZip/V2.
type impl struct {
	name string
	fn   ArchiveProcessor
}

func testImpls(t *testing.T, archivePath string, impls []impl) {
	for _, im := range impls {
		t.Run(im.name, func(t *testing.T) {
			result, err := im.fn(archivePath)
			if err != nil {
				t.Fatal(err)
			}

			checkData(t, result)
		})
	}
}

func benchImpls(b *testing.B, archivePath string, impls []impl) {
	for _, im := range impls {
		b.Run(im.name, func(b *testing.B) {
			for b.Loop() {
				if _, err := im.fn(archivePath); err != nil {
					// Should never error
					b.Fatal(err)
				}
			}
		})
	}
}
