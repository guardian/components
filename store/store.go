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
	GetAll() []Component
}

type Setter interface {
	Set(Component)
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

func (s MemoryStore) GetAll() []Component {
	return s.Components
}
