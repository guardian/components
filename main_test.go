package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicl/components/store"
)

var c101 = store.Component{
	ID: "101",
}

func TestGetID(t *testing.T) {
	type expected struct {
		statusCode int
	}

	cases := map[string]expected{
		"101": expected{statusCode: 200},
		"202": expected{statusCode: 404},
	}

	s := store.MemoryStore{Components: []store.Component{c101}}
	r := Router(&s)
	ts := httptest.NewServer(http.HandlerFunc(r.ServeHTTP))
	defer ts.Close()

	for ID, exp := range cases {
		res, err := http.Get(ts.URL + "/components/" + ID)
		assertNoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assertNoError(t, err)

		assert(t, res.StatusCode == exp.statusCode, notEq(res.StatusCode, exp.statusCode))

		if res.StatusCode == 200 {
			var c store.Component
			err = c.Unmarshall(body)
			assertNoError(t, err)
			assert(t, c.ID == ID, notEq(c.ID, ID))
		}

	}
}

func TestCreate(t *testing.T) {
	s := store.MemoryStore{Components: []store.Component{}}
	ts := httptest.NewServer(http.HandlerFunc(Create(&s)))
	defer ts.Close()

	component := store.Component{ID: "foo"}
	res, err := http.Post(ts.URL+"/components", "application/json", bytes.NewReader(component.Marshall()))
	assertNoError(t, err)

	assert(t, res.StatusCode == 201, notEq(res.StatusCode, 200))
	assert(t, len(s.Components) == 1, "expected component to be stored")
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}

func assert(t *testing.T, test bool, msg string) {
	if !test {
		t.Error(msg)
	}
}

func notEq(a, b interface{}) string {
	return fmt.Sprintf("got %v; want %v", a, b)
}
