package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// CanonicalTemplateData normalizes object-shaped template write input.
// It ignores whitespace and object key order and preserves number lexemes.
func CanonicalTemplateData(raw json.RawMessage) (json.RawMessage, error) {
	value, err := decodeJSONValue(raw)
	if err != nil {
		return nil, err
	}

	return marshalTemplateObject(value)
}

func canonicalTemplateReadback(raw json.RawMessage) (json.RawMessage, error) {
	value, err := decodeTemplateEnvelope(raw)
	if err != nil {
		return nil, err
	}

	return marshalTemplateObject(value)
}

// decodeTemplateEnvelope reads Linear templateData as a JSON value, and if that
// value is a non-empty string, as one nested JSON value. It does not require
// an object and does not apply emptiness policy.
func decodeTemplateEnvelope(raw json.RawMessage) (any, error) {
	value, err := decodeJSONValue(raw)
	if err != nil {
		return nil, err
	}
	encoded, ok := value.(string)
	if !ok || encoded == "" {
		return value, nil
	}
	unwrapped, err := decodeJSONValue(json.RawMessage(encoded))
	if err != nil {
		return nil, fmt.Errorf("%w: template data must be a JSON object", ErrWriteInvalid)
	}

	return unwrapped, nil
}

func mustJSON(value any) json.RawMessage {
	canonical, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return canonical
}

func decodeJSONValue(raw json.RawMessage) (any, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, fmt.Errorf("%w: template data must be a JSON object", ErrWriteInvalid)
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("%w: malformed template JSON", ErrWriteInvalid)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%w: malformed template JSON", ErrWriteInvalid)
	}

	return value, nil
}

func marshalTemplateObject(value any) (json.RawMessage, error) {
	object, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: template data must be a JSON object", ErrWriteInvalid)
	}

	return mustJSON(object), nil
}
