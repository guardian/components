package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	ts := httptest.NewServer(http.HandlerFunc(GetID(s)))
	defer ts.Close()

	for ID, exp := range cases {
		res, err := http.Get(ts.URL + "/components/" + ID)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Error(err.Error())
		}

		if res.StatusCode != exp.statusCode {
			t.Errorf("got (%d), want (%d)", res.StatusCode, exp.statusCode)
		}

		// UnMarshall
		if exp.statusCode != 200 {
			return
		}

		var c store.Component
		err = json.Unmarshal(body, &c)
		if err != nil {
			t.Error(err.Error())
		}

		if c.ID != ID {
			t.Errorf("got ID (%s), want (%s)", c.ID, ID)
		}

	}

	t.Error("failed (expected)")
}
