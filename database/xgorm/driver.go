package xgorm

import (
	"errors"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database"
	"gorm.io/gorm"
)

var _ database.Driver = (*driver)(nil)

// NewDriver creates a new gorm driver.
func NewDriver(baseDriver database.Driver) (database.Driver, error) {
	if baseDriver == nil {
		return nil, errors.New("xgorm provided nil base driver")
	}
	return wrapDriver(baseDriver), nil
}

func wrapDriver(drv database.Driver) database.Driver {
	d, ok := drv.(*driver)
	if ok {
		return d
	}
	return &driver{base: drv}
}

type driver struct {
	base database.Driver
}

func (d *driver) ErrorCode(err error) cgerrors.ErrorCode {
	switch err {
	case gorm.ErrRecordNotFound:
		return cgerrors.ErrorCode_NotFound
	case gorm.ErrInvalidTransaction, gorm.ErrNotImplemented, gorm.ErrMissingWhereClause, gorm.ErrUnsupportedRelation, gorm.ErrPrimaryKeyRequired, gorm.ErrModelValueRequired, gorm.ErrUnsupportedDriver, gorm.ErrDryRunModeUnsupported, gorm.ErrInvalidDB, gorm.ErrRegistered, gorm.ErrInvalidField:
		return cgerrors.ErrorCode_Internal
	case gorm.ErrEmptySlice, gorm.ErrInvalidData, gorm.ErrInvalidValue, gorm.ErrInvalidValueOfLength:
		return cgerrors.ErrorCode_InvalidArgument
	default:
		return d.base.ErrorCode(err)
	}
}

func (d *driver) CanRetry(err error) bool {
	return d.base.CanRetry(err)
}
