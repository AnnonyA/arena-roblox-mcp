package mcp

import (
	"context"
	"errors"
	"testing"
)

func TestSDKSessionRejectsBlankToolNameBeforeCallingSDK(t *testing.T) {
	session := &sdkSession{}

	_, err := session.CallTool(context.Background(), " \t\n", nil)
	if !errors.Is(err, ErrMissingToolName) {
		t.Fatalf("CallTool error = %v, want ErrMissingToolName", err)
	}
}
