package http

import (
	"encoding/json"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/pkg/errors"
)

// SignatureType is the type of signature
type SignatureType int

// Supported values for the SignatureType
const (
	AWSSignature SignatureType = iota + 1
	limit
)

// signatureTypeMappings is a table for conversion of signature type from string.
var signatureTypeMappings = map[SignatureType]string{
	AWSSignature: "aws",
}

func GetSupportedSignaturesTypes() []SignatureType {
	var result []SignatureType
	for index := SignatureType(1); index < limit; index++ {
		result = append(result, index)
	}
	return result
}

func toSignatureType(valueToMap string) (SignatureType, error) {
	normalizedValue := normalizeValue(valueToMap)
	for key, currentValue := range signatureTypeMappings {
		if normalizedValue == currentValue {
			return key, nil
		}
	}
	return -1, errors.New("invalid signature type: " + valueToMap)
}

func normalizeValue(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func (t SignatureType) String() string {
	return signatureTypeMappings[t]
}

// SignatureTypeHolder is used to hold internal type of the signature
type SignatureTypeHolder struct {
	Value SignatureType
}

func (holder SignatureTypeHolder) JSONSchemaType() *jsonschema.Type {
	gotType := &jsonschema.Type{
		Type:        "string",
		Title:       "type of the signature",
		Description: "Type of the signature",
	}
	for _, types := range GetSupportedSignaturesTypes() {
		gotType.Enum = append(gotType.Enum, types.String())
	}
	return gotType
}

func (holder *SignatureTypeHolder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var marshalledTypes string
	if err := unmarshal(&marshalledTypes); err != nil {
		return err
	}

	computedType, err := toSignatureType(marshalledTypes)
	if err != nil {
		return err
	}

	holder.Value = computedType
	return nil
}

func (holder *SignatureTypeHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(holder.Value.String())
}

func (holder SignatureTypeHolder) MarshalYAML() (interface{}, error) {
	return holder.Value.String(), nil
}