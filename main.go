package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/nicl/components/store"
)

func main() {
	s := store.MemoryStore{}
	http.HandleFunc("/components", GetID(s))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetID(s store.Getter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := strings.TrimPrefix(r.URL.Path, "/components/")

		c, err := s.Get(ID)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Not found"))
		}

		resp, _ := json.Marshal(c)
		w.Write(resp)

	}
}
