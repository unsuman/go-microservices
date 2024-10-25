package main

import (
	"fmt"

	"github.com/unsuman/go-microservices/types"
)

type Storer interface {
	InsertDistance(types.Distance) error
	GetDistance(int64) (float64, error)
}

type MemeoryStore struct {
	data map[int64]float64
}

func NewMemoryStore() *MemeoryStore {
	return &MemeoryStore{
		data: make(map[int64]float64),
	}
}

func (m *MemeoryStore) InsertDistance(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *MemeoryStore) GetDistance(obuid int64) (float64, error) {
	if _, ok := m.data[obuid]; !ok {
		return 0.0, fmt.Errorf("no distance found for obuid: %d", obuid)
	}

	return m.data[obuid], nil
}
