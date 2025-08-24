package sysinfo

import (
	"fmt"
	"net"
)

func PrintAllAddresses(port int) {

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		return
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // skip down interfaces
		}

		//fmt.Println("Interface Name:", iface.Name)

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("  Error getting addresses:", err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			if ip.To4() != nil {
				fmt.Printf("availables also via: http://%s:%d\n", ip.String(), port)
			}
		}
	}

}
