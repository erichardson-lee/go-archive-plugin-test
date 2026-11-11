package main

import "testing"

var tarGzImpls = []impl{
	{"V1", TarGzProcessorV1},
}

func TestTarGz(t *testing.T) { testImpls(t, "plugin.tar.gz", tarGzImpls) }

func BenchmarkTarGz(b *testing.B) { benchImpls(b, "plugin.tar.gz", tarGzImpls) }
