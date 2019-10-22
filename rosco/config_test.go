package rosco

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig()
	if c == nil {
		t.Errorf("Unable to create a new config struct")
	}
	if c.Port != "ttycodereader" {
		t.Errorf("Unexpected value set for config Port %+v", c)
	}
}

func TestReadConfig(t *testing.T) {
	var c = ReadConfig()
	if c == nil {
		t.Errorf("Unable to read config")
	}
	if c.Port == "ttycodereader" {
		t.Errorf("Unexpected value set for config Port %+v", c)
	}
}
