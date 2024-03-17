package db

import "fmt"

func (self *Mapping) AuditName() string {
	return fmt.Sprintf("%s ↔ %s (%s)", self.MACAddress, self.IPAddress, self.Label)
}
