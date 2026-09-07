package main

import (
	"context"
	"io"
	"strings"
	"testing"
)

func BenchmarkCLIStartup(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		in := strings.NewReader("/exit\n")
		if err := runWithStudioDependencies(ctx, in, io.Discard, []string{"--model", "arena/benchmark"}, nil, nil, nil, nil); err != nil {
			b.Fatal(err)
		}
	}
}
