package main

import (
	"github.com/unsuman/go-microservices/types"
)

type Storer interface {
	InsertDistance(types.Distance) error
}

type MemeoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemeoryStore {
	return &MemeoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemeoryStore) InsertDistance(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}
