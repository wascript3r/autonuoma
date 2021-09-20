package pgsql

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/wascript3r/autonuoma/pkg/domain"
)

type ErrCode string

const (
	UniqueViolationErrCode ErrCode = "23505"
	CheckViolationErrCode  ErrCode = "23514"
)

func ParseSQLError(err error) error {
	switch err {
	case sql.ErrNoRows:
		return domain.ErrNotFound
	}
	return err
}

func ParsePgError(err error) error {
	if e, ok := err.(*pq.Error); ok {
		switch ErrCode(e.Code) {
		case UniqueViolationErrCode:
			return domain.ErrExists

		case CheckViolationErrCode:
			return domain.ErrInvalidItem
		}
	}
	return err
}
