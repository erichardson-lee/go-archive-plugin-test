package main

import "testing"

var zipImpls = []impl{
	{"V1", ZipProcessorV1},
	{"V2", ZipProcessorV2},
}

func TestZip(t *testing.T) { testImpls(t, "plugin.zip", zipImpls) }

func BenchmarkZip(b *testing.B) { benchImpls(b, "plugin.zip", zipImpls) }
