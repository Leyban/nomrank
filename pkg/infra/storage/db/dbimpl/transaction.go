package dbimpl

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type sqlxtx struct {
	sqlxtx *sqlx.Tx
}

func (gtx *sqlxtx) NamedExec(ctx context.Context, query string, arg any) (sql.Result, error) {
	return gtx.sqlxtx.NamedExecContext(ctx, gtx.sqlxtx.Rebind(query), arg)
}

func (gtx *sqlxtx) Exec(ctx context.Context, query string, arg ...any) (sql.Result, error) {
	return gtx.sqlxtx.ExecContext(ctx, gtx.sqlxtx.Rebind(query), arg...)
}

func (gtx *sqlxtx) Query(ctx context.Context, query string, arg ...any) (*sql.Rows, error) {
	return gtx.sqlxtx.QueryContext(ctx, gtx.sqlxtx.Rebind(query), arg...)
}

func (gtx *sqlxtx) Get(ctx context.Context, dest any, query string, arg ...any) error {
	return gtx.sqlxtx.GetContext(ctx, dest, gtx.sqlxtx.Rebind(query), arg...)
}

func (gtx *sqlxtx) Select(ctx context.Context, dest any, query string, arg ...any) error {
	return gtx.sqlxtx.SelectContext(ctx, dest, gtx.sqlxtx.Rebind(query), arg...)
}

func (gtx *sqlxtx) Rollback() error {
	return gtx.sqlxtx.Rollback()
}

func (gtx *sqlxtx) Commit() error {
	return gtx.sqlxtx.Commit()
}
