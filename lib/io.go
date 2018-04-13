package ferryman

import (
	"bytes"
	"io"
	"strings"
)

// Replacement buffer copy possibly inefficient
func replaceBuffer(dst io.Writer, src io.Reader, find string, replace string) (written int64, err error) {
	var bBuff bytes.Buffer
	var chunk []byte = make([]byte, 32*1024)

	for {
		read, readErr := io.ReadFull(src, chunk)
		if readErr == nil {
			bBuff.Write(chunk)
		} else {
			if strings.Contains(readErr.Error(), "EOF") {
				bBuff.Write(chunk[:read])
			} else {
				err = readErr
			}
			break
		}
	}

	if err == nil {
		replaced := strings.Replace(bBuff.String(), find, replace, -1)
		wrote, wErr := dst.Write([]byte(replaced))
		if wErr == nil {
			return int64(wrote), wErr
		} else {
			return 0, wErr
		}
	} else {
		return 0, err
	}
}
