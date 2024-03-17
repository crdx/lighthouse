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

var portMap map[int64]string

func Ports() []int64 {
	return lo.Keys(PortMap())
}

func PortMap() map[int64]string {
	if portMap == nil {
		portMap = map[int64]string{}

		scanner := bufio.NewScanner(strings.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line[0] == '#' {
				continue
			}

			if port, name, ok := strings.Cut(line, ":"); ok {
				portMap[int64(lo.Must(strconv.Atoi(port)))] = name
			}
		}
	}

	return portMap
}

func ServiceName(port int64) string {
	return PortMap()[port]
}
