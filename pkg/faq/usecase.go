package faq

import (
	"context"
)

type Usecase interface {
	GetAll(ctx context.Context) (*GetAllRes, error)
}
