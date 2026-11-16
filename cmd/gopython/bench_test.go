package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkPrograms(b *testing.B) {
	paths, err := filepath.Glob("../../benchmarks/programs/*.py")
	if err != nil || len(paths) == 0 {
		b.Fatalf("no benchmark programs found: %v", err)
	}
	for _, path := range paths {
		source, err := os.ReadFile(path)
		if err != nil {
			b.Fatal(err)
		}
		name := strings.TrimSuffix(filepath.Base(path), ".py")
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				if err := runSource(string(source), io.Discard); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
