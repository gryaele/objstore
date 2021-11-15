package server

import (
	"sync"
)

type memStore struct {
	mu      sync.Mutex
	buckets map[string]*bucket
}

func newMemStore() *memStore {
	return &memStore{
		buckets: make(map[string]*bucket),
	}
}

func (s *memStore) get(bucket string, key string) ([]byte, error) {
	s.mu.Lock()
	b, ok := s.buckets[bucket]
	s.mu.Unlock()
	if !ok {
		return nil, errNotFound
	}
	return b.get(key)
}

func (s *memStore) put(bucketName string, key string, data []byte) error {
	s.mu.Lock()
	b, ok := s.buckets[bucketName]
	if !ok {
		b = newBucket()
		s.buckets[bucketName] = b
	}
	s.mu.Unlock()

	return b.put(key, data)
}

func (s *memStore) delete(bucket string, key string) error {
	s.mu.Lock()
	b, ok := s.buckets[bucket]
	s.mu.Unlock()
	if !ok {
		return errNotFound
	}
	if err := b.delete(key); err != nil {
		return err
	}

	// When the last object of the bucket is deleted we should delete garbage collect the bucket.
	// This is left out of scope for this task.
	return nil
}
