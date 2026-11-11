package main

import "testing"

var tarImpls = []impl{
	{"V1", TarProcessorV1},
	{"V2", TarProcessorV2},
	{"V2/Big", TarProcessorV2Big},
}

func TestTar(t *testing.T) { testImpls(t, "plugin.tar", tarImpls) }

func BenchmarkTar(b *testing.B) { benchImpls(b, "plugin.tar", tarImpls) }
