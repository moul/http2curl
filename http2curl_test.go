package http2curl

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func ExampleGetCurlCommand() {
	req, _ := http.NewRequest("PUT", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", bytes.NewBufferString(`{"hello":"world","answer":42}`))
	req.Header.Set("Content-Type", "application/json")

	command, _ := GetCurlCommand(req)
	fmt.Println(command)

	// Output:
	// curl -X PUT -d "{\"hello\":\"world\",\"answer\":42}" -H "Content-Type: application/json" 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'
}

func TestGetCurlCommand(t *testing.T) {
	uri := "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu"
	payload := new(bytes.Buffer)
	payload.Write([]byte(`{"hello":"world","answer":42}`))
	req, err := http.NewRequest("PUT", uri, payload)
	if err != nil {
		t.Fatalf("got error with request: %v", err)
	}
	req.Header.Set("X-Auth-Token", "private-token")
	req.Header.Set("Content-Type", "application/json")

	command, err := GetCurlCommand(req)
	if err != nil {
		t.Fatalf("got error with GetCurlCommand: %v", err)
	}
	expected := `curl -X PUT -d "{\"hello\":\"world\",\"answer\":42}" -H "Content-Type: application/json" -H "X-Auth-Token: private-token" 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'`
	if command.String() != expected {
		t.Fatalf("\nexpected: %s\ngot     : %s", expected, command.String())
	}
}
