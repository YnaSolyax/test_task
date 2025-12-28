package storage

import (
	"errors"
	"sync"
)

type Inmemory struct {
	mu      sync.Mutex
	storage map[int]*Todo
	nextId  int
}

func NewInmemory() *Inmemory {
	return &Inmemory{
		storage: make(map[int]*Todo, 10),
		nextId:  1,
	}
}

func (im *Inmemory) Create(todo *Todo) error {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.storage[im.nextId] = todo
	im.nextId++
	return nil
}

func (im *Inmemory) GetAll() ([]*Todo, error) {
	im.mu.Lock()
	defer im.mu.Unlock()
	if len(im.storage) == 0 {
		return nil, nil
	}

	todos := make([]*Todo, 0, len(im.storage))

	for _, t := range im.storage {
		todos = append(todos, t)
	}

	return todos, nil
}

func (im *Inmemory) Get(id int) (*Todo, error) {
	im.mu.Lock()
	defer im.mu.Unlock()
	todo := im.storage[id]
	return todo, nil
}

func (im *Inmemory) Update(id int, t Todo) error {
	im.mu.Lock()
	defer im.mu.Unlock()
	m, ok := im.storage[id]
	if !ok {
		return ErrNotFound
	}
	if t.Title != "" {
		m.Title = t.Title
	}
	if t.Description.IsSet {
		m.Description = t.Description
	}
	if t.Status.IsSet {
		m.Status.Value = t.Status.Value
	}

	return nil
}

func (im *Inmemory) Delete(id int) error {
	im.mu.Lock()
	defer im.mu.Unlock()
	_, ok := im.storage[id]
	if !ok {
		return ErrNotFound
	}
	delete(im.storage, id)
	return nil
}

var ErrNotFound = errors.New("not found")
