package server

import "testing"

func TestStoreGetObject(t *testing.T) {
	s := newMemStore()
	if err := s.put("b1", "key1", []byte("foo")); err != nil {
		t.Errorf("s.put(b1, key1, foo) = %v, want: nil", err)
	}
	val, err := s.get("b1", "key1")
	if err != nil || string(val) != "foo" {
		t.Errorf("s.get(b1, key1) = %v, %v, want: foo, nil", val, err)
	}
}

func TestStoreGetObjectNotFound(t *testing.T) {
	s := newMemStore()
	val, err := s.get("b1", "key1")
	if err != errNotFound {
		t.Errorf("s.get(b1, key1) = %v, %v, want: nil, errNotFound", val, err)
	}
}

func TestStoreDelete(t *testing.T) {
	s := newMemStore()
	if err := s.put("b1", "key1", []byte("foo")); err != nil {
		t.Errorf("s.put(b1, key1, foo) = %v, want: nil", err)
	}
	if err := s.delete("b1", "key1"); err != nil {
		t.Errorf("s.delete(b1, key1) = %v, want: nil", err)
	}
	if err := s.delete("b1", "key1"); err != errNotFound {
		t.Errorf("s.delete(b1, key1) = %v, want: errNotFound", err)
	}
}
