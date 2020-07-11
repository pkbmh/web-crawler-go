package pkg

import "sync"

// TODO Add more useful functionality like GetAllUrls with a given state etc
// Change the implementation from using url string to using urlNode.

type Storage interface {
	GetUrlState(key interface{}) (value interface{}, ok bool)
	Store(key interface{}, value interface{}) error
	isPresent(key interface{}) bool
}

type inMemoryStore struct {
	store *sync.Map
}

func NewInMemoryStore() Storage {
	syncMap := new(sync.Map)
	return &inMemoryStore{
		store:syncMap,
	}
}

func (c *inMemoryStore) GetUrlState(key interface{}) (value interface{}, ok bool) {
	value, ok = c.store.Load(key)
	return value, ok
}

func (c *inMemoryStore) Store(key interface{}, value interface{}) error {
	c.store.Store(key, value)
	return nil
}

func (c *inMemoryStore) isPresent(key interface{}) bool {
	_, ok := c.GetUrlState(key)
	return ok
}