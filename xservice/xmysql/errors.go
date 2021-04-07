package xmysql

import (
	"github.com/go-sql-driver/mysql"
)

// IsDuplicatedError checks if given error states for the mysql duplicated input error.
func IsDuplicatedError(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	if me.Number == 1062 {
		return true
	}
	return false
}
