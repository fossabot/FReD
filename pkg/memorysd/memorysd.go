package memorysd

import (
	"errors"
	"sync"
)

// Storage stores a map of keygroups by name.
type Storage struct {
	keygroups map[string]Keygroup
	sync.RWMutex
}

// Keygroup stores a map of items by id, has a maxkey to keep track of the unused ids.
type Keygroup struct {
	items  map[string]string
	sync.RWMutex
}

// New create a new Storage.
func New() (s *Storage) {
	s = &Storage{
		keygroups: make(map[string]Keygroup),
	}

	return
}

// Read returns an item with the specified id from the specified keygroup.
func (s *Storage) Read(kgname string, id string) (string, error) {
	if kgname == "" {
		return "", errors.New("invalid keygroup name")
	}

	s.RLock()
	kg, ok := s.keygroups[kgname]
	s.RUnlock()

	if !ok {
		return "", errors.New("no such keygroup")
	}

	kg.RLock()
	var value string
	value, ok = kg.items[id]
	kg.RUnlock()

	if !ok {
		return "", errors.New("no such item")
	}

	return value, nil
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(kgname string, id string, data string) error {
	if kgname == "" {
		return errors.New("invalid keygroup name")
	}

	if data == "" {
		return errors.New("empty data")
	}

	s.RLock()
	kg, ok := s.keygroups[kgname]

	if !ok {
		s.RUnlock()
		return errors.New("no such keygroup")
	}

	s.RUnlock()

	kg.Lock()

	kg.items[id] = data

	kg.Unlock()

	return nil
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(kgname string, id string) error {
	if kgname == "" {
		return errors.New("invalid keygroup name")
	}

	s.RLock()
	kg, ok := s.keygroups[kgname]

	if !ok {
		s.RUnlock()
		return errors.New("no such keygroup")
	}

	s.RUnlock()

	kg.RLock()
	_, ok = kg.items[id]
	kg.RUnlock()

	if !ok {
		return errors.New("no such item")
	}

	kg.Lock()
	delete(kg.items, id)
	kg.Unlock()

	return nil

}

// CreateKeygroup creates a new keygroup with the specified name in Storage.
func (s *Storage) CreateKeygroup(kgname string) error {
	if kgname == "" {
		return errors.New("invalid keygroup name")
	}

	s.RLock()
	kg, exists := s.keygroups[kgname]

	if exists {
		s.RUnlock()
		return errors.New("keygroup exists")
	}

	s.RUnlock()

	kg = Keygroup{
		items:  make(map[string]string),
	}

	s.Lock()
	s.keygroups[kgname] = kg
	s.Unlock()

	return nil
}

// DeleteKeygroup removes the keygroup with the specified name from Storage.
func (s *Storage) DeleteKeygroup(kgname string) error {
	if kgname == "" {
		return errors.New("invalid keygroup name")
	}

	s.RLock()
	_, ok := s.keygroups[kgname]
	s.RUnlock()

	if !ok {
		return errors.New("keygroup does not exist")
	}

	s.Lock()
	delete(s.keygroups, kgname)
	s.Unlock()

	return nil
}
