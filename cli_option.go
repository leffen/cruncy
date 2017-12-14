package cruncy

import (
	"os"
	"strings"

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

// Viper returns viper object
func (c *CliOption) Viper() *viper.Viper {
	return c.v
}

// ReadConfig reads and parses the config options set
func (c *CliOption) ReadConfig() {
	c.v.AutomaticEnv()
	c.f.Parse(os.Args)
	c.v.ReadInConfig()
}

// ReadTomlConfigFile reads config file from a given folder
func (c *CliOption) ReadTomlConfigFile(configFolder, fileName string) error {
	c.v.AutomaticEnv()
	c.v.SetConfigType("toml")
	c.v.AddConfigPath(".")
	c.v.AddConfigPath(configFolder)
	c.v.SetConfigName(strings.TrimSuffix(fileName, ".toml"))
	err := c.v.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
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

// GetString returns string value for param
func (c *CliOption) GetString(name string) string {
	return c.v.GetString(name)
}

// GetInt returns int value for param
func (c *CliOption) GetInt(name string) int {
	return c.v.GetInt(name)
}

// GetBool returns bool value for param
func (c *CliOption) GetBool(name string) bool {
	return c.v.GetBool(name)
}
