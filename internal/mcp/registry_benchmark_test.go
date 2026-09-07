package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func BenchmarkToolRegistry(b *testing.B) {
	tools := make([]Tool, 64)
	for i := range tools {
		tools[i] = Tool{
			Name:        fmt.Sprintf("tool_%02d", i),
			Description: "benchmark tool",
			InputSchema: json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"}}}`),
		}
	}
	registry := NewRegistry(func(context.Context) ([]Tool, error) {
		return tools, nil
	})
	ctx := context.Background()
	if _, err := registry.Tools(ctx); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := registry.Tools(ctx); err != nil {
			b.Fatal(err)
		}
	}
}
