// Host detection.

package embd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// The Host type represents all the supported host types.
type Host int

const (
	// HostNull reprents a null host.
	HostNull Host = iota

	// HostRPi represents the RaspberryPi.
	HostRPi

	// HostBBB represents the BeagleBone Black.
	HostBBB

	// HostCubieTruck represents the Cubie Truck.
	HostCubieTruck

	// HostGalileo represents the Intel Galileo board.
	HostGalileo
)

func execOutput(name string, arg ...string) (output string, err error) {
	var out []byte
	if out, err = exec.Command(name, arg...).Output(); err != nil {
		return
	}
	output = strings.TrimSpace(string(out))
	return
}

func nodeName() (string, error) {
	return execOutput("uname", "-n")
}

func parseVersion(str string) (major, minor, patch int, err error) {
	parts := strings.Split(str, ".")
	len := len(parts)

	if major, err = strconv.Atoi(parts[0]); err != nil {
		return 0, 0, 0, err
	}
	if minor, err = strconv.Atoi(parts[1]); err != nil {
		return 0, 0, 0, err
	}
	if len > 2 {
		part := parts[2]
		part = strings.TrimSuffix(part, "+")
		if patch, err = strconv.Atoi(part); err != nil {
			return 0, 0, 0, err
		}
	}

	return major, minor, patch, err
}

func kernelVersion() (major, minor, patch int, err error) {
	output, err := execOutput("uname", "-r")
	if err != nil {
		return 0, 0, 0, err
	}

	return parseVersion(output)
}

// DetectHost returns the detected host and its revision number.
func DetectHost() (Host, int, error) {
	major, minor, patch, err := kernelVersion()
	if err != nil {
		return HostNull, 0, err
	}

	if major < 3 || (major == 3 && minor < 8) {
		return HostNull, 0, fmt.Errorf("embd: linux kernel versions lower than 3.8 are not supported. you have %v.%v.%v", major, minor, patch)
	}

	node, err := nodeName()
	if err != nil {
		return HostNull, 0, err
	}

	var host Host
	var rev int

	switch node {
	case "raspberrypi":
		host = HostRPi
	case "beaglebone":
		host = HostBBB
	default:
		return HostNull, 0, fmt.Errorf("embd: your host %q is not supported at this moment. please request support at https://github.com/kidoman/embd/issues", node)
	}

	return host, rev, nil
}
