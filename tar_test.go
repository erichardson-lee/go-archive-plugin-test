package main

import "testing"

var tarImpls = []impl{
	{"V1", TarProcessorV1},
}

func TestTar(t *testing.T) { testImpls(t, "plugin.tar", tarImpls) }

func BenchmarkTar(b *testing.B) { benchImpls(b, "plugin.tar", tarImpls) }
