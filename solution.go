package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func SearchString(input *os.File) (string, error) {
	// Read the file contents
	content, err := io.ReadAll(input)
	if err != nil {
		return "", err
	}

	dataLen := len(content)

	if dataLen == 0 {
		return "TYPE:EMPTY|LEN:0|DATA:", nil
	}

	// Check for PNG
	if dataLen >= 4 && content[0] == 0x89 && content[1] == 0x50 && content[2] == 0x4E && content[3] == 0x47 {
		encoded := base64.StdEncoding.EncodeToString(content)
		return fmt.Sprintf("TYPE:PNG|LEN:%d|DATA:%s", dataLen, encoded), nil
	}

	// Check for XML
	if dataLen > 0 && content[0] == '<' {
		return fmt.Sprintf("TYPE:XML|LEN:%d|DATA:%s", dataLen, string(content)), nil
	}

	// Check for JSON
	if dataLen > 0 && (content[0] == '{' || content[0] == '[') {
		return fmt.Sprintf("TYPE:JSON|LEN:%d|DATA:%s", dataLen, string(content)), nil
	}

	// Check if it's printable text
	isPrintable := true
	for i := 0; i < dataLen; i++ {
		if content[i] < 32 && content[i] != '\n' && content[i] != '\r' && content[i] != '\t' {
			isPrintable = false
			break
		}
	}

	if isPrintable {
		return fmt.Sprintf("TYPE:TEXT|LEN:%d|DATA:%s", dataLen, string(content)), nil
	}

	// Binary data - base64 encode it
	encoded := base64.StdEncoding.EncodeToString(content)
	return fmt.Sprintf("TYPE:BINARY|LEN:%d|DATA:%s", dataLen, encoded), nil
}
