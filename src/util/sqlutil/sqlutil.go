package sqlutil

import (
	"database/sql"
	"database/sql/driver"
)

type NullUint struct {
	Uint  uint
	Valid bool
}

func (self *NullUint) Scan(value any) error {
	if value == nil {
		self.Uint, self.Valid = 0, false
		return nil
	}

	self.Valid = true
	i := sql.NullInt64{}
	if err := i.Scan(value); err != nil {
		return err
	}

	self.Uint = uint(i.Int64)
	return nil
}

func (self NullUint) Value() (driver.Value, error) {
	if !self.Valid {
		return nil, nil
	}

	return int64(self.Uint), nil
}
