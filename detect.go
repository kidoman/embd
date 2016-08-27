// Host detection.

package embd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

// The Host type represents all the supported host types.
type Host string

const (
	// HostNull reprents a null host.
	HostNull Host = ""

	// HostRPi represents the RaspberryPi.
	HostRPi = "Raspberry Pi"

	// HostBBB represents the BeagleBone Black.
	HostBBB = "BeagleBone Black"

	// HostGalileo represents the Intel Galileo board.
	HostGalileo = "Intel Galileo"

	// HostCubieTruck represents the Cubie Truck.
	HostCubieTruck = "CubieTruck"

	// HostRadxa represents the Radxa board.
	HostRadxa = "Radxa"

	// HostCHIP represents the NextThing C.H.I.P.
	HostCHIP = "CHIP"
)

func execOutput(name string, arg ...string) (output string, err error) {
	var out []byte
	if out, err = exec.Command(name, arg...).Output(); err != nil {
		return
	}
	output = strings.TrimSpace(string(out))
	return
}

func parseVersion(str string) (major, minor, patch int, err error) {
	versionNumber := strings.Split(str, "-")
	parts := strings.Split(versionNumber[0], ".")
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

func cpuInfo() (model, hardware string, revision int, err error) {
	output, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", "", 0, err
	}
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Split(line, ":")
		if len(fields) < 1 {
			continue
		}
		switch {
		case strings.HasPrefix(fields[0], "Revision"):
			revs := strings.TrimSpace(fields[1])
			rev, err := strconv.ParseInt(revs, 16, 32)
			if err != nil {
				continue
			}
			revision = int(rev)
		case strings.HasPrefix(fields[0], "Hardware"):
			hardware = strings.TrimSpace(fields[1])
		case strings.HasPrefix(fields[0], "model name"):
			model = fields[1]
		}
	}
	return model, hardware, revision, nil
}

// DetectHost returns the detected host and its revision number.
func DetectHost() (host Host, rev int, err error) {
	major, minor, patch, err := kernelVersion()
	if err != nil {
		return HostNull, 0, err
	}

	if major < 3 || (major == 3 && minor < 8) {
		return HostNull, 0, fmt.Errorf(
			"embd: linux kernel versions lower than 3.8 are not supported, "+
				"you have %v.%v.%v", major, minor, patch)
	}

	model, hardware, rev, err := cpuInfo()
	if err != nil {
		return HostNull, 0, err
	}

	switch {
	case strings.Contains(model, "ARMv7") && (strings.Contains(hardware, "AM33XX") || strings.Contains(hardware, "AM335X")):
		return HostBBB, rev, nil
	case strings.Contains(hardware, "BCM2708") || strings.Contains(hardware, "BCM2709"):
		return HostRPi, rev, nil
	case hardware == "Allwinner sun4i/sun5i Families":
		if major < 4 || (major == 4 && minor < 4) {
			return HostNull, 0, fmt.Errorf(
				"embd: linux kernel version 4.4+ required, you have %v.%v",
				major, minor)
		}
		return HostCHIP, rev, nil
	default:
		return HostNull, 0, fmt.Errorf(`embd: your host "%v:%v" is not supported at this moment. request support at https://github.com/kidoman/embd/issues`, host, model)
	}
}
