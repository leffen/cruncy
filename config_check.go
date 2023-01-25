package cruncy

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type ConfigCheck struct {
	errs []string
}

func (c *ConfigCheck) AddIfEmptyStr(title, value string) {
	if value != "" {
		return
	}
	if c.errs == nil {
		c.errs = []string{}
	}
	c.errs = append(c.errs, title)
	logrus.Debugf("Added %s to errors", title)
}

func (c *ConfigCheck) TerminateIfErrors() {
	if len(c.errs) == 0 {
		return
	}

	logrus.Fatalln("Missing variables ( or app params ): " + strings.Join(c.errs, " - "))
}
