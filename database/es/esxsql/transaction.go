package esxsql

import (
	"context"

	"github.com/kucjac/cleango/cgerrors"
	eventsource "github.com/kucjac/cleango/database/es"
	"github.com/kucjac/cleango/database/xsql"
)

// Compile time check if Transaction implements eventsource.Storage.
var _ eventsource.TxStorage = (*Transaction)(nil)

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
	txx, ok := dst.(**xsql.Tx)
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

// Commit commits the transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	if t.done {
		return cgerrors.ErrInternalf("transaction '%s' is already done", t.id)
	}
	tx, err := t.txConn()
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return t.Err(err)
	}
	t.done = true
	return nil
}

// Rollback the transaction.
func (t *Transaction) Rollback(_ context.Context) error {
	if t.done {
		return cgerrors.ErrInternalf("transaction '%s' is already done", t.id)
	}
	tx, err := t.txConn()
	if err != nil {
		return err
	}
	if err = tx.Rollback(); err != nil {
		return t.Err(err)
	}
	t.done = true
	return nil
}

func (t *Transaction) txConn() (*xsql.Tx, error) {
	tx, ok := t.conn.(*xsql.Tx)
	if !ok {
		return nil, cgerrors.ErrInternalf("unknown type of sqlx based eventsource transaction conn: %T", t.conn)
	}
	return tx, nil
}
