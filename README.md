# Archive Filesystem Test

This is a quick test repo to see how different archive formats work with go

## Methodology

The rough test plan is as follows (written before implementation, so TDD I guess?)
1. Create a folder with 2 target files, and 1 big file
  ```
  ❯ ls -l plugin
  total 1053376
  -rw-rw-r-- 1 edward edward 1073741824 Nov 11 11:25 big.file
  -rw-rw-r-- 1 edward edward    4903179 Nov 10 20:13 plugin.wasm
  -rw-rw-r-- 1 edward edward        512 Nov 10 14:26 plugin.yaml
  ```
2. Make a `.zip` a `.tar` a `.tar.gz`
  ```
  ❯ ls -l plugin.*
  -rw-rw-r-- 1 edward edward 1078650880 Nov 11 11:26 plugin.tar
  -rw-rw-r-- 1 edward edward 1075263074 Nov 11 11:26 plugin.tar.gz
  -rw-rw-r-- 1 edward edward 1075261541 Nov 11 11:26 plugin.zip
  ```
3. Make a small go function that tries to do the folowing on each archive type:
  - Open the archive
  - ~~Read the contents of `plugin.yaml` and parse as a string~~
  - Read the contents of `plugin.yaml` and return it as a `[]byte`
  - Read the contents of `plugin.wasm` and store in full as a `[]byte`
  - Close the archive
  - Return the `plugin.yaml` as a `string`, and `plugin.wasm` as a `[]byte`
4. Add go tests that validate the functions return as expected
5. Add go benchmarks that'll measure speed, CPU usage, and Memory usage
6. Add results to this document
7. Profit?

# Results

Below are the results of various rounds of testing, I've included the commit 
hash for the version of code that I ran to get a specific output 

1. (`8a78afc2967566526ae23146f39b140d68555c45`)
  Using `io.ReadAll` on the tar file's files
  Turns out this was using ~2.2GB of RAM, as io.ReadAll is using a temporary buffer
```
go test -run='^$' -bench=. -benchmem ./...
goos: linux
goarch: amd64
pkg: example.com/archive-fs-test
cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
BenchmarkZip-8                51          21952518 ns/op         5318240 B/op       3741 allocs/op
BenchmarkTar-8                 1        1355336519 ns/op        2202710896 B/op      129 allocs/op
BenchmarkTarGz-8               1        1330816974 ns/op        2762209160 B/op     4371 allocs/op
PASS
ok      example.com/archive-fs-test     3.946s
```

2. (`06a161d6928f11d59d5f564c8d08c3e3f6545f3e`)
  Using `io.Readfull` and pre-allocating the buffer based on the tar header
  Tar is now ~2.3x quicker, and uses half as much RAM
  Tar.gz is ~1.9x quicker, so less of an improvement but still better

```
go test -run='^$' -bench=. -benchmem ./...
goos: linux
goarch: amd64
pkg: example.com/archive-fs-test
cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
BenchmarkZip-8                50          22844869 ns/op         5318262 B/op       3741 allocs/op
BenchmarkTar-8                 2         572265646 ns/op        1078661668 B/op       56 allocs/op
BenchmarkTarGz-8               2         686186814 ns/op        1079088692 B/op     4294 allocs/op
PASS
ok      example.com/archive-fs-test     3.772s
```

3. Not really a run; but at this point I (got claude to)
  - seperate the tests into their own files (`zip.go` `tar.go` etc.)
  - Add a V1 Prefix
  - to make it so make bench runs 10 iterations
