package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed static-dhcp.conf
var conf string

var names map[string]string

func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func main() {
	names = map[string]string{}

	// Format: dhcp-host=AA:BB:CC:DD:EE:FF,192.168.1.123,hostname
	re := regexp.MustCompile(`dhcp-host=([^,]+),\d+\.\d+\.\d+\.\d+,(.*)`)

	for _, line := range strings.Split(conf, "\n") {
		matches := re.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		names[strings.ToUpper(matches[1])] = strings.TrimSpace(matches[2])
	}

	for macAddress, name := range names {
		fmt.Printf(
			"UPDATE devices D INNER JOIN adapters A ON A.device_id = D.id SET D.name = %s WHERE A.mac_address = %s;\n",
			quote(name),
			quote(macAddress),
		)
	}
}
