package cruncy

import (
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// CliOption encapsulates viper flag and values
type CliOption struct {
	v      *viper.Viper
	f      *pflag.FlagSet
	prefix string
}

// NewCliOption creates a new cli option object
func NewCliOption(prefix string) *CliOption {
	v := viper.New()
	v.SetEnvPrefix(prefix)
	f := pflag.NewFlagSet("sync", pflag.ContinueOnError)
	return &CliOption{v: v, f: f, prefix: prefix}
}

// ReadConfig reads and parses the config options set
func (c *CliOption) ReadConfig() {
	c.v.AutomaticEnv()
	c.f.Parse(os.Args)
	c.v.ReadInConfig()
}

// MakeString creates a application string variable both for env and command line
func (c *CliOption) MakeString(key, short, envName, defaultValue, description string) {
	c.v.SetDefault(key, defaultValue)
	_ = c.v.BindEnv(key, envName)
	if short != "" {
		c.f.StringP(key, short, defaultValue, description)
	} else {
		c.f.String(key, defaultValue, description)
	}
	c.v.BindPFlag(key, c.f.Lookup(key))
}

// MakeInt creates a application int variable both for env and command line
func (c *CliOption) MakeInt(key, short, envName string, defaultValue int, description string) {
	c.v.SetDefault(key, defaultValue)
	_ = c.v.BindEnv(key, envName)
	if short != "" {
		c.f.IntP(key, short, defaultValue, description)
	} else {
		c.f.Int(key, defaultValue, description)

	}
	c.v.BindPFlag(key, c.f.Lookup(key))
}

// MakeBool creates a application bool variable both for env and command line
func (c *CliOption) MakeBool(key, short, envName string, defaultValue bool, description string) {
	c.v.SetDefault(key, defaultValue)
	_ = c.v.BindEnv(key, envName)
	if short != "" {
		c.f.BoolP(key, short, defaultValue, description)
	} else {
		c.f.Bool(key, defaultValue, description)

	}
	c.v.BindPFlag(key, c.f.Lookup(key))
}
