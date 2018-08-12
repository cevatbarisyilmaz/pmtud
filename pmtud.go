//Package mtud provides Path MTU Discovery functionality.
package pmtud

import (
	"errors"
	"math"
	"net"
)

const (
	MinIPv4HeaderSize = 20
	MaxIPv4HeaderSize = 60
	IPv6HeaderSize    = 40
	ICMPHeaderSize    = 8
)

//Number of ICMP messages to send for each ping test
var NumberOfMessages = 3

//Timeout for ICPM echo replies, currently only used in Windows
var TimeOutDuration = 5000

var (
	TimeOutError  = errors.New("ping request timed out")
	Unimplemented = errors.New("pmtud is  not implemented for this OS")
	NoRecord      = errors.New("no IP record is found for given host")
)

//Performs Path MTU Discovery between given address. Address can be an IP address or domain name.
func Pmtud(addr string) (int, error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		records, err := net.LookupIP(addr)
		if err != nil {
			return 0, err
		}
		if len(records) == 0 {
			return 0, NoRecord
		}
		ip = records[0]
	}
	var overhead int
	if ip.To4() == nil {
		overhead = IPv6HeaderSize + ICMPHeaderSize
	} else {
		overhead = MinIPv4HeaderSize + ICMPHeaderSize
	}
	max := 65500
	min := 0
	for max != min {
		pmtu := int(math.Round(float64(max+min) / 2))
		good, rec, err := check(ip.String(), pmtu)
		if err != nil {
			return 0, err
		}
		if good {
			min = pmtu
		} else {
			if rec != 0 {
				good, _, err := check(ip.String(), rec-overhead)
				if err != nil {
					return 0, err
				}
				if good {
					good, _, err := check(ip.String(), rec+1-overhead)
					if err != nil {
						return 0, err
					}
					if !good {
						return rec, nil
					}
				}
			}
			max = pmtu - 1
		}
	}
	return min + overhead, nil
}
