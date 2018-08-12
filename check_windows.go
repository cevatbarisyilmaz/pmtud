package pmtud

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

func check(addr string, size int) (bool, int, error) {
	for i := 0; i < 7; i++ {
		outByte, _ := exec.Command("ping", "-n", strconv.Itoa(NumberOfMessages), "-l", strconv.Itoa(size), "-f", "-w", strconv.Itoa(TimeOutDuration), addr).Output()
		outString := string(outByte)
		if strings.Contains(outString, "Packet needs to be fragmented but DF set") {
			return false, 0, nil
		}
		if strings.Contains(outString, "(0% loss)") {
			return true, 0, nil
		}
		if strings.Contains(outString, "(100% loss)") {
			return false, 0, nil
		}
		if strings.Contains(outString, "Request timed out") {
			return false, 0, TimeOutError
		}
	}
	return false, 0, errors.New("something went wrong with pinging")
}
