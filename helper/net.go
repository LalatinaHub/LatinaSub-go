package helper

import "net"

func GetFreePort() uint {
	var l *net.TCPListener
	for {
		addr, _ := net.ResolveTCPAddr("tcp", "localhost:0")
		if addr == nil {
			return GetFreePort()
		} else {
			l, _ = net.ListenTCP("tcp", addr)
			defer l.Close()

			break
		}
	}

	return uint(l.Addr().(*net.TCPAddr).Port)
}
