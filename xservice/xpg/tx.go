package xpg

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/kucjac/cleango/errors"
)

// RunInTransaction is just like orm's RunInTransaction except
// it enforces errdef.ErrSet
func RunInTransaction(ctx context.Context, db orm.DB, fn func(*pg.Tx) error) error {
	switch x := db.(type) {
	case *pg.DB:
		err := x.RunInTransaction(ctx, func(tx *pg.Tx) error {
			return fn(tx)
		})
		switch x := err.(type) {
		case nil:
			return nil
		case *errors.Error:
			return x
		default:
			return errors.ErrInternalf("postgres database unknown error: %v", x)
		}
	case *pg.Tx:
		return fn(x)
	default:
		return errors.ErrInternal("transaction type unknown. requires *pg.DB orf *pg.TX")
	}
}
