package generator

import (
	"bytes"
	"go/format"
)

func FormatResult(buffer *bytes.Buffer) (*[]byte, error) {
	formattedBytes, err := format.Source(buffer.Bytes())
	if err != nil {
		return &[]byte{}, err
	}
	return &formattedBytes, nil
}
