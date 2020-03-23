package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nicl/components/store"
)

func main() {
	s := store.MemoryStore{}
	log.Fatal(http.ListenAndServe(":8080", Router(&s)))
}

func Router(s *store.MemoryStore) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/components", Get(s))
	r.Get("/components/{ID}", GetID(s))
	r.Post("/components", Create(s))
	r.Delete("/components/{ID}", DeleteID(s))
	return r
}

func Get(s store.Getter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		components, err := s.GetAll()
		if err != nil {
			response(w, 500, []byte(err.Error()))
			return
		}

		resp, _ := json.Marshal(components)
		w.Write(resp)
	}
}

func GetID(s store.Getter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "ID")
		c, err := s.Get(ID)
		if err != nil {
			response(w, 404, nil)
			return
		}

		resp, _ := json.Marshal(c)
		response(w, 200, resp)
	}
}

func DeleteID(s store.Deleter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "ID")
		err := s.Delete(ID)
		if err != nil {
			response(w, 500, nil)
			return
		}

		response(w, 200, nil)
	}
}

func Create(s store.Setter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			response(w, 500, []byte(err.Error()))
			return
		}

		var c store.Component
		err = c.Unmarshall(body)
		if err != nil {
			response(w, 400, []byte(err.Error()))
			return
		}

		err = s.Set(c)
		if err != nil {
			response(w, 500, []byte(err.Error()))
			return
		}

		response(w, 201, nil)
	}
}

func response(w http.ResponseWriter, code int, msg []byte) {
	w.WriteHeader(code)
	w.Write(msg)
}
