package pmtud

import (
	"errors"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

func check(ip net.IP, size int) (bool, int, error) {
	for i := 0; i < 7; i++ {
		var outByte []byte
		var err error
		if ip.To4() == nil {
			outByte, err = exec.Command("ping", "-n", strconv.Itoa(NumberOfMessages), "-l", strconv.Itoa(size), "-w", strconv.Itoa(TimeOutDuration), ip.String()).Output()
		} else {
			outByte, err = exec.Command("ping", "-n", strconv.Itoa(NumberOfMessages), "-l", strconv.Itoa(size), "-f", "-w", strconv.Itoa(TimeOutDuration), ip.String()).Output()
		}
		outString := string(outByte)
		if strings.Contains(outString, "Packet needs to be fragmented but DF set") {
			return false, 0, nil
		}
		if strings.Contains(outString, "(0% loss)") && err == nil {
			return true, 0, nil
		}
		if strings.Contains(outString, "(100% loss)") {
			return false, 0, nil
		}
	}
	return false, 0, errors.New("something went wrong with pinging")
}
