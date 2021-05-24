package sqlxes

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kucjac/cleango/cgerrors"
)

// Transaction is the
type Transaction struct {
	id string
	storage
	done bool
}

func (t *Transaction) Conn() (*sqlx.Tx, error) {
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
	tx, err := t.Conn()
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
