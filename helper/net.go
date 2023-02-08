package helper

import "net"

func GetFreePort() uint {
	var l *net.TCPListener
	for {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return GetFreePort()
		} else {
			l, err = net.ListenTCP("tcp", addr)
			if err != nil {
				return GetFreePort()
			}
			defer l.Close()

			break
		}
	}

	return uint(l.Addr().(*net.TCPAddr).Port)
}

func GetFreePortsLength() int {
	var freePorts []uint
	for {
		freePort := GetFreePort()
		for _, port := range freePorts {
			if port == freePort {
				return len(freePorts)
			}
		}
		freePorts = append(freePorts, freePort)
	}
}
