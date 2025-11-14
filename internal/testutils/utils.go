package testutils

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func PrintBody(t *testing.T, res *http.Response) {
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
