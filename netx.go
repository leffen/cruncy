package cruncy

import "os"

// Hostname returns primary env HOST_HOSTNAME  secondary os.Hostname
func Hostname() string {
	name := os.Getenv("HOST_HOSTNAME")
	if len(name) > 0 {
		return name
	}
	name, _ = os.Hostname()
	return name
}
