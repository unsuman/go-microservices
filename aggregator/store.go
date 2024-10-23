package main

import "github.com/unsuman/go-microservices/types"

type Storer interface {
	InsertDistance(types.Distance) error
}

type MemeoryStore struct {
}

func NewMemoryStore() *MemeoryStore {
	return &MemeoryStore{}
}

func (m *MemeoryStore) InsertDistance(d types.Distance) error {
	return nil
}
