package server

import (
	"crypto/sha256"
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	errNotFound      = errors.New("object not found")
	errAlreadyExists = errors.New("object already exists")
)

type content struct {
	data []byte
	refs int // Protected by bucket.mu
}

type hash string

type bucket struct {
	mu           sync.Mutex
	blobs        map[hash]*content
	objectHashes map[string]hash
}

func newBucket() *bucket {
	return &bucket{
		mu:           sync.Mutex{},
		blobs:        make(map[hash]*content),
		objectHashes: make(map[string]hash),
	}
}

// Add object to a bucket
func (b *bucket) put(key string, data []byte) error {
	// Calculate hash sum
	hs := sha256.Sum256(data)
	hv := hash(hs[:])

	b.mu.Lock()
	defer b.mu.Unlock()

	// We can potentially allow object overwrites, this would require deleting previous value.
	// This version doesn't permit putting an object with already existing key
	// Check if the key with such value already exists
	if _, exists := b.objectHashes[key]; exists {
		return errAlreadyExists
	}

	b.objectHashes[key] = hv
	// Check if object with this hash sum already exists
	if blob, ok := b.blobs[hv]; ok {
		// Update references count
		blob.refs++
		return nil
	}
	// Save object
	b.blobs[hv] = &content{
		data: data,
		refs: 1,
	}
	return nil
}

// delete object with the give key from the bucket
func (b *bucket) delete(key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	hv, found := b.objectHashes[key]
	if !found {
		return errNotFound
	}
	delete(b.objectHashes, key)

	blob, ok := b.blobs[hv]
	if !ok {
		log.Fatalf("internal error: no content file for key %s", key)
	}
	blob.refs--
	// Check if last reference to the object is removed
	if blob.refs == 0 {
		delete(b.blobs, hv)
	}
	return nil
}

// get object from the bucket
func (b *bucket) get(key string) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	hv, ok := b.objectHashes[key]
	if !ok {
		return nil, errNotFound
	}
	blob, ok := b.blobs[hv]
	if !ok {
		log.Fatalf("invariant failure: no blob for %s", key)
	}
	return blob.data, nil
}
