package xgorm

import (
	"errors"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database"
	"gorm.io/gorm"
)

// Compile time check if Driver implements database.Driver.
var _ database.Driver = (*Driver)(nil)

// NewDriver creates a new gorm driver.
func NewDriver(baseDriver database.Driver) (*Driver, error) {
	if baseDriver == nil {
		return nil, errors.New("xgorm provided nil base driver")
	}
	return wrapDriver(baseDriver), nil
}

func wrapDriver(drv database.Driver) *Driver {
	d, ok := drv.(*Driver)
	if ok {
		return d
	}
	return &Driver{base: drv}
}

// Driver is the gorm implementation for the driver name.
type Driver struct {
	base database.Driver
}

// DriverName gets the name of the driver.
func (d *Driver) DriverName() string {
	return "gorm+" + d.base.DriverName()
}

// ErrorCode gets the error code for given error.
func (d *Driver) ErrorCode(err error) cgerrors.ErrorCode {
	switch err {
	case gorm.ErrRecordNotFound:
		return cgerrors.CodeNotFound
	case gorm.ErrInvalidTransaction, gorm.ErrNotImplemented, gorm.ErrMissingWhereClause, gorm.ErrUnsupportedRelation, gorm.ErrPrimaryKeyRequired, gorm.ErrModelValueRequired, gorm.ErrUnsupportedDriver, gorm.ErrDryRunModeUnsupported, gorm.ErrInvalidDB, gorm.ErrRegistered, gorm.ErrInvalidField:
		return cgerrors.CodeInternal
	case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrInvalidValue, gorm.ErrInvalidValueOfLength:
		return cgerrors.CodeInvalidArgument
	default:
		return d.base.ErrorCode(err)
	}
}

// CanRetry checks if the error relted query could be retried.
func (d *Driver) CanRetry(err error) bool {
	return d.base.CanRetry(err)
}
