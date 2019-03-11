package di

import "testing"

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("container == nil")
	}
}

func TestGetContainer(t *testing.T) {
	c := New()
	container, err := c.Container()
	if err != nil {
		t.Fatal(err)
	}
	if container.GRPCServer() == nil {
		t.Fatal("grpcServer == nil")
	}
	if container.GRPCListener() == nil {
		t.Fatal("grpcListener == nil")
	}
	if container.Logger() == nil {
		t.Fatal("logger == nil")
	}
}
