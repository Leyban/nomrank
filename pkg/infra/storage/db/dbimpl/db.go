package dbimpl

import (
	"context"
	"database/sql"
	"fmt"
	"nomrank/pkg/infra/storage/db"

	"github.com/jmoiron/sqlx"
)

type ContextSessionKey struct{}

type sqlxdb struct {
	db *sqlx.DB
}

func (gs *sqlxdb) Get(ctx context.Context, dest any, query string, args ...any) error {
	return gs.db.GetContext(ctx, dest, gs.db.Rebind(query), args...)
}

func (gs *sqlxdb) Select(ctx context.Context, dest any, query string, args ...any) error {
	return gs.db.SelectContext(ctx, dest, gs.db.Rebind(query), args...)
}

func (gs *sqlxdb) Query(ctx context.Context, dest any, query string, args ...any) (*sql.Rows, error) {
	return gs.db.QueryContext(ctx, gs.db.Rebind(query), args...)
}

func (gs *sqlxdb) Exec(ctx context.Context, dest any, query string, args ...any) (sql.Result, error) {
	return gs.db.ExecContext(ctx, gs.db.Rebind(query), args...)
}

func (gs *sqlxdb) NamedExec(ctx context.Context, dest any, query string, args any) (sql.Result, error) {
	return gs.db.NamedExecContext(ctx, gs.db.Rebind(query), args)
}

func (gs *sqlxdb) Close() error {
	if err := gs.db.Close(); err != nil {
		return err
	}
	return nil
}

func (gs *sqlxdb) Ping() error {
	if err := gs.db.Ping(); err != nil {
		return err
	}
	return nil
}

// Transaction Mode
func (gs *sqlxdb) Beginx() (db.Tx, error) {
	tx, err := gs.db.Beginx()
	return &sqlxtx{sqlxtx: tx}, err
}

func (gs *sqlxdb) checkSession(ctx context.Context) (db.Tx, bool) {
	value := ctx.Value(ContextSessionKey{})

	var tx db.Tx

	sess, ok := value.(db.Tx)
	if ok {
		return sess, false
	}

	tx, err := gs.Beginx()
	if err != nil {
		return tx, false
	}

	return tx, true
}

func (gs *sqlxdb) inTransaction(ctx context.Context, callback func(db.Tx) error) error {
	tx, isNew := gs.checkSession(ctx)
	if !isNew {
		return fmt.Errorf("Multiple transaction instances")
	}

	err := callback(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (gs *sqlxdb) WithTransaction(ctx context.Context, fn func(sess db.Tx) error) error {
	return gs.inTransaction(ctx, fn)
}

func (gs *sqlxdb) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return gs.inTransaction(ctx, func(sess db.Tx) error {
		withValue := context.WithValue(ctx, ContextSessionKey{}, sess)
		return fn(withValue)
	})
}
