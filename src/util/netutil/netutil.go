package netutil

import (
	"net"
	"strings"

	"github.com/google/gopacket/macs"
)

func GetVendor(macAddress string) (string, bool) {
	hardwareAddr, err := net.ParseMAC(macAddress)
	if err != nil {
		return "", false
	}

	var prefix [3]byte
	copy(prefix[:], hardwareAddr[:3])

	vendor, found := macs.ValidMACPrefixMap[prefix]
	return vendor, found
}

func UnqualifyHostname(hostname string) string {
	hostname = strings.TrimSuffix(hostname, ".")

	index := strings.LastIndex(hostname, ".")
	if index == -1 {
		return hostname
	}

	return hostname[:index]
}
