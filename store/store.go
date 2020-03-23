package store

import (
	"encoding/json"
	"errors"
)

type Component struct {
	ID          string
	Name        string
	Description string
	Owner       string
	HTML        string
	CSS         string
	Version     int
	Tags        []string
}

func (c Component) Marshall() []byte {
	bytes, _ := json.Marshal(c)
	return bytes
}

func (c *Component) Unmarshall(data []byte) error {
	return json.Unmarshal(data, c)
}

type Getter interface {
	Get(string) (Component, error)
	GetAll() ([]Component, error)
}

type Setter interface {
	Set(Component) error
}

type Deleter interface {
	Delete(string) error
}

type MemoryStore struct {
	Components []Component
}

func (s MemoryStore) Get(ID string) (Component, error) {
	for _, c := range s.Components {
		if c.ID == ID {
			return c, nil
		}
	}

	return Component{}, errors.New("item not found")
}

func (s MemoryStore) GetAll() ([]Component, error) {
	return s.Components, nil
}

func (s *MemoryStore) Set(c Component) error {
	s.Components = append(s.Components, c)
	return nil
}

func (s *MemoryStore) Delete(ID string) error {
	var data []Component
	for _, c := range s.Components {
		if c.ID != ID {
			data = append(data, c)
		}
	}

	s.Components = data
	return nil
}
