package arena

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxToolCallIdentifierBytes = 4096
	maxToolCallArgumentsBytes  = 1024 * 1024
)

func (call *ToolCall) UnmarshalJSON(data []byte) error {
	if !utf8.Valid(data) {
		return fmt.Errorf("tool call contains invalid UTF-8")
	}

	type toolCallAlias ToolCall

	var decoded toolCallAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	if strings.TrimSpace(decoded.ID) != decoded.ID {
		return fmt.Errorf("tool call id contains surrounding whitespace")
	}
	if len(decoded.ID) > maxToolCallIdentifierBytes {
		return fmt.Errorf("tool call id exceeds %d bytes", maxToolCallIdentifierBytes)
	}
	for _, r := range decoded.ID {
		if unicode.IsControl(r) {
			return fmt.Errorf("tool call id contains control character")
		}
		if r == '\u2028' || r == '\u2029' {
			return fmt.Errorf("tool call id contains line separator")
		}
		if unicode.Is(unicode.Cf, r) {
			return fmt.Errorf("tool call id contains format character")
		}
	}
	if strings.TrimSpace(decoded.Type) != decoded.Type {
		return fmt.Errorf("tool call type contains surrounding whitespace")
	}
	if len(decoded.Type) > maxToolCallIdentifierBytes {
		return fmt.Errorf("tool call type exceeds %d bytes", maxToolCallIdentifierBytes)
	}
	for _, r := range decoded.Type {
		if unicode.IsControl(r) {
			return fmt.Errorf("tool call type contains control character")
		}
		if r == '\u2028' || r == '\u2029' {
			return fmt.Errorf("tool call type contains line separator")
		}
		if unicode.Is(unicode.Cf, r) {
			return fmt.Errorf("tool call type contains format character")
		}
	}
	if strings.TrimSpace(decoded.Function.Name) != decoded.Function.Name {
		return fmt.Errorf("tool call name contains surrounding whitespace")
	}
	if len(decoded.Function.Name) > maxToolCallIdentifierBytes {
		return fmt.Errorf("tool call name exceeds %d bytes", maxToolCallIdentifierBytes)
	}
	for _, r := range decoded.Function.Name {
		if unicode.IsControl(r) {
			return fmt.Errorf("tool call name contains control character")
		}
		if r == '\u2028' || r == '\u2029' {
			return fmt.Errorf("tool call name contains line separator")
		}
		if unicode.Is(unicode.Cf, r) {
			return fmt.Errorf("tool call name contains format character")
		}
	}
	if len(decoded.Function.Arguments) > maxToolCallArgumentsBytes {
		return fmt.Errorf("tool call arguments exceed %d bytes", maxToolCallArgumentsBytes)
	}
	*call = ToolCall(decoded)
	return nil
}
