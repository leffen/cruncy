package cruncy

import (
	"os"
	"strconv"
)

// GetEnvUnlessFlagIsNotDefault Getting String env variable unless flagvalue differs from defaultvalue
func GetEnvUnlessFlagIsNotDefault(envVar, defaultValue, flagValue, flagDefaultValue string) string {
	// log.Infof("%s defaultValue :%s flagValue:%s flagDefaultvalue: %s ENV=%s\n", envVar, defaultValue, flagValue, flagDefaultValue, os.Getenv(envVar))
	if flagValue != flagDefaultValue && flagDefaultValue != "" {
		return flagValue
	}

	if os.Getenv(envVar) != "" {
		return os.Getenv(envVar)
	}

	if flagValue != "" {
		return flagValue
	}

	if defaultValue != "" {
		return defaultValue
	}

	return flagDefaultValue

}

// GetEnvIntUnlessFlagIsNotDefault Getting INT env variable unless flagvalue differs from defaultvalue
func GetEnvIntUnlessFlagIsNotDefault(envVar string, defaultValue, flagValue, flagDefaultValue int) int {
	if flagValue != flagDefaultValue {
		return flagValue
	}

	if os.Getenv(envVar) != "" {
		i, err := strconv.Atoi(os.Getenv(envVar))
		if err != nil {
			return defaultValue
		}

		return i
	}
	return defaultValue
}

// GetEnvBoolOverrideFlag if environment flag is set to 1 it returns true
func GetEnvBoolOverrideFlag(currValue bool, flagName string) bool {
	if currValue {
		return currValue
	}

	return os.Getenv(flagName) == "1"
}
