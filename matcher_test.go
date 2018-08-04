package ipmatcher

import (
	"net"
	"testing"
)

func TestIPRange(t *testing.T) {
	matcher := NewIPMatcher()
	item1, _ := StringToIPRange("86.100.32.0/24", "A")
	item2, _ := StringToIPRange("86.100.32.5-86.100.32.100", "B")
	item3, _ := StringToIPRange("86.100.32.5-86.100.33.50", "C")
	matcher.Add(item1)
	matcher.Add(item2)
	matcher.Add(item3)
	matcher.AddRange(net.ParseIP("10.1.2.1"), net.ParseIP("10.1.2.50"), "D")
	matcher.AddRange(net.ParseIP("10.1.2.1"), net.ParseIP("10.16.3.100"), "E")
	matcher.AddRange(net.ParseIP("10.1.2.1"), net.ParseIP("10.1.2.50"), "F")
	matcher.AddRange(net.ParseIP("10.1.2.1"), net.ParseIP("10.1.2.255"), "G")

	ip1 := "10.1.2.40"
	got := matcher.Match(net.ParseIP(ip1))
	if got == nil {
		t.Errorf("%v should match but not", ip1)
	} else {
		if got.Data.(string) != "F" {
			t.Errorf("got unexcepted data got %v, want %v", got.Data.(string), "F")
		}
	}

	ip2 := "10.1.2.100"
	got = matcher.Match(net.ParseIP(ip2))
	if got == nil {
		t.Errorf("%v should match but not", ip2)
	} else {
		if got.Data.(string) != "G" {
			t.Errorf("got unexcepted data got %v, want %v", got.Data.(string), "G")
		}
	}

	ip3 := "86.100.32.255"
	got = matcher.Match(net.ParseIP(ip3))
	if got == nil {
		t.Errorf("%v should match but not", ip3)
	} else {
		if got.Data.(string) != "C" {
			t.Errorf("got unexcepted data got %v, want %v", got.Data.(string), "C")
		}
	}

	ip4 := "86.100.32.1"
	got = matcher.Match(net.ParseIP(ip4))
	if got == nil {
		t.Errorf("%v should match but not", ip4)
	} else {
		if got.Data.(string) != "A" {
			t.Errorf("got unexcepted data got %v, want %v", got.Data.(string), "A")
		}
	}
}

func TestPortmatcher(t *testing.T) {
	matcher := NewPortMatcher()
	item1, _ := StringToPortRange("0-500", nil)
	item2, _ := StringToPortRange("600-700", nil)
	item3, _ := StringToPortRange("802", nil)
	item4, _ := StringToPortRange("501-700", nil)
	item5, _ := StringToPortRange("555-800", nil)
	matcher.Add(item1)
	matcher.Add(item2)
	matcher.Add(item3)
	matcher.Add(item4)
	matcher.Add(item5)

	got := matcher.Match(501)
	if got == nil {
		t.Errorf("port %v should match but not", 501)
	}
	got = matcher.Match(801)
	if got != nil {
		t.Errorf("port %v should not match but match", 801)
	}
	got = matcher.Match(600)
	if got == nil {
		t.Errorf("port %v should match but not", 600)
	}
	got = matcher.Match(900)
	if got != nil {
		t.Errorf("port %v should not match but match", 900)
	}
}
