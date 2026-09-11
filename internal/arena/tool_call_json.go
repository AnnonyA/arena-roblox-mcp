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
	for _, r := range decoded.Function.Name {
		if unicode.IsControl(r) {
			return fmt.Errorf("tool call name contains control character")
		}
	}
	*call = ToolCall(decoded)
	return nil
}
