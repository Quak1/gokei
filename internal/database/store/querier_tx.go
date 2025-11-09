package store

import (
	"database/sql"
)

type QuerierTx interface {
	WithTx(tx *sql.Tx) QuerierTx
	Querier
}

type QueriesWrapper struct {
	Queries
}

func NewQueriesWrapper(db DBTX) *QueriesWrapper {
	return NewQueriesWrapperTx(New(db))
}

func NewQueriesWrapperTx(q *Queries) *QueriesWrapper {
	return &QueriesWrapper{
		*q,
	}
}

func (q *QueriesWrapper) WithTx(tx *sql.Tx) QuerierTx {
	return NewQueriesWrapperTx(q.Queries.WithTx(tx))
}

var _ QuerierTx = (*QueriesWrapper)(nil)
