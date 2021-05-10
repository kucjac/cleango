package xgorm

import (
	"errors"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/xservice"
	"gorm.io/gorm"
)

var _ xservice.Driver = (*driver)(nil)

// NewDriver creates a new gorm driver.
func NewDriver(baseDriver xservice.Driver) (xservice.Driver, error) {
	if baseDriver == nil {
		return nil, errors.New("xgorm provided nil base driver")
	}
	return &driver{base: baseDriver}, nil
}

type driver struct {
	base xservice.Driver
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
