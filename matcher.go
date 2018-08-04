package ipmatcher

import (
	"errors"
	"net"

	"github.com/google/btree"
)

// ErrInvalidItem message
var ErrInvalidItem = errors.New("Invalid IP range item")

// IPMatcher base
type IPMatcher struct {
	root *btree.BTree
}

// New tree object
func NewIPMatcher() *IPMatcher {
	return &IPMatcher{root: btree.New(2)}
}

// AddRange IPs vith value
func (t *IPMatcher) AddRange(ip1, ip2 net.IP, val interface{}) error {
	return t.Add(&IPRange{StartIP: ip1, EndIP: ip2, Data: val})
}

// Add IP range item
func (t *IPMatcher) Add(item *IPRange) error {
	if item == nil || (item.StartIP == nil && item.EndIP == nil) {
		return ErrInvalidItem
	}
	if item.StartIP == nil {
		item.StartIP = item.EndIP
	}
	if item.EndIP == nil {
		item.EndIP = item.StartIP
	}
	item.Normalize()
	t.root.ReplaceOrInsert(item)
	return nil
}

// Match item by IP
func (t *IPMatcher) Match(ip net.IP) (response *IPRange) {
	ipVal := ip2int(ip)
	t.root.AscendGreaterOrEqual(IP(ip), func(item btree.Item) bool {
		it := item.(*IPRange)
		switch compareExt(ip, ipVal, it) {
		case 1:
			if compare(ip, it.EndIP) <= 0 {
				response = it
				return false
			}
		case 0:
			response = it
			return false
		case -1:
			return false
		}
		return true
	})
	return
}

type PortMatcher struct {
	root *btree.BTree
}

func NewPortMatcher() *PortMatcher {
	return &PortMatcher{root: btree.New(2)}
}

func (t *PortMatcher) Add(item *PortRange) error {
	if item == nil {
		return ErrInvalidItem
	}
	t.root.ReplaceOrInsert(item)
	return nil
}

func (t *PortMatcher) Match(port Port) (response *PortRange) {
	t.root.AscendGreaterOrEqual(port, func(item btree.Item) bool {
		it := item.(*PortRange)
		if port < it.StartPort {
			return false
		} else if port == it.StartPort {
			response = it
			return false
		} else {
			if port < it.EndPort {
				response = it
				return false
			}
		}
		return true
	})
	return
}
