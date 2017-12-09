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

func TestBool(t *testing.T) {
	o := NewCliOption("asn")

	o.MakeBool("sync_up", "u", "SYNC_UP", false, "Sync directory up to S3")
	if o.GetBool("sync_up") {
		t.Error("Expected flag to be false by default")
	}
}
