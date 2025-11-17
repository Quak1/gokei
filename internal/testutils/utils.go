package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Quak1/gokei/internal/appcontext"
	"github.com/Quak1/gokei/internal/database/store"
)

func PrintBody(t *testing.T, res *http.Response) {
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func CreatePostRequest(t *testing.T, route string, requestBody any, user *store.User) *http.Request {
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(body))
	req.Header.Set("Contenty-Type", "application/json")
	req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
		ID:       user.ID,
		Username: user.Username,
	})

	return req
}

func CreateGetRequest(t *testing.T, route string, user *store.User) *http.Request {
	req := httptest.NewRequest(http.MethodGet, route, nil)
	req = appcontext.SetContextUser(req, &store.GetUserFromTokenRow{
		ID:       user.ID,
		Username: user.Username,
	})

	return req
}
