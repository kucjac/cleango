package sqlxes

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/eventsource"
)

// Compile time check if Transaction implements eventsource.Storage.
var _ eventsource.Storage = (*Transaction)(nil)

// Transaction is the implementation of the
type Transaction struct {
	id string
	storage
	done bool
}

// As sets the destination with the *sqlx.Tx implementation.
func (t *Transaction) As(dst interface{}) error {
	tx, err := t.txConn()
	if err != nil {
		return err
	}
	txx, ok := dst.(**sqlx.Tx)
	if !ok {
		return cgerrors.ErrInternalf("provided invalid input type: %T, wanted: **sqlx.Tx", dst)
	}
	*txx = tx
	return nil
}

// Done checks if the transaction is already done.
func (t *Transaction) Done() bool {
	return t.done
}

func (t *Transaction) txConn() (*sqlx.Tx, error) {
	tx, ok := t.conn.(*sqlx.Tx)
	if !ok {
		return nil, cgerrors.ErrInternalf("unknown type of sqlx based eventsource transaction conn: %T", t.conn)
	}
	return tx, nil
}

// Commit commits the transaction.
func (t *Transaction) Commit() error {
	if t.done {
		return cgerrors.ErrInternalf("transaction '%s' is already done", t.id)
	}
	tx, err := t.txConn()
	if err != nil {
		return err
	}
	err = t.tryTx(context.Background(), tx, func(ctx context.Context, tx *sqlx.Tx) error {
		return tx.Commit()
	})
	if err != nil {
		return t.Err(err)
	}
	t.done = true
	return nil
}

// Rollback the transaction.
func (t *Transaction) Rollback() error {
	if t.done {
		return cgerrors.ErrInternalf("transaction '%s' is already done", t.id)
	}
	tx, err := t.txConn()
	if err != nil {
		return err
	}
	err = t.tryTx(context.Background(), tx, func(ctx context.Context, tx *sqlx.Tx) error {
		return tx.Rollback()
	})
	if err != nil {
		return t.Err(err)
	}
	t.done = true
	return nil
}
