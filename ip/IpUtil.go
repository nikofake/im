package ip

import (
	"net"
)

func GetInternal() (string) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, a := range address {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return  ""
}
