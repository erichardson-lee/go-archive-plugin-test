#!/usr/bin/env bash
# Generates a file of random data at plugin/big.file.
# Usage: ./1.generate-big-file.sh [size_in_mb]   (default: 1024 MB = 1 GiB)
set -euo pipefail

SIZE_MB="${1:-1024}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUT_DIR="$SCRIPT_DIR/plugin"
OUT_FILE="$OUT_DIR/big.file"

mkdir -p "$OUT_DIR"

echo "Writing ${SIZE_MB} MB of random data to $OUT_FILE ..."
dd if=/dev/urandom of="$OUT_FILE" bs=1M count="$SIZE_MB" iflag=fullblock status=progress

echo "Done:"
ls -lh "$OUT_FILE"
