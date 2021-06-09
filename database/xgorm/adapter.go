package xgorm

import (
	"context"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/database"
	"gorm.io/gorm"
)

// Adapter is the gorm structure used as a template for the gorm based adapters.
type Adapter struct {
	DB         *gorm.DB
	RetryCount int
	driver     database.Driver
}

// NewAdapter creates a new gorm based adapter. If the input driver is not a GORM driver,
// it would be wrapped by the xgorm.driver.
func NewAdapter(db *gorm.DB, drv database.Driver, retryCount int) (*Adapter, error) {
	wrappedDriver, err := NewDriver(drv)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		DB:         db,
		driver:     wrappedDriver,
		RetryCount: retryCount,
	}, nil
}

// Health implements health check for the service.
func (g *Adapter) Health(ctx context.Context) error {
	return g.Do(ctx, g.DB, func(ctx context.Context, db *gorm.DB) error {
		sqlDB, err := db.WithContext(ctx).DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	})
}

// Do executes given function 'f' so that if there would be a connection error it would try to retry the query
// up to the limit set on initialization.
func (g *Adapter) Do(ctx context.Context, db *gorm.DB, f func(ctx context.Context, db *gorm.DB) error) error {
	var err error
	for i := 0; i < g.RetryCount; i++ {
		if err = f(ctx, db); err != nil {
			if g.driver.CanRetry(err) {
				continue
			}
			return err
		}
		break
	}
	return err
}

// Err parses err using adapter driver.
func (g *Adapter) Err(err error) error {
	code := g.driver.ErrorCode(err)
	switch code {
	case cgerrors.ErrorCode_AlreadyExists:
		return cgerrors.ErrAlreadyExistsf("already exists: %v", err)
	case cgerrors.ErrorCode_NotFound:
		return cgerrors.ErrNotFound("not found")
	default:
		return cgerrors.New("", err.Error(), code)
	}
}

// Driver gets the adapter driver.
func (g *Adapter) Driver() database.Driver {
	return g.driver
}
