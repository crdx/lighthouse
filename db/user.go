package db

import (
	"fmt"

	"crdx.org/lighthouse/pkg/constants"
)

func (self *User) AuditName() string {
	return fmt.Sprintf("%s (ID: %d)", self.Username, self.ID)
}

func (self *User) IsAdmin() bool  { return self.Role >= constants.RoleAdmin }
func (self *User) IsEditor() bool { return self.Role >= constants.RoleEditor }
func (self *User) IsViewer() bool { return self.Role >= constants.RoleViewer }
