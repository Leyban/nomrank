package db

import (
	"context"
	"database/sql"
)

type DB interface {
	Get(ctx context.Context, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	Query(ctx context.Context, dest any, query string, args ...any) (*sql.Rows, error)
	Exec(ctx context.Context, dest any, query string, args ...any) (sql.Result, error)
	NamedExec(ctx context.Context, dest any, query string, args any) (sql.Result, error)
	Ping() error
	WithTransaction(ctx context.Context, fn func(sess Tx) error) error
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Tx interface {
	NamedExec(ctx context.Context, query string, arg any) (sql.Result, error)
	Exec(ctx context.Context, query string, arg ...any) (sql.Result, error)
	Query(ctx context.Context, query string, arg ...any) (*sql.Rows, error)
	Get(ctx context.Context, dest any, query string, arg ...any) error
	Select(ctx context.Context, dest any, query string, arg ...any) error
	Rollback() error
	Commit() error
}
