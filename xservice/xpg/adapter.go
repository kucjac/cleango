package xpg

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/kucjac/cleango/cgerrors"
)

// Adapter is the go-pg adapter implementation.
type Adapter struct {
	DB         orm.DB
	Driver     *PGDriver
	RetryCount int
}

// NewAdapter creates a new go-pg adapter.
func NewAdapter(db orm.DB) *Adapter {
	return &Adapter{
		DB:         db,
		Driver:     NewDriver(),
		RetryCount: 10,
	}
}

// Do executes given function with taking care about the retries of the connection.
func (a *Adapter) Do(ctx context.Context, db orm.DB, fn func(context.Context, orm.DB) error) error {
	for i := 0; i < a.RetryCount; i++ {
		err := fn(ctx, db)
		if err != nil {
			if a.Driver.CanRetry(err) {
				continue
			}
			return err
		}
		break
	}
	return nil
}

// DoTx executes given function in a transaction taking care about the retries of the connection.
func (a *Adapter) DoTx(ctx context.Context, db orm.DB, fn func(context.Context, *pg.Tx) error) error {
	switch x := db.(type) {
	case *pg.DB:
		for i := 0; i < a.RetryCount; i++ {
			err := x.RunInTransaction(ctx, func(tx *pg.Tx) error {
				return fn(ctx, tx)
			})
			if err != nil {
				if a.Driver.CanRetry(err) {
					continue
				}
				return err
			}
			break
		}
		return nil
	case *pg.Tx:
		for i := 0; i < a.RetryCount; i++ {
			if err := fn(ctx, x); err != nil {
				if a.Driver.CanRetry(err) {
					continue
				}
				return err
			}
			break
		}
		return nil
	default:
		return cgerrors.ErrInternal("transaction type unknown. requires *pg.DB orf *pg.Tx")
	}
}

// Err handlers given error returning a cgerrors.Error
func (a *Adapter) Err(err error) *cgerrors.Error {
	return a.Driver.Err(err)
}
