package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func getExpectedData(files *PluginFiles) PluginData {
	manifestData, err := os.ReadFile(filepath.Join("./plugin/", files.manifestPath))
	if err != nil {
		panic("AAAAAAAAA!!!!")
	}

	wasmData, err := os.ReadFile(filepath.Join("./plugin/", files.wasmPath))
	if err != nil {
		panic("PANIC!")
	}

	return PluginData{
		manifest: manifestData,
		wasm:     wasmData,
	}
}

func checkData(t *testing.T, actual *PluginData, files *PluginFiles) {
	expected := getExpectedData(files)
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

var normalPaths = PluginFiles{
	manifestPath: "plugin.yaml",
	wasmPath:     "plugin.wasm",
}
var bigPaths = PluginFiles{
	manifestPath: "plugin.yaml",
	wasmPath:     "big.file",
}

func testImpls(t *testing.T, archivePath string, impls []impl) {

	for _, im := range impls {
		t.Run(im.name, func(t *testing.T) {
			result, err := im.fn(archivePath, &normalPaths)
			if err != nil {
				t.Fatal(err)
			}

			checkData(t, result, &normalPaths)
		})

		t.Run(im.name+"/Big", func(t *testing.T) {
			// I tried testing the big ones, but it nearly crashed my PC lowkey
			// (I think because it's loading the data into memory twice)
			t.SkipNow()

			result, err := im.fn(archivePath, &bigPaths)
			if err != nil {
				t.Fatal(err)
			}

			checkData(t, result, &bigPaths)
		})
	}
}

func benchImpls(b *testing.B, archivePath string, impls []impl) {
	for _, im := range impls {
		b.Run(im.name, func(b *testing.B) {
			for b.Loop() {
				if _, err := im.fn(archivePath, &normalPaths); err != nil {
					// Should never error
					b.Fatal(err)
				}
			}
		})

		b.Run(im.name+"/Big", func(b *testing.B) {
			for b.Loop() {
				if _, err := im.fn(archivePath, &bigPaths); err != nil {
					// Should never error
					b.Fatal(err)
				}
			}
		})
	}
}
