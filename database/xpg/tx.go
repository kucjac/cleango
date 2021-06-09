package xpg

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/kucjac/cleango/cgerrors"
)

// RunInTransaction is just like orm's RunInTransaction except
// it enforces errdef.ErrSet
func RunInTransaction(ctx context.Context, db orm.DB, fn func(*pg.Tx) error) error {
	var err error
	switch x := db.(type) {
	case *pg.DB:
		err = x.RunInTransaction(ctx, func(tx *pg.Tx) error {
			return fn(tx)
		})
	case *pg.Tx:
		err = fn(x)
	default:
		return cgerrors.ErrInternal("transaction type unknown. requires *pg.DB orf *pg.TX")
	}
	if err != nil {
		return defaultDriver.Err(err)
	}
	return nil
}
