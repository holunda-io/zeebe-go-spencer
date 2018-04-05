package common

import "os"

const zeebeHostEnvVar = "ZEEBE_HOST"

// Get the zeebe hostname to connect to.
// Either from the first commandline argument
// or environment variable ZEEBE_HOST
func GetZeebeHost() string {
	zeebeHost := ""
	if len(os.Args) >= 2 {
		zeebeHost = os.Args[1]
	}
	if zeebeHost == "" {
		zeebeHost = os.Getenv(zeebeHostEnvVar)
	}
	return zeebeHost
}
