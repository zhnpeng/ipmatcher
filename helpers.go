package ipmatcher

import (
	"net"
)

func compare(ip1, ip2 net.IP) int {
	var i1, i2 int

	if len(ip1) > len(ip2) {
		i1 = len(ip1) - len(ip2)
		for i := i1; i > 0; i-- {
			if ip1[i] != 0 {
				return 1
			}
		}
	} else if len(ip2) > len(ip1) {
		i2 = len(ip2) - len(ip1)
		for i := i2; i > 0; i-- {
			if ip2[i] != 0 {
				return -1
			}
		}
	}

	for ; i1 < len(ip1); i1++ {
		if ip1[i1] < ip2[i2] {
			return -1
		} else if ip1[i1] > ip2[i2] {
			return 1
		}
		i2++
	}
	return 0
}

func lastIP(ip net.IP, mask net.IPMask) net.IP {
	var (
		n   = len(mask)
		j   = len(ip) - n
		out = make(net.IP, n)
	)

	for i := 0; i < n; i++ {
		out[i] = ip[j] | ^mask[i]
		j++
	}
	return out
}

func compareExt(ip net.IP, ipUint uint32, it *IPRange) int {
	if ipUint > 0 {
		if ipUint < it.StartIPUint {
			return -1
		}
		if ipUint == it.StartIPUint {
			return 0
		}
		if it.StartIPUint > 0 {
			return 1
		}
	} else if it.StartIPUint > 0 {
		return 1
	}
	return compare(ip, it.StartIP)
}
