package host

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Host int

const (
	Null Host = iota
	RPi
	BBB
)

func execOutput(name string, arg ...string) (output string, err error) {
	var out []byte
	if out, err = exec.Command(name, arg...).Output(); err != nil {
		return
	}
	output = string(out)
	return
}

func nodeName() (string, error) {
	return execOutput("uname", "-n")
}

func kernelVersion() (major, minor, patch int, err error) {
	output, err := execOutput("uname", "-r")
	if err != nil {
		return
	}

	parts := strings.Split(output, ".")

	if major, err = strconv.Atoi(parts[0]); err != nil {
		return
	}
	if minor, err = strconv.Atoi(parts[1]); err != nil {
		return
	}
	if patch, err = strconv.Atoi(parts[2]); err != nil {
		return
	}

	return
}

func Detect() (host Host, rev int, err error) {
	major, minor, patch, err := kernelVersion()
	if err != nil {
		return
	}

	if major < 3 || (major == 3 && minor < 8) {
		err = fmt.Errorf("embd: linux kernel versions lower than 3.8 are not supported. you have %v.%v.%v", major, minor, patch)
		return
	}

	node, err := nodeName()
	if err != nil {
		return
	}

	switch node {
	case "raspberrypi":
		host = RPi
	case "beaglebone":
		host = BBB
	default:
		err = fmt.Errorf("embd: your host %q is not supported at this moment. please request support at https://github.com/kidoman/embd/issues", node)
	}

	return
}
