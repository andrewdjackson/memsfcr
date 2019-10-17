package rosco

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Errorf("Unable to create a new Mems struct")
	}
}

func TestMemsConnect(t *testing.T) {
	c := config.ReadConfig()
	c.Port = "/dev/ttys0"
	m := New()
	MemsConnect(m, c)

	if m.SerialPort != nil {
		t.Errorf("Failed to connect")
	}
}
