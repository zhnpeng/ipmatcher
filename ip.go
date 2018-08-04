package ipmatcher

import (
	"encoding/binary"
	"errors"
	"net"
	"strings"

	"github.com/google/btree"
)

// ErrInvalidItemParse by string
var ErrInvalidItemParse = errors.New("Invalid parse item")

// IP implementation btree campare
type IP net.IP

// Less camparing for btree
func (ip IP) Less(then btree.Item) bool {
	switch v := then.(type) {
	case *IPRange:
		return compare(net.IP(ip), v.StartIP) == -1
	case IP:
		return compare(net.IP(ip), net.IP(v)) == -1
	}
	return false
}

// IPRange IP range
type IPRange struct {
	StartIPUint uint32
	StartIP     net.IP
	EndIP       net.IP
	Data        interface{}
}

// StringToIPRange tranform string to IPRange
func StringToIPRange(s string, val interface{}) (item *IPRange, err error) {
	s = strings.Trim(s, " \t\n_–+")
	if strings.Contains(s, "-") {
		if arr := strings.Split(s, "-"); 2 == len(arr) {
			item = &IPRange{
				StartIP: net.ParseIP(arr[0]),
				EndIP:   net.ParseIP(arr[1]),
				Data:    val,
			}
		}
	} else if strings.Contains(s, "/") {
		if ip, inet, e := net.ParseCIDR(s); nil == e {
			lip := lastIP(inet.IP, inet.Mask)
			item = &IPRange{
				StartIP: ip,
				EndIP:   net.ParseIP(lip.String()),
				Data:    val,
			}
		} else {
			err = e
		}
	} else {
		item = &IPRange{StartIP: net.ParseIP(s)}
		item.EndIP = item.StartIP
	}

	if err == nil {
		if item == nil || item.StartIP == nil || item.EndIP == nil {
			err = ErrInvalidItemParse
		} else if compare(item.StartIP, item.EndIP) > 0 {
			item.StartIP, item.EndIP = item.EndIP, item.StartIP
		}
	}
	return
}

// Less camparing for btree
func (i *IPRange) Less(then btree.Item) bool {
	switch ip := then.(type) {
	case *IPRange:
		if i.StartIPUint > 0 {
			if i.StartIPUint < ip.StartIPUint {
				return true
			}
		} else if ip.StartIPUint > 0 {
			return false
		}
		// replace only if both start ip and end ip are equals
		// or store all item if their start ip are the same
		// then sort by their end ip
		return compare(i.EndIP, ip.EndIP) == -1
	case IP:
		return compare(i.EndIP, net.IP(ip)) == -1
	}
	return false
}

// Compare with the second item
func (i *IPRange) Compare(it interface{}) int {
	switch ip := it.(type) {
	case *IPRange:
		if i.StartIPUint > 0 {
			if i.StartIPUint < ip.StartIPUint {
				return -1
			}
			if i.StartIPUint == ip.StartIPUint {
				return 0
			}
			if ip.StartIPUint > 0 {
				return 1
			}
		} else if ip.StartIPUint > 0 {
			return -1
		}
		return compare(i.StartIP, ip.StartIP)
	case IP:
		return compare(i.EndIP, net.IP(ip))
	case net.IP:
		return compare(i.EndIP, ip)
	}
	return 0
}

// Has IP in range
func (i *IPRange) Has(ip net.IP) bool {
	return compare(i.StartIP, ip) <= 0 && compare(i.EndIP, ip) >= 0
}

// Normalize IP values
func (i *IPRange) Normalize() {
	if i.StartIP != nil {
		if ip := i.StartIP.To4(); ip != nil {
			i.StartIPUint = ip2int(i.StartIP)
		}
	}
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}
