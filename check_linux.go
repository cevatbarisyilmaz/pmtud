package pmtud

import (
	"errors"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func check(ip net.IP, size int) (bool, int, error) {
	for i := 0; i < 7; i++ {
		outByte, _ := exec.Command("ping", "-c", strconv.Itoa(NumberOfMessages), "-s", strconv.Itoa(size), "-M", "do", ip.String()).Output()
		outString := string(outByte)
		if strings.Contains(outString, "Message too long") {
			regex := regexp.MustCompile(`mtu=\d+`)
			inv := regex.FindString(outString)
			mtuString := strings.TrimPrefix(inv, "mtu=")
			if mtuString != "" {
				mtu, err := strconv.Atoi(mtuString)
				if err == nil {
					return false, mtu, nil
				}
			}
			return false, 0, nil
		}
		if strings.Contains(outString, " 0% packet loss") {
			return true, 0, nil
		}
		if strings.Contains(outString, " 100% packet loss") {
			return false, 0, nil
		}
	}
	return false, 0, errors.New("something went wrong with pinging")
}
