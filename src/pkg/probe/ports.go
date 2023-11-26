package probe

import (
	"bufio"
	_ "embed"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

//go:embed ports.txt
var data string

var portMap map[uint]string

func Ports() []uint {
	return lo.Keys(PortMap())
}

func PortMap() map[uint]string {
	if portMap == nil {
		portMap = map[uint]string{}

		scanner := bufio.NewScanner(strings.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line[0] == '#' {
				continue
			}

			if port, name, ok := strings.Cut(line, ":"); ok {
				portMap[uint(lo.Must(strconv.Atoi(port)))] = name
			}
		}
	}

	return portMap
}

func ServiceName(port uint) string {
	return PortMap()[port]
}
