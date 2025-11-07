package netutil

import (
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/google/gopacket/macs"
	"github.com/samber/lo"
)

// GetVendor returns the vendor for a mac address, and true if it was found in the gopacket
// database.
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

// UnqualifyHostname returns a hostname without the domain portion and removes any trailing periods,
// if any.
func UnqualifyHostname(hostname string) string {
	hostname = strings.TrimSuffix(hostname, ".")

	index := strings.LastIndex(hostname, ".")
	if index == -1 {
		return hostname
	}

	return hostname[:index]
}

// ExpandIPNet takes a network and returns a slice of all IPs that belong to that network.
func ExpandIPNet(ipNet *net.IPNet) []net.IP {
	var ips []net.IP

	n := binary.BigEndian.Uint32([]byte(ipNet.IP))
	mask := binary.BigEndian.Uint32([]byte(ipNet.Mask))
	ip := n & mask
	broadcast := ip | ^mask

	for ip++; ip < broadcast; ip++ {
		var buffer [4]byte
		binary.BigEndian.PutUint32(buffer[:], ip)
		ips = append(ips, net.IP(buffer[:]))
	}

	return ips
}

// IPNetTooLarge returns true if a network is considered too large to scan.
//
// Anything larger than a /24 is considered too large.
func IPNetTooLarge(ipNet *net.IPNet) bool {
	return ipNet.Mask[0] != 0xff ||
		ipNet.Mask[1] != 0xff ||
		ipNet.Mask[2] != 0xff
}

// findInterface returns the first interface that is considered "up" and is not the loopback
// interface.
func findInterface() (*net.Interface, bool) {
	interfaces := lo.Must(net.Interfaces())

	for _, iface := range interfaces {
		isRunning := iface.Flags&net.FlagRunning == net.FlagRunning
		hasAssignedIps := lo.CountBy(lo.Must(iface.Addrs()), func(addr net.Addr) bool {
			network, ok := addr.(*net.IPNet)
			return ok && network.IP.To4() != nil
		}) > 0

		if iface.Name != "lo" && isRunning && hasAssignedIps {
			return &iface, true
		}
	}

	return nil, false
}

// findIPNet returns the first IPv4 network for an interface.
func findIPNet(iface *net.Interface) (*net.IPNet, bool) {
	addresses := lo.Must(iface.Addrs())

	for _, address := range addresses {
		if network, ok := address.(*net.IPNet); ok {
			if ip := network.IP.To4(); ip != nil {
				return &net.IPNet{
					IP:   ip,
					Mask: network.Mask[len(network.Mask)-4:],
				}, true
			}
		}
	}

	return nil, false
}

// FindNetwork returns the first active interface and network.
func FindNetwork() (*net.Interface, *net.IPNet, error) {
	iface, found := findInterface()
	if !found {
		return nil, nil, fmt.Errorf("no interface found")
	}

	ipNet, found := findIPNet(iface)
	if !found {
		return nil, nil, fmt.Errorf("no network found")
	} else if IPNetTooLarge(ipNet) {
		return nil, nil, fmt.Errorf("network too large")
	}

	return iface, ipNet, nil
}

func IsValidMAC(s string) bool {
	re := regexp.MustCompile(`^[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}:[0-9A-F]{2}$`)
	return re.MatchString(strings.TrimSpace(s))
}

func ParseMACList(values string) ([]string, bool) {
	var macs []string
	for value := range strings.SplitSeq(values, ",") {
		value := strings.TrimSpace(value)
		if !IsValidMAC(value) {
			return nil, false
		}
		macs = append(macs, value)
	}
	return macs, true
}
