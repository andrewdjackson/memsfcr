package rosco

import (
	"testing"
)

func TestNew(t *testing.T) {
	mems := New()
	if mems == nil {
		t.Errorf("Unable to create a new Mems struct")
	}
}

func TestMemsConnect(t *testing.T) {
	port := "/Users/ajackson/ttyecu"
	mems := New()
	MemsConnect(mems, port)

	if mems.SerialPort != nil {
		t.Errorf("Failed to connect")
	}
}
