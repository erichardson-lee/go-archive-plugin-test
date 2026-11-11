# Archive Filesystem Test

This is a quick test repo to see how different archive formats work with go

See [TEST_PLAN.md](./TEST_PLAN.md) for information on methodology / results

Main thing to know is that I wanted to see how different methods compare, and if

# Results

These results are based on 50 iterations, all values are per operation; and for
the sake of sanity, each value is rounded to 4 signifigant figures <br/>
(Full values available in [TEST_PLAN.md](./TEST_PLAN.md#results))

This has all been ran with a `big.file`™ in the archive

| Method      | Runtime (ms) | Memory Usage (kb) | Allocs |
| ----------- | ------------ | ----------------- | ------ |
| Tar V2      | 0.7889       | 4913              | 79     |
| Zip V2      | 22.66        | 5406              | 3752   |
| Tar Gzip V2 | 547.76       | 5771              | 8555   |

From this, we can see that Tar beats Zip and Tar Gzip in terms of speed, memory,
and allocs

When running without the big file, we see that Tar Gzip catches up to Zip

| Method      | Runtime (ms) | Memory Usage (kb) | Allocs |
| ----------- | ------------ | ----------------- | ------ |
| Tar V2      | 0.5530       | 4912              | 53     |
| Zip V2      | 21.89        | 5406              | 3747   |
| Tar Gzip V1 | 21.97        | 5363              | 3759   |

With the big file included, the output file sizes are essentially the same, with
zip being the smallest by 3mb

However, when not including the big file, the archive file sizes are:

| Archive Method | Size (mb) | Compared to Zip |
| -------------- | --------- | --------------- |
| Tar            | 4.915     | 3.651x          |
| Zip            | 1.346     | 1x              |
| Tar Gzip       | 1.346     | 1x              |

(note: tar.gz is larger than a zip by 102 bytes)

I also tested to see what the performance was like when we were loading the big
file as the WASM module

| Method (n=10)   | Runtime (ms) | Memory Usage (kb) | Allocs |
| --------------- | ------------ | ----------------- | ------ |
| Tar V2 Big      | 482.6        | 1073750           | 68     |
| Tar Gzip V2 Big | 601.0        | 1074235           | 4836   |
| Zip V2 Big      | 659.0        | 1073849           | 589    |

