#!/usr/bin/env bash
# Formats `go test -bench ... -benchmem` output as a markdown table.
# Usage:
#   go test -run='^$' -bench=. -benchtime=50x -benchmem ./... | ./3.format-bench.sh
#   ./3.format-bench.sh bench-output.txt
set -euo pipefail

awk '
function sf4(x,   d) {
  # 4 significant figures, keeping trailing zeros (43.003 -> 43.00)
  if (x == 0) return "0.000"
  x = sprintf("%.3e", x) + 0   # round to 4sf first so 0.099999 -> 0.1000
  d = 3 - int(log(x) / log(10) + (x < 1 ? -1 : 0) + 1e-9)
  if (d < 0) d = 0
  return sprintf("%." d "f", x)
}
BEGIN {
  names["Tar"] = "Tar"; names["TGz"] = "Tar Gzip"; names["TarGz"] = "Tar Gzip"; names["Zip"] = "Zip"
}
/^Benchmark/ && /ns\/op/ {
  # e.g. BenchmarkTar/V1-8   50   711418 ns/op   4910708 B/op   38 allocs/op
  name = $1
  sub(/^Benchmark/, "", name)
  sub(/-[0-9]+$/, "", name)            # strip GOMAXPROCS suffix
  np = split(name, parts, "/")
  method = (parts[1] in names) ? names[parts[1]] : parts[1]
  for (j = 2; j <= np; j++) method = method " " parts[j]

  n = $2; ms = ""; kb = ""; allocs = ""
  for (i = 3; i < NF; i++) {
    if ($(i+1) == "ns/op")     ms = sf4($i / 1e6)
    if ($(i+1) == "B/op")      kb = sprintf("%d", $i / 1000 + 0.5)
    if ($(i+1) == "allocs/op") allocs = $i
  }
  if (count == 0) first_n = n
  else if (n != first_n) mixed_n = 1
  meth[++count] = method; ns[count] = n
  rest[count] = sprintf(" | %s | %s | %s |", ms, kb, allocs)
}
END {
  if (count == 0) { print "error: no benchmark lines found" > "/dev/stderr"; exit 1 }
  # iteration count goes in the header, unless rows differ (then per row)
  print "| Method" (mixed_n ? "" : " (n=" first_n ")") " | Runtime (ms) | Memory Usage (kb) | Allocs |"
  print "| ------ | ------------ | ----------------- | ------ |"
  for (i = 1; i <= count; i++)
    print "| " meth[i] (mixed_n ? " (n=" ns[i] ")" : "") rest[i]
}
' "${1:-/dev/stdin}" | column -t -s '|' -o '|' | sed -E '2s/ /-/g'
