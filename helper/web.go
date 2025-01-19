package helper

import (
	"bytes"
	"io"
	"net/http"
)

func ReadRequestBody(request *http.Request) ([]byte, error) {
	body := request.Body

	// Read the body
	var buf bytes.Buffer
	if body != nil {
		_, err := io.Copy(&buf, body)
		if err != nil {
			return []byte{}, err
		}
	}
	// Reset the body so it can be read again
	request.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	return buf.Bytes(), nil
}
