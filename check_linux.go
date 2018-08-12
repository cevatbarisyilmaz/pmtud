package pmtud

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func check(addr string, size int) (bool, int, error) {
	for i := 0; i < 7; i++ {
		outByte, _ := exec.Command("ping", "-c", strconv.Itoa(NumberOfPings), "-s", strconv.Itoa(size), "-M", "do", addr).Output()
		outString := string(outByte)
		if strings.Contains(outString, "Message too long") {
			regex := regexp.MustCompile(`mut=\d+`)
			inv := regex.FindString(outString)
			mtuString := strings.TrimPrefix(inv, "mut=")
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
