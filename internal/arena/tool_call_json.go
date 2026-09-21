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

func validateToolCallField(field, value string) error {
	if strings.TrimSpace(value) != value {
		return fmt.Errorf("tool call %s contains surrounding whitespace", field)
	}
	if len(value) > maxToolCallIdentifierBytes {
		return fmt.Errorf("tool call %s exceeds %d bytes", field, maxToolCallIdentifierBytes)
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return fmt.Errorf("tool call %s contains control character", field)
		}
		if r == '\u2028' || r == '\u2029' {
			return fmt.Errorf("tool call %s contains line separator", field)
		}
		if unicode.Is(unicode.Cf, r) {
			return fmt.Errorf("tool call %s contains format character", field)
		}
	}
	return nil
}

func (call *ToolCall) UnmarshalJSON(data []byte) error {
	if !utf8.Valid(data) {
		return fmt.Errorf("tool call contains invalid UTF-8")
	}

	type toolCallAlias ToolCall

	var decoded toolCallAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	if err := validateToolCallField("id", decoded.ID); err != nil {
		return err
	}
	if err := validateToolCallField("type", decoded.Type); err != nil {
		return err
	}
	if err := validateToolCallField("name", decoded.Function.Name); err != nil {
		return err
	}
	if len(decoded.Function.Arguments) > maxToolCallArgumentsBytes {
		return fmt.Errorf("tool call arguments exceed %d bytes", maxToolCallArgumentsBytes)
	}
	*call = ToolCall(decoded)
	return nil
}
