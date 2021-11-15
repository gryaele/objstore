package server

import (
	"testing"
)

const (
	Key1 = "key1"
	Key2 = "key2"
	Data = "data"
)

// Test Object Deduplication
func TestBucketDeduplication(t *testing.T) {
	b := newBucket()
	data := []byte(Data)
	if err := b.put(Key1, data); err != nil {
		t.Errorf("first put failed: %v", err)
	}
	if err := b.put(Key2, data); err != nil {
		t.Errorf("second put failed: %v", err)
	}
	if len(b.blobs) != 1 {
		t.Errorf("len(b.blobs) = %v, want: 1", len(b.blobs))
	}
	if len(b.objectHashes) != 2 {
		t.Errorf("len(b.objectHashes) = %v, want: 2", len(b.objectHashes))
	}
	got, err := b.get(Key1)
	if err != nil {
		t.Errorf("failed to get key1: %v", err)
	}
	want := Data
	if string(got) != want {
		t.Errorf("b.get(%v) = %v, want: %v", Key1, string(got), want)
	}
	got, err = b.get(Key2)
	if err != nil {
		t.Errorf("failed to get key1: %v", err)
	}
	if string(got) != want {
		t.Errorf("b.get(%v) = %v, want: %v", Key2, string(got), want)
	}
}

// Check deletion of the object referenced once
func TestBucketUniqueObjectDelete(t *testing.T) {
	b := newBucket()
	data := []byte(Data)
	if err := b.put(Key1, data); err != nil {
		t.Errorf("first put failed: %v", err)
	}
	if err := b.delete(Key1); err != nil {
		t.Errorf("delete failed: %v", err)
	}
	if len(b.blobs) != 0 {
		t.Errorf("len(b.blobs) = %v, want: 0", len(b.blobs))
	}
	if len(b.objectHashes) != 0 {
		t.Errorf("len(b.objectHashes) = %v, want: 0", len(b.objectHashes))
	}
}

// Check deletion of the object referenced twice
func TestBucketReferencedObjectDelete(t *testing.T) {
	b := newBucket()
	data := []byte(Data)
	if err := b.put(Key1, data); err != nil {
		t.Errorf("first put failed: %v", err)
	}
	if err := b.put(Key2, data); err != nil {
		t.Errorf("second put failed: %v", err)
	}
	if err := b.delete(Key1); err != nil {
		t.Errorf("delete failed: %v", err)
	}
	if len(b.blobs) != 1 {
		t.Errorf("len(b.blobs) = %v, want: 1", len(b.blobs))
	}
	if len(b.objectHashes) != 1 {
		t.Errorf("len(b.objectHashes) = %v, want: 1", len(b.objectHashes))
	}
	if val, err := b.get(Key1); err != errNotFound {
		t.Errorf("b.get(\"key1\") = %v, %v: want: nil, %v", val, err, errNotFound)
	}
	val, err := b.get(Key2)
	want := Data
	if string(val) != want || err != nil {
		t.Errorf("b.get(\"key2\") = %v, %v: want: %v, nil", val, err, want)
	}
	if err := b.delete(Key2); err != nil {
		t.Errorf("delete failed: %v", err)
	}
	if len(b.blobs) != 0 {
		t.Errorf("len(b.blobs) = %v, want: 0", len(b.blobs))
	}
	if len(b.objectHashes) != 0 {
		t.Errorf("len(b.objectHashes) = %v, want: 0", len(b.objectHashes))
	}
}
