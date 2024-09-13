package generator

import (
	"bytes"
	"go-postgres-generator-example/logger"
	"os"
)

func WriteBufferToFile(buffer *bytes.Buffer, path string) error {
	bytesArray := buffer.Bytes()
	return WriteBytesToFile(&bytesArray, path)
}

func WriteBytesToFile(bytes *[]byte, path string) error {
	outputFile, err := os.Create(path)
	if err != nil {
		logger.Error("failed to create file", err)
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(*bytes)
	if err != nil {
		logger.Error("failed to write file", err)
		return err
	}
	return nil
}
