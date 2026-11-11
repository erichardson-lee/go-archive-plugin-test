.PHONY: test bench clean

ARCHIVES := plugin.zip plugin.tar plugin.tar.gz

test: $(ARCHIVES)
	go test -v ./...

bench: $(ARCHIVES)
	go test -run='^$$' -bench=. -benchmem ./...

# Only built if missing (no timestamp-based rebuilds).
plugin/big.file:
	./1.generate-big-file.sh

# Grouped target: one script run makes all three. Order-only dep on big.file
# means archives are only built if missing, not when big.file is newer.
$(ARCHIVES) &: | plugin/big.file
	./2.make-archives.sh

clean:
	rm -f plugin/big.file $(ARCHIVES)
