package grpc

import "testing"

func TestNew(t *testing.T) {
	s := New(nil)
	if s == nil {
		t.Fatal("s was nil")
	}
}
