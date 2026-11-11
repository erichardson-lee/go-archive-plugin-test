package main

import "testing"

var tarGzImpls = []impl{
	{"V1", TarGzProcessorV1},
	{"V2", TarGzProcessorV2},
	{"V2/Big", TarGzProcessorV2Big},
}

func TestTGz(t *testing.T) { testImpls(t, "plugin.tar.gz", tarGzImpls) }

func BenchmarkTGz(b *testing.B) { benchImpls(b, "plugin.tar.gz", tarGzImpls) }
