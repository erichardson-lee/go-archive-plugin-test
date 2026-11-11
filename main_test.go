package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

// pluginDir holds the original files that the test archives were built from.
const pluginDir = "plugin"

// impl is a named ArchiveProcessor version. Tests and benchmarks run each one
// as a sub-test, e.g. BenchmarkZip/V1, BenchmarkZip/V2.
type impl struct {
	name string
	fn   ArchiveProcessor
}

// fileSets are the files each impl is asked to extract. Each set gets its own
// sub-test, named after the impl plus the suffix, e.g. TestZip/V1/Big.
var fileSets = []struct {
	suffix string
	files  PluginFiles
}{
	{"", PluginFiles{manifestPath: "plugin.yaml", wasmPath: "plugin.wasm"}},
	{"/Big", PluginFiles{manifestPath: "plugin.yaml", wasmPath: "big.file"}},
}

func testImpls(t *testing.T, archivePath string, impls []impl) {
	for _, im := range impls {
		for _, set := range fileSets {
			t.Run(im.name+set.suffix, func(t *testing.T) {
				// Free this sub-test's result (up to ~1GB) before the next
				// sub-test allocates its own; otherwise the GC can hold both.
				t.Cleanup(debug.FreeOSMemory)

				got, err := im.fn(archivePath, &set.files)
				if err != nil {
					t.Fatal(err)
				}

				assertMatchesFile(t, got.manifest, set.files.manifestPath)
				assertMatchesFile(t, got.wasm, set.files.wasmPath)
			})
		}
	}
}

func benchImpls(b *testing.B, archivePath string, impls []impl) {
	for _, im := range impls {
		for _, set := range fileSets {
			b.Run(im.name+set.suffix, func(b *testing.B) {
				for b.Loop() {
					if _, err := im.fn(archivePath, &set.files); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// assertMatchesFile fails the test unless data is identical to the file at
// name inside pluginDir.
func assertMatchesFile(t *testing.T, data []byte, name string) {
	t.Helper()

	equal, err := equalsFile(data, filepath.Join(pluginDir, name))
	if err != nil {
		t.Fatal(err)
	}
	if !equal {
		t.Fatalf("extracted %s does not match the original", name)
	}
}

// equalsFile reports whether data is byte-for-byte identical to the file at
// path. The file is read and compared one chunk at a time, so checking a 1GB
// file only needs a 1MiB buffer rather than a second copy of the whole file.
func equalsFile(data []byte, path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return false, err
	}
	if info.Size() != int64(len(data)) {
		return false, nil
	}

	buf := make([]byte, 1<<20)
	for len(data) > 0 {
		chunk := buf[:min(len(buf), len(data))]
		if _, err := io.ReadFull(f, chunk); err != nil {
			return false, err
		}
		if !bytes.Equal(chunk, data[:len(chunk)]) {
			return false, nil
		}
		data = data[len(chunk):]
	}
	return true, nil
}
