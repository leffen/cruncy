package cruncy

import "testing"

func TestOptionCli(t *testing.T) {
	o := NewCliOption("test")
	o.MakeString("test", "t", "TEST", "OleBrum", "Test env variable")
	o.ReadConfig()
	x := o.v.GetString("test")
	if x != "OleBrum" {
		t.Fatalf("Unable to read correct var value %s is supposed to be OleBrum", x)
	}
}
