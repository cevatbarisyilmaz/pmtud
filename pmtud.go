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
	max := 65520
	min := 0
	good, rec, err := check(ip, 1500-overhead)
	if err != nil {
		return 0, err
	}
	if good {
		min = 1500
		good, rec, err = check(ip, 1501-overhead)
		if err != nil {
			return 0, err
		}
		if !good {
			return 1500, nil
		}
		min = 1501
	}
	if min == 1501 {
		good, rec, err = check(ip, 65520-overhead)
		if err != nil {
			return 0, err
		}
		if good {
			return 65520, nil
		}
	}
	for rec != 0 {
		good, nRec, err := check(ip, rec-overhead)
		if err != nil {
			return 0, err
		}
		if good {
			good, _, err = check(ip, rec+1-overhead)
			if err != nil {
				return 0, err
			}
			if !good {
				return rec, nil
			}
			if min < rec+1 {
				min = rec + 1
				break
			}
		} else {
			if nRec >= rec {
				break
			}
		}
		rec = nRec
	}
	for max != min {
		pmtu := int(math.Round(float64(max+min) / 2))
		good, _, err := check(ip, pmtu)
		if err != nil {
			return 0, err
		}
		if good {
			min = pmtu
		} else {
			max = pmtu - 1
		}
	}
	return min + overhead, nil
}
