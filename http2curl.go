package http2curl

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// CurlCommand contains exec.Command compatible slice + helpers
type CurlCommand struct {
	slice []string
}

// append appends a string to the CurlCommand
func (c *CurlCommand) append(newSlice ...string) {
	c.slice = append(c.slice, newSlice...)
}

// String returns a ready to copy/paste command
func (c *CurlCommand) String() string {
	slice := make([]string, len(c.slice))
	copy(slice, c.slice)
	for i := range slice {
		quoted := fmt.Sprintf("%q", slice[i])
		if len(quoted) != len(slice[i])+2 {
			slice[i] = quoted
		}
	}
	return strings.Join(slice, " ")
}

// GetCurlCommand returns a CurlCommand corresponding to an http.Request
func GetCurlCommand(req *http.Request) (*CurlCommand, error) {
	command := CurlCommand{}

	command.append("curl")

	command.append("-X", req.Method)

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	if buf.Len() > 0 {
		command.append("-d", buf.String())
	}

	for key, values := range req.Header {
		command.append("-H", key, strings.Join(values, " "))
	}

	command.append(req.URL.String())

	return &command, nil
}
