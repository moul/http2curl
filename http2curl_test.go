package http2curl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func ExampleGetCurlCommand() {
	form := url.Values{}
	form.Add("age", "10")
	form.Add("name", "Hudson")
	body := form.Encode()

	req, _ := http.NewRequest(http.MethodPost, "http://foo.com/cats", ioutil.NopCloser(bytes.NewBufferString(body)))
	req.Header.Set("API_KEY", "123")

	command, _ := GetCurlCommand(req)
	fmt.Println(command)

	// Output:
	// curl -X 'POST' -d 'age=10&name=Hudson' -H 'Api_key: 123' 'http://foo.com/cats'
}

func ExampleGetCurlCommand_json() {
	req, _ := http.NewRequest("PUT", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", bytes.NewBufferString(`{"hello":"world","answer":42}`))
	req.Header.Set("Content-Type", "application/json")

	command, _ := GetCurlCommand(req)
	fmt.Println(command)

	// Output:
	// curl -X 'PUT' -d '{"hello":"world","answer":42}' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'
}

func ExampleGetCurlCommand_slice() {
	// See https://github.com/moul/http2curl/issues/12
	req, _ := http.NewRequest("PUT", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", bytes.NewBufferString(`{"hello":"world","answer":42}`))
	req.Header.Set("Content-Type", "application/json")

	command, _ := GetCurlCommand(req)
	fmt.Println(strings.Join(*command, " \\\n  "))

	// Output:
	// curl \
	//   -X \
	//   'PUT' \
	//   -d \
	//   '{"hello":"world","answer":42}' \
	//   -H \
	//   'Content-Type: application/json' \
	//   'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'
}

func ExampleGetCurlCommand_noBody() {
	req, _ := http.NewRequest("PUT", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", nil)
	req.Header.Set("Content-Type", "application/json")

	command, _ := GetCurlCommand(req)
	fmt.Println(command)

	// Output:
	// curl -X 'PUT' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'
}

func ExampleGetCurlCommand_emptyStringBody() {
	req, _ := http.NewRequest("PUT", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	command, _ := GetCurlCommand(req)
	fmt.Println(command)

	// Output:
	// curl -X 'PUT' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'
}

func ExampleGetCurlCommand_newlineInBody() {
	req, _ := http.NewRequest("POST", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", bytes.NewBufferString("hello\nworld"))
	req.Header.Set("Content-Type", "application/json")

	command, _ := GetCurlCommand(req)
	fmt.Println(command)

	// Output:
	// curl -X 'POST' -d 'hello
	// world' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'
}

func ExampleGetCurlCommand_specialCharsInBody() {
	req, _ := http.NewRequest("POST", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", bytes.NewBufferString(`Hello $123 o'neill -"-`))
	req.Header.Set("Content-Type", "application/json")

	command, _ := GetCurlCommand(req)
	fmt.Println(command)

	// Output:
	// curl -X 'POST' -d 'Hello $123 o'\''neill -"-' -H 'Content-Type: application/json' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'
}

func ExampleGetCurlCommand_other() {
	uri := "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu"
	payload := new(bytes.Buffer)
	payload.Write([]byte(`{"hello":"world","answer":42}`))
	req, err := http.NewRequest("PUT", uri, payload)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", "private-token")
	req.Header.Set("Content-Type", "application/json")

	command, err := GetCurlCommand(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(command)
	// Output: curl -X 'PUT' -d '{"hello":"world","answer":42}' -H 'Content-Type: application/json' -H 'X-Auth-Token: private-token' 'http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu'
}

// Benchmark test for GetCurlCommand
func BenchmarkGetCurlCommand(b *testing.B) {
	form := url.Values{}

	for i := 0; i <= b.N; i++ {
		form.Add("number", strconv.Itoa(i))
		body := form.Encode()
		req, _ := http.NewRequest(http.MethodPost, "http://foo.com", ioutil.NopCloser(bytes.NewBufferString(body)))
		_, err := GetCurlCommand(req)
		if err != nil {
			panic(err)
		}
	}
}
