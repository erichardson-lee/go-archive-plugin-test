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
  - Read the contents of `plugin.yaml` and parse as a string
  - Read the contents of `plugin.wasm` and store in full as a `[]byte`
  - Close the archive
  - Return the `plugin.yaml` as a `string`, and `plugin.wasm` as a `[]byte`
4. Add go tests that validate the functions return as expected
5. Add go benchmarks that'll measure speed, CPU usage, and Memory usage
6. Add results to this document
7. Profit?
