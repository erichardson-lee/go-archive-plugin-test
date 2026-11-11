package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//  ____
// < v1 >
//  ----
//         \   ^__^
//          \  (oo)\_______
//             (__)\       )\/\
//                 ||----w |
//                 ||     ||
//
// Builds an in-memory filesystem with the full contents of each file stored in
// a map

type TarMemFs map[string][]byte

func (fs TarMemFs) Open(path string) ([]byte, error) {
	data, ok := fs[filepath.Clean(path)]

	if !ok {
		return nil, fmt.Errorf("file not found: %v: valid paths: %v", path, fs)
	}

	return data, nil
}

func getTarFS(tarFile io.Reader) (TarMemFs, error) {
	tarReader := tar.NewReader(tarFile)

	var tarFs = make(TarMemFs)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error getting next file: %w", err)
		}

		fileName := filepath.Clean(header.Name)
		fileData := make([]byte, header.Size)

		if _, err := io.ReadFull(tarReader, fileData); err != nil {
			return nil, fmt.Errorf("error reading file %v: %w", fileName, err)
		}

		tarFs[fileName] = fileData
	}

	return tarFs, nil
}

//	____
//
// < v2 >
//
//	----
//	       \   ^__^
//	        \  (oo)\_______
//	           (__)\       )\/\
//	               ||----w |
//	               ||     ||
//
// Now reads the full tar file, skipping everything except the target file
// Because tars don't have random access, requires a full re-read each time
//
// We don't have the wasm path until we read the config, so we need to find the
// config and then find the wasm file after
// This causes two full file reads from disk
// Or 'full' reads, as they stop early once they find the file they want

func getFileFromTarFile(tarPath string, path string) ([]byte, error) {
	tarFile, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()

	return getFileFromTarReader(tarFile, path)
}
func getFileFromTarGzFile(tarGzPath string, path string) ([]byte, error) {
	gzFile, err := os.Open(tarGzPath)
	if err != nil {
		return nil, err
	}
	defer gzFile.Close()

	tarFile, err := gzip.NewReader(gzFile)
	if err != nil {
		return nil, err
	}

	return getFileFromTarReader(tarFile, path)
}

func getFileFromTarReader(reader io.Reader, path string) ([]byte, error) {
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("file not found: %v", path)
		}
		if err != nil {
			return nil, fmt.Errorf("error getting next file: %w", err)
		}

		fileName := filepath.Clean(header.Name)

		if fileName != path {
			// Skip to next file
			continue
		}

		fileData := make([]byte, header.Size)

		if _, err := io.ReadFull(tarReader, fileData); err != nil {
			return nil, fmt.Errorf("error reading file %v: %w", fileName, err)
		}

		return fileData, nil
	}
}
