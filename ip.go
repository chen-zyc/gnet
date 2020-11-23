package gnet

import (
	"net"
)

// IsIPv4 returns true if ip is v4.
func IsIPv4(ip net.IP) bool {
	return ip != nil && ip.To4() != nil
}

// IsIPv6 returns true if ip is v6.
func IsIPv6(ip net.IP) bool {
	return ip != nil && ip.To4() == nil
}
