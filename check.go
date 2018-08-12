// +build !windows,!linux

package pmtud

func check(addr string, size int) (bool, int, error) {
	return false, 0, Unimplemented
}
