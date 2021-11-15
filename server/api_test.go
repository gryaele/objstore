package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateObject(t *testing.T) {
	api := NewApi()
	ts := httptest.NewServer(api.NewRouter())
	defer ts.Close()
	content := "foo"
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/objects/bucket1/object1", ts.URL), bytes.NewBuffer([]byte(content)))
	if err != nil {
		t.Errorf("put request failed: %v", err)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("put request failed: %v", err)
	}
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed to read body: %v", err)
	}
	if err = rsp.Body.Close(); err != nil {
		t.Errorf("failed to close response: %v", err)
	}
	want := `{"id":"object1"}`
	if string(got) != want {
		t.Errorf("Get(/bucket1/object1) = %v, want: %v", string(got), want)
	}
	if rsp.StatusCode != http.StatusCreated {
		t.Errorf("Status Code = %v, want: 201", rsp.StatusCode)
	}

}

func TestObjectNotFound(t *testing.T) {
	api := NewApi()
	ts := httptest.NewServer(api.NewRouter())
	defer ts.Close()

	rsp, err := http.Get(fmt.Sprintf("%s/bucket1/object1", ts.URL))
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	if rsp.StatusCode != http.StatusNotFound {
		t.Errorf("Status Code = %v, want: 404", rsp.StatusCode)
	}
}

func TestDeleteObject(t *testing.T) {
	store := newMemStore()
	if err := store.put("b1", "key1", []byte("foo")); err != nil {
		t.Errorf(" store.put(b1, key1, foo) = %v, want nil", err)
	}
	api := &Api{
		objStore: store,
	}
	ts := httptest.NewServer(api.NewRouter())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/objects/b1/key1", ts.URL), nil)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("delete request failed: %v", err)
	}

	if rsp.StatusCode != http.StatusOK {
		t.Errorf("Status Code = %v, want: 200", rsp.StatusCode)
	}
}

func TestDeleteObjectNotFound(t *testing.T) {
	api := NewApi()
	ts := httptest.NewServer(api.NewRouter())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/objects/b1/object1", ts.URL), nil)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("delete request failed: %v", err)
	}

	if rsp.StatusCode != http.StatusNotFound {
		t.Errorf("Status Code = %v, want: 404", rsp.StatusCode)
	}
}

func TestGetObject(t *testing.T) {
	store := newMemStore()
	if err := store.put("b1", "key1", []byte("foo")); err != nil {
		t.Errorf(" store.put(b1, key1, foo) = %v, want nil", err)
	}
	api := &Api{
		objStore: store,
	}
	ts := httptest.NewServer(api.NewRouter())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/objects/bucket1/object1", ts.URL), nil)
	if err != nil {
		t.Errorf("http request failed: %v", err)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("delete request failed: %v", err)
	}

	if rsp.StatusCode != http.StatusNotFound {
		t.Errorf("Status Code = %v, want: 404", rsp.StatusCode)
	}
}
