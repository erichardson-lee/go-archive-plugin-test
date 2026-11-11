.PHONY: test bench bench-raw clean

# bash + pipefail so a failing `go test` still fails `make bench`.
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -c

ARCHIVES := plugin.zip plugin.tar plugin.tar.gz

# Iterations per benchmark (exact count, via -benchtime=Nx). Override with
# e.g. `make bench BENCHTIME=50x` or `make bench BENCHTIME=5s`.
BENCHTIME ?= 10x
# Filter which benchmarks run, e.g. `make bench BENCH=Tar` or `BENCH=/V2`.
BENCH ?= .

test: $(ARCHIVES)
	go test -v ./...

BENCH_CMD = go test -run='^$$' -bench='$(BENCH)' -benchtime=$(BENCHTIME) -benchmem ./...

# Results as a markdown table (see 3.format-bench.sh).
bench: $(ARCHIVES)
	$(BENCH_CMD) | ./3.format-bench.sh

# Unformatted `go test` output, e.g. to see failure details.
bench-raw: $(ARCHIVES)
	$(BENCH_CMD)

# Only built if missing (no timestamp-based rebuilds).
plugin/big.file:
	./1.generate-big-file.sh

# Grouped target: one script run makes all three. Order-only dep on big.file
# means archives are only built if missing, not when big.file is newer.
$(ARCHIVES) &: | plugin/big.file
	./2.make-archives.sh

clean:
	rm -f plugin/big.file $(ARCHIVES)
