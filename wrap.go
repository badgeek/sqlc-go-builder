package sqlc

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/pkg/errors"
)

var _ DBTX = (*wrappedDB)(nil)

func Wrap(db DBTX) DBTX {
	return &wrappedDB{db}
}

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type wrappedDB struct {
	DBTX
}

func (w *wrappedDB) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	var err error
	if b, ok := BuilderFrom(ctx); ok {
		query, args, err = b.Build(query, args...)
	}
	if err != nil {
		return pgconn.CommandTag{}, errors.Wrap(err, "could not build query")
	}
	return w.DBTX.Exec(ctx, query, args...)
}

func (w *wrappedDB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	var err error
	if b, ok := BuilderFrom(ctx); ok {
		query, args, err = b.Build(query, args...)
	}
	if err != nil {
		return nil, errors.Wrap(err, "could not build query")
	}
	return w.DBTX.Query(ctx, query, args...)
}

func (w *wrappedDB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	var err error
	if b, ok := BuilderFrom(ctx); ok {
		if queryNew, argsNew, err := b.Build(query, args...); err == nil {
			query = queryNew
			args = argsNew
		}
	}
	if err != nil {
		fmt.Printf("could not build query: %s", err)
	}
	return w.DBTX.QueryRow(ctx, query, args...)
}
