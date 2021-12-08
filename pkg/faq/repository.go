package faq

import (
	"context"

	"github.com/wascript3r/autonuoma/pkg/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*domain.FAQ, error)
}
