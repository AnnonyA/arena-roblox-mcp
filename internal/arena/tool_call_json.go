package arena

import (
	"encoding/json"
	"fmt"
	"unicode"
)

func (call *ToolCall) UnmarshalJSON(data []byte) error {
	type toolCallAlias ToolCall

	var decoded toolCallAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	for _, r := range decoded.ID {
		if unicode.IsControl(r) {
			return fmt.Errorf("tool call id contains control character")
		}
		if isBidirectionalFormatting(r) {
			return fmt.Errorf("tool call id contains bidirectional formatting")
		}
	}
	for _, r := range decoded.Function.Name {
		if unicode.IsControl(r) {
			return fmt.Errorf("tool call name contains control character")
		}
		if isBidirectionalFormatting(r) {
			return fmt.Errorf("tool call name contains bidirectional formatting")
		}
	}
	*call = ToolCall(decoded)
	return nil
}
