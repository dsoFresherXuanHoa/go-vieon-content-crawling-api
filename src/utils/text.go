package utils

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrWriteString2File = errors.New("write string to text/plain file failure")
	ErrReadString2File  = errors.New("read string from text/plain file failure")
)

type textUtil struct{}

func NewTextUtil() *textUtil {
	return &textUtil{}
}

func (textUtil) String2File(content string, filePath string) error {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		fmt.Println("Error while write string context to text/plain file: " + err.Error())
		return ErrWriteString2File
	}
	return nil
}

func (textUtil) File2String(filePath string) (*string, error) {
	if content, err := os.ReadFile(filePath); err != nil {
		fmt.Println("Error while read string context from text/plain file: " + err.Error())
		return nil, ErrReadString2File
	} else {
		result := string(content)
		return &result, nil
	}
}
