package ipmatcher

import (
	"strconv"
	"strings"

	"github.com/google/btree"
)

type Port int32

func (p Port) Less(then btree.Item) bool {
	switch np := then.(type) {
	case *PortRange:
		return p < np.StartPort
	case Port:
		return p < np
	}
	return false
}

type PortRange struct {
	StartPort Port
	EndPort   Port
	Data      interface{}
}

func (p *PortRange) Less(then btree.Item) bool {
	switch np := then.(type) {
	case *PortRange:
		if p.StartPort < np.StartPort {
			return true
		} else if p.StartPort == np.StartPort {
			return p.EndPort < np.EndPort
		} else {
			return false
		}
	case Port:
		return p.EndPort < np
	}
	return false
}

// StringToPortRange translate string to port range
func StringToPortRange(s string, val interface{}) (item *PortRange, err error) {
	s = strings.Trim(s, " \t\n_–+")
	if strings.Contains(s, "-") {
		if arr := strings.Split(s, "-"); len(arr) == 2 {
			sp, _ := strconv.Atoi(arr[0])
			ep, _ := strconv.Atoi(arr[1])
			item = &PortRange{
				StartPort: Port(sp),
				EndPort:   Port(ep),
				Data:      val,
			}
		}
	} else {
		sp, _ := strconv.Atoi(s)
		item = &PortRange{
			StartPort: Port(sp),
			EndPort:   Port(sp),
			Data:      val,
		}
	}
	return
}
