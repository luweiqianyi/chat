package ws

import (
	"fmt"
	"net"
)

func GetServerIP() (string, error) {
	// get all interfaces of network
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// find first ipv4 address
	for _, netInterface := range interfaces {
		if netInterface.Flags&net.FlagLoopback == 0 && netInterface.Flags&net.FlagUp != 0 {
			addrArray, _ := netInterface.Addrs()
			for _, addr := range addrArray {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						return ipNet.IP.String(), nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("no IP address found")
}

func GenerateWebsocketClientID(ip string, port string, connTime string) string {
	return fmt.Sprintf("%v:%v %v", ip, port, connTime)
}
