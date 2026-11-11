#!/usr/bin/env bash
# Packages the contents of plugin/ into plugin.zip, plugin.tar and plugin.tar.gz.
# Files sit at the top level of each archive (plugin.yaml, not plugin/plugin.yaml).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$SCRIPT_DIR/plugin"
OUT_DIR="$SCRIPT_DIR"

if [[ ! -d "$SRC_DIR" ]]; then
  echo "error: $SRC_DIR does not exist" >&2
  exit 1
fi

cd "$SRC_DIR"

# Collect top-level entries (including dotfiles) so archive paths have no ./ prefix.
shopt -s nullglob dotglob
entries=(*)
shopt -u nullglob dotglob

if [[ ${#entries[@]} -eq 0 ]]; then
  echo "error: $SRC_DIR is empty" >&2
  exit 1
fi

# Remove old archives so we don't append to / update stale ones.
rm -f "$OUT_DIR/plugin.zip" "$OUT_DIR/plugin.tar" "$OUT_DIR/plugin.tar.gz"

echo "Creating plugin.zip ..."
zip -r -q "$OUT_DIR/plugin.zip" "${entries[@]}"

echo "Creating plugin.tar ..."
tar -cf "$OUT_DIR/plugin.tar" "${entries[@]}"

echo "Creating plugin.tar.gz ..."
tar -czf "$OUT_DIR/plugin.tar.gz" "${entries[@]}"

echo "Done:"
ls -lh "$OUT_DIR/plugin.zip" "$OUT_DIR/plugin.tar" "$OUT_DIR/plugin.tar.gz"
