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
(Each result below is the output of `make bench` except where stated otherwise)

1. (`8a78afc2967566526ae23146f39b140d68555c45`)
  Using `io.ReadAll` on the tar file's files
  Turns out this was using ~2.2GB of RAM, as io.ReadAll is using a temporary
  buffer (doubling memory usage)
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
```
go test -run='^$' -bench='.' -benchtime=10x -benchmem ./...
goos: linux
goarch: amd64
pkg: example.com/archive-fs-test
cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
BenchmarkTar/V1-8             10         190695317 ns/op        1078656828 B/op       54 allocs/op
BenchmarkTarGz/V1-8                   10         725309938 ns/op        1079088188 B/op     4292 allocs/op
BenchmarkZip/V1-8                     10          22352958 ns/op         5321957 B/op       3748 allocs/op
PASS
ok      example.com/archive-fs-test     9.479s
```

4. I had an idea while grabbing a bite to eat, if we just search the tar for a
   target file and skip the rest, we can save on loading anything into memory
   other than the target file.
   This does come at the cost of requiring a new file open per read, however.
   (also at this point, I renamed the tar.gz test to TGZ as the alignment was
   annoying)

   Anyway, that made 
```
go test -run='^$' -bench='.' -benchtime=10x -benchmem ./...
goos: linux
goarch: amd64
pkg: example.com/archive-fs-test
cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
BenchmarkTar/V1-8                     10         267320647 ns/op        1078659004 B/op       55 allocs/op
BenchmarkTar/V2-8                     10            798046 ns/op         4915784 B/op         80 allocs/op
BenchmarkTGz/V1-8                     10         768202258 ns/op        1079088184 B/op     4292 allocs/op
BenchmarkTGz/V2-8                     10         588651060 ns/op         5773452 B/op       8555 allocs/op
BenchmarkZip/V1-8                     10          22109717 ns/op         5321936 B/op       3748 allocs/op
PASS
ok      example.com/archive-fs-test     16.479s
```

5. To make things fair, I also added a V2 of the ZIP implementation that tires
   pre-allocating memory instead of using `io.readFull`
   This has no meaningfull difference
```
go test -run='^$' -bench='.' -benchtime=10x -benchmem ./...
goos: linux
goarch: amd64
pkg: example.com/archive-fs-test
cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
BenchmarkTar/V1-8                     10         468493807 ns/op        1078656817 B/op       54 allocs/op
BenchmarkTar/V2-8                     10            642817 ns/op         4914119 B/op         80 allocs/op
BenchmarkTGz/V1-8                     10         756886218 ns/op        1079089005 B/op     4292 allocs/op
BenchmarkTGz/V2-8                     10         596504599 ns/op         5772641 B/op       8555 allocs/op
BenchmarkZip/V1-8                     10          25153705 ns/op         5321920 B/op       3748 allocs/op
BenchmarkZip/V2-8                     10          24032298 ns/op         5406277 B/op       3751 allocs/op
PASS
ok      example.com/archive-fs-test     18.730s
```

6. Finally, to get some decent data for later comparisons, I ran
  `make bench BENCH=/V2 BENCHTIME=50x` to run 50 iterations of only the V2
  implementations
```
go test -run='^$' -bench='/V2' -benchtime=50x -benchmem ./...
goos: linux
goarch: amd64
pkg: example.com/archive-fs-test
cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz
BenchmarkTar/V2-8                     50            778928 ns/op         4913254 B/op         79 allocs/op
BenchmarkTGz/V2-8                     50         547552835 ns/op         5771328 B/op       8555 allocs/op
BenchmarkZip/V2-8                     50          22662637 ns/op         5406278 B/op       3752 allocs/op
PASS
ok      example.com/archive-fs-test     28.554s
```
