package cruncy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionCli(t *testing.T) {
	o := NewCliOption("test")
	o.MakeString("test", "t", "TEST", "OleBrum", "Test env variable")
	o.ReadConfig()
	x := o.v.GetString("test")
	if x != "OleBrum" {
		t.Fatalf("Unable to read correct var value %s is supposed to be OleBrum", x)
	}
}

func TestBoolEnv(t *testing.T) {
	os.Setenv("SYNC_UP", "true")

	o := NewCliOption("test")
	o.MakeBool("sync_up", "u", "SYNC_UP", false, "Sync directory up to S3")
	o.ReadConfig()

	assert.Equal(t, true, o.GetBool("sync_up"))
}
func TestBoolCli(t *testing.T) {
	os.Clearenv()
	os.Args = []string{os.Args[0], "-u", "true"}

	o := NewCliOption("test")
	o.MakeBool("sync_up", "u", "SYNC_UP", false, "Sync directory up to S3")
	o.ReadConfig()

	assert.Equal(t, true, o.GetBool("sync_up"))

}

func TestIntEnv(t *testing.T) {
	os.Setenv("SYNC_UP", "177")

	o := NewCliOption("test")
	o.MakeInt("sync_up", "u", "SYNC_UP", 1, "Sync directory up to S3")
	o.ReadConfig()
	assert.Equal(t, 177, o.GetInt("sync_up"))
}

func TestIntCli(t *testing.T) {
	os.Clearenv()
	os.Args = []string{os.Args[0], "-u", "66"}

	o := NewCliOption("test")
	o.MakeInt("sync_up", "u", "SYNC_UP", 99, "Sync directory up to S3")
	o.ReadConfig()
	assert.Equal(t, 66, o.GetInt("sync_up"))
}
