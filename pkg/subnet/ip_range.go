package subnet

import (
	"fmt"
	"net"
	"strings"
)

// GetUsedIPsInSubnet spočítá počet dostupných IP adres v daném rozsahu
func GetUsedIPsInSubnet(ipRange string) (int, error) {
	parts := strings.Split(ipRange, "-")

	if len(parts) == 1 {
		// Single IP address or CIDR notation
		_, ipnet, err := net.ParseCIDR(parts[0])
		if err != nil {
			return 0, err
		}

		if ipnet != nil {
			// CIDR notation
			ones, bits := ipnet.Mask.Size()
			return 1 << uint(bits-ones), nil
		}

		// Single IP address
		return 1, nil
	}

	// IP range
	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	endIP := net.ParseIP(strings.TrimSpace(parts[1]))

	if startIP == nil || endIP == nil {
		return 0, fmt.Errorf("invalid IP range")
	}

	start := ipToInt(startIP)
	end := ipToInt(endIP)

	return end - start + 1, nil
}

func ipToInt(ip net.IP) int {
	// Convert IP address to a 32-bit integer
	ip = ip.To4()
	return int(ip[0])<<24 | int(ip[1])<<16 | int(ip[2])<<8 | int(ip[3])
}
