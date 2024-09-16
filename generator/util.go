package generator

import (
	"bytes"
	"go-postgres-generator-example/logger"
	"os"
	"strings"
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

func firstLetterToLower(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}

func firstLetterToUpper(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}

func hasPrismaReference(field ParsedStructField) bool {
	return getPrismaReference(field) != ""
}
