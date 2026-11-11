# Methodology

The rough test plan is as follows (written before implementation, so TDD I guess?)

1. Create a folder with 2 target files, and 1 big file

```sh
❯ ls -l plugin
total 1053376
-rw-rw-r-- 1 edward edward 1073741824 Nov 11 11:25 big.file
-rw-rw-r-- 1 edward edward    4903179 Nov 10 20:13 plugin.wasm
-rw-rw-r-- 1 edward edward        512 Nov 10 14:26 plugin.yaml
```

2. Make a `.zip` a `.tar` a `.tar.gz` of that plugin

```sh
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

> When this file got to ~200 lines I realised these outputs are hard to scan
> through. Through the power of claude, I made `3.formatBench.sh` that transforms
> from go bench output to a markdown table, if you want to see the original raw
> results, look at (`748701b144e6f95b8548a4a62e058954beec419f`)

1. (`8a78afc2967566526ae23146f39b140d68555c45`)
   Using `io.ReadAll` on the tar file's files
   Turns out this was using ~2.2GB of RAM, as io.ReadAll is using a temporary
   buffer (doubling memory usage)

| Method         | Runtime (ms) | Memory Usage (kb) | Allocs |
| -------------- | ------------ | ----------------- | ------ |
| Zip (n=51)     | 21.95        | 5318              | 3741   |
| Tar (n=1)      | 1355         | 2202711           | 129    |
| Tar Gzip (n=1) | 1331         | 2762209           | 4371   |

2. (`06a161d6928f11d59d5f564c8d08c3e3f6545f3e`)
   Using `io.Readfull` and pre-allocating the buffer based on the tar header
   Tar is now ~2.3x quicker, and uses half as much RAM
   Tar.gz is ~1.9x quicker, so less of an improvement but still better

| Method         | Runtime (ms) | Memory Usage (kb) | Allocs |
| -------------- | ------------ | ----------------- | ------ |
| Zip (n=50)     | 22.84        | 5318              | 3741   |
| Tar (n=2)      | 572.3        | 1078662           | 56     |
| Tar Gzip (n=2) | 686.2        | 1079089           | 4294   |

3. Not really a run; but at this point I (got claude to)

- seperate the tests into their own files (`zip.go` `tar.go` etc.)
- Add a V1 Prefix
- to make it so make bench runs 10 iterations

| Method (n=10) | Runtime (ms) | Memory Usage (kb) | Allocs |
| ------------- | ------------ | ----------------- | ------ |
| Tar V1        | 190.7        | 1078657           | 54     |
| Tar Gzip V1   | 725.3        | 1079088           | 4292   |
| Zip V1        | 22.35        | 5322              | 3748   |

4. (`d3963470dd9f746828f22d516fa8ad10796370c5`)
   I had an idea while grabbing a bite to eat, if we just search the tar for a
   target file and skip the rest, we can save on loading anything into memory
   other than the target file.
   This does come at the cost of requiring a new file open per read, however.
   (also at this point, I renamed the tar.gz test to TGZ as the alignment was
   annoying)

Anyway, that made

| Method (n=10) | Runtime (ms) | Memory Usage (kb) | Allocs |
| ------------- | ------------ | ----------------- | ------ |
| Tar V1        | 267.3        | 1078659           | 55     |
| Tar V2        | 0.7980       | 4916              | 80     |
| Tar Gzip V1   | 768.2        | 1079088           | 4292   |
| Tar Gzip V2   | 588.7        | 5773              | 8555   |
| Zip V1        | 22.11        | 5322              | 3748   |

5. To make things fair, I also added a V2 of the ZIP implementation that tires
   pre-allocating memory instead of using `io.readFull`
   This has no meaningfull difference

| Method (n=10) | Runtime (ms) | Memory Usage (kb) | Allocs |
| ------------- | ------------ | ----------------- | ------ |
| Tar V1        | 468.5        | 1078657           | 54     |
| Tar V2        | 0.6428       | 4914              | 80     |
| Tar Gzip V1   | 756.9        | 1079089           | 4292   |
| Tar Gzip V2   | 596.5        | 5773              | 8555   |
| Zip V1        | 25.15        | 5322              | 3748   |
| Zip V2        | 24.03        | 5406              | 3751   |

6. Finally, to get some decent data for later comparisons, I ran
   `make bench BENCH=/V2 BENCHTIME=50x` to run 50 iterations of only the V2
   implementations

| Method (n=50) | Runtime (ms) | Memory Usage (kb) | Allocs |
| ------------- | ------------ | ----------------- | ------ |
| Tar V2        | 0.7789       | 4913              | 79     |
| Tar Gzip V2   | 547.6        | 5771              | 8555   |
| Zip V2        | 22.66        | 5406              | 3752   |

7. As an extra bit, I re-ran the 10x, but without a `big.file` in the archives,
   and got:

| Method (n=10) | Runtime (ms) | Memory Usage (kb) | Allocs |
| ------------- | ------------ | ----------------- | ------ |
| Tar V1        | 0.9842       | 4913              | 40     |
| Tar V2        | 0.7403       | 4916              | 54     |
| Tar Gzip V1   | 22.53        | 5367              | 3760   |
| Tar Gzip V2   | 43.13        | 5812              | 7493   |
| Zip V1        | 21.78        | 5322              | 3744   |
| Zip V2        | 21.48        | 5406              | 3747   |

since it was so much quicker, I re-ran @ 50x

| Method (n=50) | Runtime (ms) | Memory Usage (kb) | Allocs |
| ------------- | ------------ | ----------------- | ------ |
| Tar V1        | 0.7114       | 4911              | 38     |
| Tar V2        | 0.5530       | 4912              | 53     |
| Tar Gzip V1   | 21.97        | 5363              | 3759   |
| Tar Gzip V2   | 43.00        | 5811              | 7492   |
| Zip V1        | 21.50        | 5318              | 3736   |
| Zip V2        | 21.89        | 5406              | 3747   |

on the smaller files, it seems like these are the fastest options:

- BenchmarkTar/V2-8
- BenchmarkTGz/V1-8
- BenchmarkZip/V2-8 <- barely, same either way tbh

8. I figured I would see what it's like when the big file is the target file
   (and not the ignored one)
   `make bench BENCH=/V2/Big`

| Method (n=10)   | Runtime (ms) | Memory Usage (kb) | Allocs |
| --------------- | ------------ | ----------------- | ------ |
| Tar V2 Big      | 482.6        | 1073750           | 68     |
| Tar Gzip V2 Big | 601.0        | 1074235           | 4836   |
| Zip V2 Big      | 659.0        | 1073849           | 589    |
