package arena

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (chunk *streamChunk) UnmarshalJSON(data []byte) error {
	if err := rejectDuplicateStreamJSONKeys(data); err != nil {
		return err
	}

	type streamChunkAlias streamChunk
	return json.Unmarshal(data, (*streamChunkAlias)(chunk))
}

func rejectDuplicateStreamJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	return scanStreamJSONValue(decoder)
}

func scanStreamJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}

	switch delim {
	case '{':
		seen := make(map[string]struct{})
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("invalid JSON object key")
			}
			if _, exists := seen[key]; exists {
				return fmt.Errorf("duplicate JSON key %q", key)
			}
			seen[key] = struct{}{}
			if err := scanStreamJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	case '[':
		for decoder.More() {
			if err := scanStreamJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	default:
		return nil
	}
}
